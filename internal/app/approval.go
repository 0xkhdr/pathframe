package app

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

var now = func() time.Time { return time.Now().UTC() }

func (s Service) Check(change string) (artifacts.CheckResult, error) {
	current, err := s.Orient(change)
	if err != nil {
		return artifacts.CheckResult{}, err
	}
	if current.Change == "" {
		return artifacts.CheckResult{}, fmt.Errorf("cannot check without a selected change")
	}
	root, err := store.Discover(s.Dir)
	if err != nil {
		return artifacts.CheckResult{}, err
	}
	dir, err := store.ChangeDir(root, current.Change)
	if err != nil {
		return artifacts.CheckResult{}, err
	}
	plan, issues := artifacts.Validate(dir)
	replay, replayErr := store.ReplayAndRepair(dir)
	if replayErr != nil {
		return artifacts.CheckResult{}, replayErr
	}
	if replay.State.PlanMode != "" && modeRank(plan.Change.Mode) < modeRank(artifacts.Mode(replay.State.PlanMode)) {
		issues = append(issues, artifacts.Issue{Path: "change.yaml", Code: "mode_lowering", Message: "mode cannot be lowered below the approved " + replay.State.PlanMode + " mode"})
		plan.Identity = ""
	}
	result := artifacts.CheckResult{Schema: artifacts.CheckSchema, Change: current.Change, Mode: plan.Change.Mode, Valid: len(issues) == 0, Identity: plan.Identity, Issues: issues, Recovery: []string{}}
	if len(issues) > 0 {
		result.Recovery = []string{"edit the named authored artifact", "run pathframe check again"}
	}
	return result, nil
}

func (s Service) Approve(change string) (artifacts.ApprovalResult, error) {
	check, err := s.Check(change)
	if err != nil {
		return artifacts.ApprovalResult{}, err
	}
	if !check.Valid {
		return artifacts.ApprovalResult{}, fmt.Errorf("plan validation failed; run check for precise fixes")
	}
	root, err := store.Discover(s.Dir)
	if err != nil {
		return artifacts.ApprovalResult{}, err
	}
	dir, err := store.ChangeDir(root, check.Change)
	if err != nil {
		return artifacts.ApprovalResult{}, err
	}
	replay, err := store.ReplayAndRepair(dir)
	if err != nil {
		return artifacts.ApprovalResult{}, err
	}
	if replay.State.Phase == workflow.PhaseReplanning {
		after, _ := workflow.Apply(replay.State, workflow.ActionPlan, workflow.ActorBrain)
		event, eventErr := store.NewEvent(replay.State, after, workflow.ActionPlan, workflow.ActorBrain, "revalidated material plan")
		if eventErr != nil {
			return artifacts.ApprovalResult{}, eventErr
		}
		if err := store.Append(filepath.Join(dir, "history.jsonl"), event); err != nil {
			return artifacts.ApprovalResult{}, err
		}
		replay, err = store.ReplayAndRepair(dir)
		if err != nil {
			return artifacts.ApprovalResult{}, err
		}
	}
	after, err := workflow.Apply(replay.State, workflow.ActionApprove, workflow.ActorHuman)
	if err != nil {
		return artifacts.ApprovalResult{}, err
	}
	event, err := store.NewEvent(replay.State, after, workflow.ActionApprove, workflow.ActorHuman, "explicit human approval of validated plan")
	if err != nil {
		return artifacts.ApprovalResult{}, err
	}
	event.PlanIdentity = check.Identity
	event.PlanMode = string(check.Mode)
	event.Timestamp = now()
	if err := store.Append(filepath.Join(dir, "history.jsonl"), event); err != nil {
		return artifacts.ApprovalResult{}, err
	}
	if _, err := store.ReplayAndRepair(dir); err != nil {
		return artifacts.ApprovalResult{}, err
	}
	return artifacts.ApprovalResult{Schema: artifacts.ApprovalSchema, Change: check.Change, Mode: check.Mode, Identity: check.Identity, Approved: event.Timestamp, Phase: string(workflow.PhaseReady)}, nil
}
