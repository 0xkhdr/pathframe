package app

import (
	"fmt"
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func (s Service) GetNext(change string) (workflow.Result, error) { return s.Orient(change) }

func (s Service) Transition(change string, action workflow.Action) (workflow.Result, error) {
	current, err := s.Orient(change)
	if err != nil {
		return workflow.Result{}, err
	}
	if current.Change == "" {
		return current, fmt.Errorf("cannot %s without a selected change", action)
	}
	before := workflow.State{Change: current.Change, Phase: current.Phase, Completed: current.Progress.Completed, Total: current.Progress.Total}
	// Resume target is persisted in the projection and reconstructed by replay.
	root, err := store.Discover(s.Dir)
	if err != nil {
		return workflow.Result{}, err
	}
	dir, err := store.ChangeDir(root, current.Change)
	if err != nil {
		return workflow.Result{}, err
	}
	replayed, err := store.ReplayAndRepair(dir)
	if err != nil {
		return workflow.Result{}, err
	}
	before = replayed.State
	after, err := workflow.Apply(before, action, workflow.ActorHuman)
	if transitionErr, ok := err.(*workflow.TransitionError); ok {
		current.Blockers = []workflow.Blocker{{Code: transitionErr.Code, Message: transitionErr.Message, Recovery: actionResults(transitionErr.Recovery)}}
		current.HumanRequired = true
		return current, err
	}
	if err != nil {
		return workflow.Result{}, err
	}
	event, err := store.NewEvent(before, after, action, workflow.ActorHuman, "requested through canonical navigation operation")
	if err != nil {
		return workflow.Result{}, err
	}
	if err := store.Append(filepath.Join(dir, "history.jsonl"), event); err != nil {
		return workflow.Result{}, err
	}
	replayed, err = store.RefreshProjection(dir)
	if err != nil {
		return workflow.Result{}, err
	}
	return result(root, replayed.State, replayed.Diagnostics), nil
}

func actionResults(actions []workflow.Action) []workflow.ActionResult {
	results := make([]workflow.ActionResult, len(actions))
	for i, action := range actions {
		results[i] = workflow.ActionResult{Action: action}
	}
	return results
}
