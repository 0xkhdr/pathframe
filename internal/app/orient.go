package app

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

type Service struct {
	Dir string
}

func (s Service) Orient(change string) (workflow.Result, error) {
	start := s.Dir
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return workflow.Result{}, err
		}
	}
	root, err := store.Discover(start)
	if errors.Is(err, store.ErrProjectNotFound) {
		root, _ = filepath.Abs(start)
		return workflow.Result{
			Schema: workflow.Schema, ProjectRoot: root, Alternatives: []workflow.ActionResult{}, Diagnostics: []workflow.ReasonCode{}, Blockers: []workflow.Blocker{{Code: workflow.ReasonNoProject, Message: ".pathframe project not found", Recovery: []workflow.ActionResult{}}}, HumanRequired: true,
		}, nil
	}
	if err != nil {
		return workflow.Result{}, err
	}

	changes, err := store.Changes(root)
	if err != nil {
		return workflow.Result{}, err
	}
	if change == "" {
		switch len(changes) {
		case 0:
			return workflow.Result{Schema: workflow.Schema, ProjectRoot: root, Configured: true, Alternatives: []workflow.ActionResult{}, Diagnostics: []workflow.ReasonCode{}, Blockers: []workflow.Blocker{{Code: workflow.ReasonNoChange, Message: "no active change found", Recovery: []workflow.ActionResult{}}}, HumanRequired: true}, nil
		case 1:
			change = changes[0]
		default:
			return workflow.Result{Schema: workflow.Schema, ProjectRoot: root, Configured: true, Alternatives: []workflow.ActionResult{}, Diagnostics: []workflow.ReasonCode{}, Blockers: []workflow.Blocker{{Code: workflow.ReasonAmbiguousChange, Message: "multiple changes found; select one with --change", Recovery: []workflow.ActionResult{{Action: workflow.ActionSelect}}}}, HumanRequired: true}, nil
		}
	}
	dir, err := store.ChangeDir(root, change)
	if err != nil {
		return workflow.Result{}, err
	}
	replay, err := store.ReplayAndRepair(dir)
	if err != nil {
		return workflow.Result{}, err
	}
	if _, statErr := os.Stat(filepath.Join(dir, "change.yaml")); errors.Is(statErr, os.ErrNotExist) {
		view := result(root, replay.State, replay.Diagnostics)
		view.Blockers = []workflow.Blocker{{Code: workflow.ReasonArtifactsMissing, Message: "planning artifacts are not initialized", Recovery: []workflow.ActionResult{{Action: workflow.ActionNew}}}}
		view.Recommended = workflow.ActionResult{Action: workflow.ActionNew}
		view.HumanRequired = true
		return view, nil
	} else if statErr != nil {
		return workflow.Result{}, statErr
	}
	if replay.State.PlanIdentity != "" && !replay.State.Phase.Terminal() {
		plan, issues := artifacts.Validate(dir)
		if len(issues) > 0 || plan.Identity != replay.State.PlanIdentity {
			before := replay.State
			after, applyErr := workflow.Apply(before, workflow.ActionReplan, workflow.ActorSystem)
			if applyErr != nil {
				return workflow.Result{}, applyErr
			}
			event, eventErr := store.NewEvent(before, after, workflow.ActionReplan, workflow.ActorSystem, "material plan change invalidated approval")
			if eventErr != nil {
				return workflow.Result{}, eventErr
			}
			if appendErr := store.Append(filepath.Join(dir, "history.jsonl"), event); appendErr != nil {
				return workflow.Result{}, appendErr
			}
			replay, err = store.ReplayAndRepair(dir)
			if err != nil {
				return workflow.Result{}, err
			}
		}
	}
	view := result(root, replay.State, replay.Diagnostics)
	if plan, issues := artifacts.Validate(dir); len(issues) == 0 {
		completed, completedErr := delegation.CompletedTasks(filepath.Join(dir, "runs"))
		if completedErr != nil {
			return workflow.Result{}, completedErr
		}
		view.Progress = workflow.Progress{Completed: len(completed), Total: len(plan.Tasks)}
	}
	return view, nil
}

func result(root string, state workflow.State, diagnostics []workflow.ReasonCode) workflow.Result {
	actions := workflow.LegalActions(state)
	view := workflow.Result{
		Schema: workflow.Schema, ProjectRoot: root, Configured: true, Change: state.Change, Phase: state.Phase,
		Progress: workflow.Progress{Completed: state.Completed, Total: state.Total}, Alternatives: []workflow.ActionResult{}, Blockers: []workflow.Blocker{}, Diagnostics: diagnostics,
		HumanRequired: state.Phase == workflow.PhasePaused || state.Phase == workflow.PhasePlanning || state.Phase == workflow.PhaseReplanning,
	}
	if view.Diagnostics == nil {
		view.Diagnostics = []workflow.ReasonCode{}
	}
	if len(actions) > 0 {
		view.Recommended = workflow.ActionResult{Action: actions[0]}
		for _, action := range actions[1:] {
			view.Alternatives = append(view.Alternatives, workflow.ActionResult{Action: action})
		}
	}
	return view
}
