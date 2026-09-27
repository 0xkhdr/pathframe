package workflow

import "fmt"

type TransitionError struct {
	State    State
	Code     ReasonCode
	Message  string
	Recovery []Action
}

func (e *TransitionError) Error() string { return e.Message }

var normalTransitions = map[Phase]map[Action]Phase{
	PhaseExploring:  {ActionPlan: PhasePlanning},
	PhasePlanning:   {ActionApprove: PhaseReady},
	PhaseReady:      {ActionExecute: PhaseExecuting},
	PhaseExecuting:  {ActionReview: PhaseReviewing, ActionBlock: PhaseBlocked},
	PhaseReviewing:  {ActionComplete: PhaseDone, ActionBlock: PhaseBlocked},
	PhaseBlocked:    {ActionExecute: PhaseExecuting},
	PhaseReplanning: {ActionPlan: PhasePlanning},
}

var allowedActors = map[Action]map[Actor]bool{
	ActionInitialize: {ActorSystem: true},
	ActionPlan:       {ActorHuman: true, ActorBrain: true},
	ActionApprove:    {ActorHuman: true},
	ActionExecute:    {ActorBrain: true},
	ActionReview:     {ActorBrain: true, ActorPinky: true},
	ActionComplete:   {ActorHuman: true, ActorBrain: true},
	ActionBlock:      {ActorBrain: true, ActorPinky: true, ActorSystem: true},
	ActionPause:      {ActorHuman: true},
	ActionResume:     {ActorHuman: true},
	ActionReplan:     {ActorHuman: true, ActorBrain: true},
	ActionCancel:     {ActorHuman: true},
}

func Apply(state State, action Action, actor Actor) (State, error) {
	if !allowedActors[action][actor] {
		return state, rejected(state, ReasonActorNotAllowed, fmt.Sprintf("%s cannot request %s", actor, action))
	}
	if state.Phase.Terminal() {
		return state, rejected(state, ReasonTerminalState, fmt.Sprintf("%s is terminal", state.Phase))
	}

	next := state
	switch action {
	case ActionPause:
		if state.Phase == PhasePaused {
			return state, rejected(state, ReasonIllegalTransition, "change is already paused")
		}
		next.ResumePhase, next.Phase = state.Phase, PhasePaused
	case ActionResume:
		if state.Phase != PhasePaused || !ValidPhase(state.ResumePhase) || state.ResumePhase.Terminal() || state.ResumePhase == PhasePaused {
			return state, rejected(state, ReasonNotPaused, "change has no resumable paused phase")
		}
		next.Phase, next.ResumePhase = state.ResumePhase, ""
	case ActionReplan:
		if state.Phase == PhaseReplanning {
			return state, rejected(state, ReasonIllegalTransition, "change is already replanning")
		}
		next.Phase, next.ResumePhase = PhaseReplanning, ""
	case ActionCancel:
		next.Phase, next.ResumePhase = PhaseCancelled, ""
	default:
		target, ok := normalTransitions[state.Phase][action]
		if !ok {
			return state, rejected(state, ReasonIllegalTransition, fmt.Sprintf("%s is not legal from %s", action, state.Phase))
		}
		next.Phase = target
	}
	return next, nil
}

func rejected(state State, code ReasonCode, message string) *TransitionError {
	recovery := LegalActions(state)
	if len(recovery) == 0 {
		recovery = []Action{ActionOrient}
	}
	return &TransitionError{State: state, Code: code, Message: message, Recovery: recovery}
}
