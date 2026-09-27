package app

import (
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func (s Service) Pause(change string) (workflow.Result, error) {
	return s.Transition(change, workflow.ActionPause)
}

func (s Service) Resume(change string) (workflow.Result, error) {
	return s.Transition(change, workflow.ActionResume)
}

func (s Service) Replan(change string) (workflow.Result, error) {
	return s.Transition(change, workflow.ActionReplan)
}

func (s Service) Cancel(change string) (workflow.Result, error) {
	result, err := s.Transition(change, workflow.ActionCancel)
	if err != nil || result.Change == "" {
		return result, err
	}
	root, err := store.Discover(s.Dir)
	if err != nil {
		return result, err
	}
	dir, err := store.ChangeDir(root, result.Change)
	if err != nil {
		return result, err
	}
	runsDir := filepath.Join(dir, "runs")
	lease, active, err := delegation.ActiveLease(runsDir)
	if err != nil || !active {
		return result, err
	}
	if err := delegation.AppendRun(runsDir, lease.ID, delegation.RunEvent{Schema: delegation.RunSchema, LeaseID: lease.ID, State: "cancelled", Timestamp: now()}); err != nil {
		return result, err
	}
	return result, delegation.Release(runsDir, lease)
}
