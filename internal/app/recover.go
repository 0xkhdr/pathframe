package app

import "github.com/0xkhdr/pathframe/internal/workflow"

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
	return s.Transition(change, workflow.ActionCancel)
}
