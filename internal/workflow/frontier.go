package workflow

func LegalActions(state State) []Action {
	if state.Phase.Terminal() {
		return nil
	}
	if state.Phase == PhasePaused {
		return []Action{ActionResume, ActionReplan, ActionCancel}
	}
	actions := make([]Action, 0, 5)
	for _, action := range []Action{ActionPlan, ActionApprove, ActionExecute, ActionReview, ActionComplete, ActionBlock} {
		if _, ok := normalTransitions[state.Phase][action]; ok {
			actions = append(actions, action)
		}
	}
	return append(actions, ActionPause, ActionReplan, ActionCancel)
}

func Recommended(state State) Action {
	actions := LegalActions(state)
	if len(actions) == 0 {
		return ""
	}
	return actions[0]
}
