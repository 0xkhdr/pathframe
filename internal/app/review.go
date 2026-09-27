package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/verification"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

type AcceptTaskInput struct {
	Change              string `json:"change"`
	Task                string `json:"task"`
	Reason              string `json:"reason"`
	KeepScopeViolations bool   `json:"keep_scope_violations,omitempty"`
}

type ReviewOutput struct {
	Schema     string                `json:"schema"`
	Change     string                `json:"change"`
	Task       string                `json:"task"`
	Decision   string                `json:"decision"`
	Phase      workflow.Phase        `json:"phase"`
	Progress   workflow.Progress     `json:"progress"`
	Projection delegation.Projection `json:"projection"`
	Recovery   []string              `json:"recovery"`
}

type RequestChangesInput struct {
	Change string `json:"change"`
	Task   string `json:"task"`
	Reason string `json:"reason"`
}

func (s Service) AcceptTask(input AcceptTaskInput) (ReviewOutput, error) {
	if strings.TrimSpace(input.Reason) == "" {
		return ReviewOutput{}, fmt.Errorf("semantic acceptance requires a reason")
	}
	root, dir, plan, leaseID, events, err := s.reviewContext(input.Change, input.Task)
	if err != nil {
		return ReviewOutput{}, err
	}
	var latest *delegation.Verification
	for _, event := range events {
		if event.Verification != nil {
			copy := *event.Verification
			latest = &copy
		}
	}
	if latest == nil || !latest.Passed {
		return ReviewOutput{}, fmt.Errorf("task requires a passing Pathframe verification before semantic acceptance")
	}
	current, err := verification.ContentIdentity(root)
	if err != nil {
		return ReviewOutput{}, err
	}
	if current != latest.ContentIdentity {
		return ReviewOutput{}, fmt.Errorf("repository content changed after verification; rerun verification")
	}
	if len(latest.ScopeViolations) > 0 && !input.KeepScopeViolations {
		return ReviewOutput{}, fmt.Errorf("scope violations require an explicit keep decision, revert, or replan: %s", strings.Join(latest.ScopeViolations, ", "))
	}
	review := delegation.Review{Decision: "accepted", Reason: input.Reason, Recovery: []string{}}
	if err := delegation.AppendRun(filepath.Join(dir, "runs"), leaseID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: leaseID, State: "completed", Timestamp: now(), Review: &review}); err != nil {
		return ReviewOutput{}, err
	}
	completed, err := delegation.CompletedTasks(filepath.Join(dir, "runs"))
	if err != nil {
		return ReviewOutput{}, err
	}
	projection := delegation.Project(plan.Tasks, completed)
	action := workflow.ActionAdvance
	if len(completed) == len(plan.Tasks) {
		action = workflow.ActionComplete
	}
	if err := s.appendDelegationTransition(dir, action, workflow.ActorBrain, "Brain semantically accepted task "+input.Task); err != nil {
		return ReviewOutput{}, err
	}
	view, err := s.Orient(input.Change)
	if err != nil {
		return ReviewOutput{}, err
	}
	return ReviewOutput{Schema: delegation.RunSchema, Change: input.Change, Task: input.Task, Decision: "accepted", Phase: view.Phase, Progress: workflow.Progress{Completed: len(completed), Total: len(plan.Tasks)}, Projection: projection, Recovery: []string{}}, nil
}

func (s Service) RequestChanges(input RequestChangesInput) (ReviewOutput, error) {
	if strings.TrimSpace(input.Reason) == "" {
		return ReviewOutput{}, fmt.Errorf("requesting changes requires a reason")
	}
	_, dir, plan, leaseID, _, err := s.reviewContext(input.Change, input.Task)
	if err != nil {
		return ReviewOutput{}, err
	}
	review := delegation.Review{Decision: "changes_requested", Reason: input.Reason, Recovery: []string{"retry the same execution policy", "replan the task"}}
	if err := delegation.AppendRun(filepath.Join(dir, "runs"), leaseID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: leaseID, State: "changes_requested", Timestamp: now(), Review: &review}); err != nil {
		return ReviewOutput{}, err
	}
	if err := s.appendDelegationTransition(dir, workflow.ActionBlock, workflow.ActorBrain, "Brain requested changes for task "+input.Task); err != nil {
		return ReviewOutput{}, err
	}
	completed, err := delegation.CompletedTasks(filepath.Join(dir, "runs"))
	if err != nil {
		return ReviewOutput{}, err
	}
	return ReviewOutput{Schema: delegation.RunSchema, Change: input.Change, Task: input.Task, Decision: "changes_requested", Phase: workflow.PhaseBlocked, Progress: workflow.Progress{Completed: len(completed), Total: len(plan.Tasks)}, Projection: delegation.Project(plan.Tasks, completed), Recovery: review.Recovery}, nil
}

func (s Service) reviewContext(change, task string) (string, string, artifacts.Plan, string, []delegation.RunEvent, error) {
	view, err := s.Orient(change)
	if err != nil {
		return "", "", artifacts.Plan{}, "", nil, err
	}
	if view.Phase != workflow.PhaseReviewing {
		return "", "", artifacts.Plan{}, "", nil, fmt.Errorf("review decision requires reviewing phase")
	}
	root, err := store.Discover(s.Dir)
	if err != nil {
		return "", "", artifacts.Plan{}, "", nil, err
	}
	dir, err := store.ChangeDir(root, view.Change)
	if err != nil {
		return "", "", artifacts.Plan{}, "", nil, err
	}
	plan, issues := artifacts.Validate(dir)
	if len(issues) > 0 {
		return "", "", artifacts.Plan{}, "", nil, fmt.Errorf("approved plan is invalid")
	}
	found := false
	for _, candidate := range plan.Tasks {
		found = found || candidate.ID == task
	}
	if !found {
		return "", "", artifacts.Plan{}, "", nil, fmt.Errorf("task not found: %s", task)
	}
	leaseID, events, err := delegation.TaskRun(filepath.Join(dir, "runs"), task)
	return root, dir, plan, leaseID, events, err
}
