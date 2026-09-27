package app

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/verification"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

type RunVerificationInput struct {
	Change       string   `json:"change"`
	Task         string   `json:"task"`
	ChangedFiles []string `json:"changed_files,omitempty"`
	Workdir      string   `json:"workdir,omitempty"`
	TimeoutMS    int64    `json:"timeout_ms"`
	MaxBytes     int      `json:"max_bytes,omitempty"`
}

type RunVerificationOutput struct {
	Schema       string                  `json:"schema"`
	Change       string                  `json:"change"`
	Task         string                  `json:"task"`
	LeaseID      string                  `json:"lease_id"`
	Verification delegation.Verification `json:"verification"`
	Recovery     []string                `json:"recovery"`
}

func (s Service) RunVerification(ctx context.Context, input RunVerificationInput) (RunVerificationOutput, error) {
	if input.TimeoutMS <= 0 {
		return RunVerificationOutput{}, fmt.Errorf("timeout_ms must be positive")
	}
	view, err := s.Orient(input.Change)
	if err != nil {
		return RunVerificationOutput{}, err
	}
	root, err := store.Discover(s.Dir)
	if err != nil {
		return RunVerificationOutput{}, err
	}
	dir, err := store.ChangeDir(root, view.Change)
	if err != nil {
		return RunVerificationOutput{}, err
	}
	plan, issues := artifacts.Validate(dir)
	if len(issues) > 0 {
		return RunVerificationOutput{}, fmt.Errorf("approved plan is invalid; repair and reapprove before verification")
	}
	var task *artifacts.Task
	for i := range plan.Tasks {
		if plan.Tasks[i].ID == input.Task {
			task = &plan.Tasks[i]
			break
		}
	}
	if task == nil {
		return RunVerificationOutput{}, fmt.Errorf("task not found: %s", input.Task)
	}
	runsDir := filepath.Join(dir, "runs")
	if (view.Phase == workflow.PhaseReady || view.Phase == workflow.PhaseBlocked) && task.ExecutionPolicy == "brain" {
		lease, acquireErr := delegation.Acquire(runsDir, view.Change, task.ID, "brain", now())
		if acquireErr != nil {
			return RunVerificationOutput{}, acquireErr
		}
		if err := s.appendDelegationTransition(dir, workflow.ActionExecute, workflow.ActorBrain, "Brain began directly authorized task "+task.ID); err != nil {
			_ = delegation.Release(runsDir, lease)
			return RunVerificationOutput{}, err
		}
		result := delegation.Result{Schema: delegation.ResultSchema, LeaseID: lease.ID, Change: view.Change, Task: task.ID, Status: delegation.ResultCompleted, Summary: "Brain submitted directly authorized work for verification", ChangedFiles: input.ChangedFiles, Verification: []delegation.VerificationReport{}, Discoveries: []string{}, Questions: []string{}, Risks: []string{}}
		if err := delegation.AppendRun(runsDir, lease.ID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: lease.ID, State: "submitted", Timestamp: now(), Result: &result}); err != nil {
			return RunVerificationOutput{}, err
		}
		if err := s.appendDelegationTransition(dir, workflow.ActionReview, workflow.ActorBrain, "Brain submitted direct work for verification"); err != nil {
			return RunVerificationOutput{}, err
		}
		if err := delegation.Release(runsDir, lease); err != nil {
			return RunVerificationOutput{}, err
		}
	} else if view.Phase != workflow.PhaseReviewing {
		return RunVerificationOutput{}, fmt.Errorf("verification requires reviewing phase or a ready brain task")
	}
	leaseID, events, err := delegation.TaskRun(runsDir, task.ID)
	if err != nil {
		return RunVerificationOutput{}, err
	}
	var submission *delegation.Result
	for _, event := range events {
		if event.Result != nil {
			submission = event.Result
		}
	}
	if submission == nil || submission.Status != delegation.ResultCompleted {
		return RunVerificationOutput{}, fmt.Errorf("task has no completed Pinky submission")
	}
	limit := input.MaxBytes
	if limit == 0 {
		limit = verification.DefaultMaxBytes
	}
	record := delegation.Verification{Commands: []verification.Result{}, ScopeViolations: verification.ScopeViolations(submission.ChangedFiles, sectionBullets(task.Sections["Write Scope"]))}
	for _, argv := range task.Verification {
		result := verification.Run(ctx, root, verification.Command{Argv: argv, Workdir: input.Workdir, Timeout: time.Duration(input.TimeoutMS) * time.Millisecond, MaxBytes: limit})
		record.Commands = append(record.Commands, result)
		if result.ExitCode != 0 || result.TimedOut || result.Interrupted {
			record.Passed = false
			break
		}
		record.Passed = true
	}
	record.ContentIdentity, err = verification.ContentIdentity(root)
	if err != nil {
		return RunVerificationOutput{}, err
	}
	state := "verification_failed"
	recovery := []string{"fix the failure and rerun verification", "request changes", "replan the task"}
	if record.Passed {
		state = "verified"
		recovery = []string{"semantically review and accept the task"}
	}
	if len(record.ScopeViolations) > 0 {
		recovery = append(recovery, "keep the reported files explicitly", "revert them outside Pathframe", "replan the task scope")
	}
	if err := delegation.AppendRun(runsDir, leaseID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: leaseID, State: state, Timestamp: now(), Verification: &record}); err != nil {
		return RunVerificationOutput{}, err
	}
	return RunVerificationOutput{Schema: delegation.RunSchema, Change: view.Change, Task: task.ID, LeaseID: leaseID, Verification: record, Recovery: recovery}, nil
}
