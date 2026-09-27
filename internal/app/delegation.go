package app

import (
	"fmt"
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

type SubmitResultInput struct {
	Result delegation.Result `json:"result"`
}

type SubmitResultOutput struct {
	Schema         string                    `json:"schema"`
	Change         string                    `json:"change"`
	Task           string                    `json:"task"`
	LeaseID        string                    `json:"lease_id"`
	Reconciliation delegation.Reconciliation `json:"reconciliation"`
	Phase          workflow.Phase            `json:"phase"`
}

type ReleaseLeaseInput struct {
	Change  string `json:"change"`
	LeaseID string `json:"lease_id"`
}
type ReleaseLeaseOutput struct {
	Change   string         `json:"change"`
	Task     string         `json:"task"`
	LeaseID  string         `json:"lease_id"`
	Phase    workflow.Phase `json:"phase"`
	Recovery []string       `json:"recovery"`
}

type EditGuardInput struct {
	Change string `json:"change"`
	Task   string `json:"task"`
}
type EditGuardOutput struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

func (s Service) CheckBrainEdit(input EditGuardInput) (EditGuardOutput, error) {
	allowed, reason, err := s.BrainEditGuard(input.Change, input.Task)
	return EditGuardOutput{Allowed: allowed, Reason: reason}, err
}

func (s Service) ReleaseLease(input ReleaseLeaseInput) (ReleaseLeaseOutput, error) {
	root, err := store.Discover(s.Dir)
	if err != nil {
		return ReleaseLeaseOutput{}, err
	}
	dir, err := store.ChangeDir(root, input.Change)
	if err != nil {
		return ReleaseLeaseOutput{}, err
	}
	runsDir := filepath.Join(dir, "runs")
	lease, active, err := delegation.ActiveLease(runsDir)
	if err != nil {
		return ReleaseLeaseOutput{}, err
	}
	if !active || lease.ID != input.LeaseID {
		return ReleaseLeaseOutput{}, fmt.Errorf("active lease does not match release request")
	}
	if err := delegation.AppendRun(runsDir, lease.ID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: lease.ID, State: "released", Timestamp: now()}); err != nil {
		return ReleaseLeaseOutput{}, err
	}
	if err := s.appendDelegationTransition(dir, workflow.ActionBlock, workflow.ActorSystem, "lost delegation lease explicitly released"); err != nil {
		return ReleaseLeaseOutput{}, err
	}
	if err := delegation.Release(runsDir, lease); err != nil {
		return ReleaseLeaseOutput{}, err
	}
	return ReleaseLeaseOutput{Change: lease.Change, Task: lease.Task, LeaseID: lease.ID, Phase: workflow.PhaseBlocked, Recovery: []string{"explicitly retry delegation", "replan the task"}}, nil
}

func (s Service) SubmitResult(input SubmitResultInput) (SubmitResultOutput, error) {
	root, err := store.Discover(s.Dir)
	if err != nil {
		return SubmitResultOutput{}, err
	}
	dir, err := store.ChangeDir(root, input.Result.Change)
	if err != nil {
		return SubmitResultOutput{}, err
	}
	runsDir := filepath.Join(dir, "runs")
	lease, active, err := delegation.ActiveLease(runsDir)
	if err != nil {
		return SubmitResultOutput{}, err
	}
	if !active {
		return SubmitResultOutput{}, fmt.Errorf("no active lease exists for result submission")
	}
	if err := delegation.ValidateResult(input.Result, lease); err != nil {
		return SubmitResultOutput{}, err
	}
	reconciliation := delegation.Reconcile(input.Result.Status)
	if err := delegation.AppendRun(runsDir, lease.ID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: lease.ID, State: reconciliation.TaskState, Timestamp: now(), Result: &input.Result}); err != nil {
		return SubmitResultOutput{}, err
	}
	action := workflow.ActionBlock
	actor := workflow.ActorBrain
	if input.Result.Status == delegation.ResultCompleted {
		action = workflow.ActionReview
	}
	if err := s.appendDelegationTransition(dir, action, actor, "Pinky result "+string(input.Result.Status)+" for lease "+lease.ID); err != nil {
		return SubmitResultOutput{}, err
	}
	if err := delegation.Release(runsDir, lease); err != nil {
		return SubmitResultOutput{}, err
	}
	view, err := s.Orient(input.Result.Change)
	if err != nil {
		return SubmitResultOutput{}, err
	}
	return SubmitResultOutput{Schema: delegation.ResultSchema, Change: lease.Change, Task: lease.Task, LeaseID: lease.ID, Reconciliation: reconciliation, Phase: view.Phase}, nil
}

func (s Service) appendDelegationTransition(changeDir string, action workflow.Action, actor workflow.Actor, reason string) error {
	replay, err := store.ReplayAndRepair(changeDir)
	if err != nil {
		return err
	}
	after, err := workflow.Apply(replay.State, action, actor)
	if err != nil {
		return err
	}
	event, err := store.NewEvent(replay.State, after, action, actor, reason)
	if err != nil {
		return err
	}
	if err := store.Append(filepath.Join(changeDir, "history.jsonl"), event); err != nil {
		return err
	}
	_, err = store.ReplayAndRepair(changeDir)
	return err
}

func (s Service) BrainEditGuard(change, task string) (bool, string, error) {
	root, err := store.Discover(s.Dir)
	if err != nil {
		return false, "", err
	}
	dir, err := store.ChangeDir(root, change)
	if err != nil {
		return false, "", err
	}
	plan, issues := artifacts.Validate(dir)
	if len(issues) > 0 {
		return false, "approved plan is invalid", nil
	}
	policy := ""
	for _, candidate := range plan.Tasks {
		if candidate.ID == task {
			policy = candidate.ExecutionPolicy
			break
		}
	}
	if policy == "" {
		return false, "task not found", nil
	}
	_, active, err := delegation.ActiveLease(filepath.Join(dir, "runs"))
	if err != nil {
		return false, "", err
	}
	allowed := delegation.BrainEditAllowed(policy, active)
	if !allowed {
		return false, "Brain cannot edit a delegated task; delegation failure does not change execution_policy", nil
	}
	return true, "Brain execution is explicitly authorized", nil
}
