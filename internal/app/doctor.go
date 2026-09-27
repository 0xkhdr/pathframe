package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	pfcontext "github.com/0xkhdr/pathframe/internal/context"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/integrations/claude"
	"github.com/0xkhdr/pathframe/internal/integrations/codex"
	"github.com/0xkhdr/pathframe/internal/recovery"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

const abandonedLeaseAfter = 24 * time.Hour

type DoctorInput struct {
	Change string `json:"change,omitempty"`
	Repair bool   `json:"repair,omitempty"`
}

func (s Service) Doctor(input DoctorInput) (recovery.Report, error) {
	start := s.Dir
	if start == "" {
		start = "."
	}
	root, err := store.Discover(start)
	if errors.Is(err, store.ErrProjectNotFound) {
		absolute, _ := filepath.Abs(start)
		return recovery.New(absolute, "", []recovery.Diagnosis{{Code: string(workflow.ReasonNoProject), Severity: recovery.SeverityError, Subject: ".pathframe", Evidence: err.Error(), HumanRequired: true}}), nil
	}
	if err != nil {
		return recovery.Report{}, err
	}
	changes, err := store.Changes(root)
	if err != nil {
		return recovery.Report{}, err
	}
	change := input.Change
	if change == "" {
		if len(changes) == 0 {
			return recovery.New(root, "", []recovery.Diagnosis{{Code: string(workflow.ReasonNoChange), Severity: recovery.SeverityWarning, Subject: ".pathframe/changes", Evidence: "no change exists", HumanRequired: true}}), nil
		}
		if len(changes) > 1 {
			return recovery.New(root, "", []recovery.Diagnosis{{Code: string(workflow.ReasonAmbiguousChange), Severity: recovery.SeverityError, Subject: ".pathframe/changes", Evidence: "multiple changes exist; select one", HumanRequired: true}}), nil
		}
		change = changes[0]
	}
	dir, err := store.ChangeDir(root, change)
	if err != nil {
		return recovery.Report{}, err
	}
	diagnoses, err := recovery.DiagnoseProjection(dir)
	if err != nil {
		return recovery.Report{}, err
	}
	leaseDiagnoses, err := recovery.DiagnoseLease(dir, now(), abandonedLeaseAfter)
	if err != nil {
		return recovery.Report{}, err
	}
	diagnoses = append(diagnoses, leaseDiagnoses...)
	if _, statErr := os.Stat(filepath.Join(root, ".pathframe", "integrations", "codex")); statErr == nil {
		for _, item := range codex.Doctor(root).Diagnostics {
			diagnoses = append(diagnoses, recovery.Diagnosis{Code: item.Code, Severity: recovery.SeverityError, Subject: item.Path, Evidence: item.Message, HumanRequired: true})
		}
	}
	if _, statErr := os.Stat(filepath.Join(root, ".pathframe", "integrations", "claude")); statErr == nil {
		for _, item := range claude.Doctor(root).Diagnostics {
			diagnoses = append(diagnoses, recovery.Diagnosis{Code: item.Code, Severity: recovery.SeverityError, Subject: item.Path, Evidence: item.Message, HumanRequired: true})
		}
	}
	if _, statErr := os.Stat(filepath.Join(dir, "change.yaml")); errors.Is(statErr, os.ErrNotExist) {
		diagnoses = append(diagnoses, recovery.Diagnosis{Code: "artifacts_missing", Severity: recovery.SeverityError, Subject: "change.yaml", Evidence: "authored change metadata is missing", HumanRequired: true})
	} else if statErr != nil {
		return recovery.Report{}, statErr
	} else if plan, issues := artifacts.Validate(dir); len(issues) > 0 {
		for _, issue := range issues {
			diagnoses = append(diagnoses, recovery.Diagnosis{Code: issue.Code, Severity: recovery.SeverityError, Subject: issue.Path, Evidence: issue.Message, HumanRequired: true})
		}
	} else {
		if _, completedErr := completedTasks(dir, plan); completedErr != nil {
			diagnoses = append(diagnoses, recovery.Diagnosis{Code: "run_history_invalid", Severity: recovery.SeverityError, Subject: "runs", Evidence: completedErr.Error(), HumanRequired: true})
		}
		for _, task := range plan.Tasks {
			if task.Role != "" && task.Role != "none" {
				role, roleErr := pfcontext.LoadRole(filepath.Join(root, ".pathframe", "roles", task.Role+".yaml"))
				if roleErr != nil || role.ID != task.Role {
					evidence := "role id does not match task"
					if roleErr != nil {
						evidence = roleErr.Error()
					}
					diagnoses = append(diagnoses, recovery.Diagnosis{Code: "unknown_role", Severity: recovery.SeverityError, Subject: task.ID, Evidence: evidence, HumanRequired: true})
				}
			}
			for _, ref := range sectionBullets(task.Sections["Required Reads"]) {
				_, unresolved, resolveErr := pfcontext.Resolve(root, dir, []pfcontext.Reference{{Layer: pfcontext.Change, Path: ref, Required: true}})
				if resolveErr != nil || len(unresolved) > 0 {
					evidence := "required reference cannot be resolved"
					if resolveErr != nil {
						evidence = resolveErr.Error()
					}
					diagnoses = append(diagnoses, recovery.Diagnosis{Code: "broken_reference", Severity: recovery.SeverityError, Subject: ref, Evidence: evidence, HumanRequired: true})
				}
			}
		}
		states, stateErr := delegation.RunStates(filepath.Join(dir, "runs"))
		if stateErr != nil {
			diagnoses = append(diagnoses, recovery.Diagnosis{Code: "run_history_invalid", Severity: recovery.SeverityError, Subject: "runs", Evidence: stateErr.Error(), HumanRequired: true})
		} else {
			for id, state := range states {
				if state == "verification_failed" {
					diagnoses = append(diagnoses, recovery.Diagnosis{Code: "verification_failed", Severity: recovery.SeverityWarning, Subject: id, Evidence: "latest verification failed; fix or request changes, then retry", HumanRequired: true})
				}
			}
		}
	}
	report := recovery.New(root, change, diagnoses)
	if !input.Repair {
		return report, nil
	}
	for _, diagnosis := range diagnoses {
		switch diagnosis.Repair {
		case recovery.RepairProjectionCode, recovery.RepairJournalTailCode:
			if err := recovery.RepairProjection(dir); err != nil {
				return report, fmt.Errorf("repair %s: %w", diagnosis.Repair, err)
			}
			report.Repaired = true
		case recovery.RepairLeaseCode:
			lease, active, leaseErr := delegation.ActiveLease(filepath.Join(dir, "runs"))
			if leaseErr != nil {
				return report, fmt.Errorf("repair abandoned lease: %w", leaseErr)
			}
			if !active {
				return report, errors.New("repair abandoned lease: active lease disappeared")
			}
			if _, releaseErr := s.ReleaseLease(ReleaseLeaseInput{Change: change, LeaseID: lease.ID}); releaseErr != nil {
				return report, releaseErr
			}
			report.Repaired = true
		}
	}
	refreshed, err := s.Doctor(DoctorInput{Change: change})
	if err != nil {
		return report, err
	}
	refreshed.Repaired = report.Repaired
	return refreshed, nil
}
