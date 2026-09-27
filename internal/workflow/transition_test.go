package workflow

import "testing"

func TestRecoveryFromEveryNonTerminalPhase(t *testing.T) {
	for _, phase := range Phases {
		if phase.Terminal() {
			continue
		}
		state := State{Change: "change", Phase: phase}
		if phase == PhasePaused {
			state.ResumePhase = PhaseExecuting
		}
		for _, action := range LegalActions(state) {
			if _, err := Apply(state, action, ActorHuman); err == nil {
				goto recovered
			}
		}
		t.Fatalf("phase %s has no executable human recovery", phase)
	recovered:
	}
}

func TestTransitionMatrixAndActors(t *testing.T) {
	tests := []struct {
		phase  Phase
		action Action
		actor  Actor
		want   Phase
		ok     bool
	}{
		{PhaseExploring, ActionPlan, ActorBrain, PhasePlanning, true},
		{PhasePlanning, ActionApprove, ActorHuman, PhaseReady, true},
		{PhaseReady, ActionExecute, ActorBrain, PhaseExecuting, true},
		{PhaseExecuting, ActionReview, ActorPinky, PhaseReviewing, true},
		{PhaseExecuting, ActionBlock, ActorSystem, PhaseBlocked, true},
		{PhaseReviewing, ActionComplete, ActorBrain, PhaseDone, true},
		{PhaseBlocked, ActionExecute, ActorBrain, PhaseExecuting, true},
		{PhaseReplanning, ActionPlan, ActorHuman, PhasePlanning, true},
		{PhaseReady, ActionExecute, ActorHuman, PhaseReady, false},
		{PhasePlanning, ActionReview, ActorBrain, PhasePlanning, false},
		{PhaseDone, ActionReplan, ActorHuman, PhaseDone, false},
	}
	for _, test := range tests {
		got, err := Apply(State{Change: "change", Phase: test.phase}, test.action, test.actor)
		if (err == nil) != test.ok || got.Phase != test.want {
			t.Errorf("Apply(%s, %s, %s) = (%s, %v), want (%s, ok=%t)", test.phase, test.action, test.actor, got.Phase, err, test.want, test.ok)
		}
		if err != nil {
			transitionErr, ok := err.(*TransitionError)
			if !ok || transitionErr.Code == "" || (!test.phase.Terminal() && len(transitionErr.Recovery) == 0) {
				t.Errorf("rejection lacks state, code, or recovery: %#v", err)
			}
		}
	}
}

func TestPauseResumeRoundTrip(t *testing.T) {
	start := State{Change: "change", Phase: PhaseExecuting}
	paused, err := Apply(start, ActionPause, ActorHuman)
	if err != nil {
		t.Fatal(err)
	}
	resumed, err := Apply(paused, ActionResume, ActorHuman)
	if err != nil {
		t.Fatal(err)
	}
	if resumed != start {
		t.Fatalf("resumed = %#v, want %#v", resumed, start)
	}
}

func TestAllDeclaredStatesAreValid(t *testing.T) {
	for _, phase := range Phases {
		if !ValidPhase(phase) {
			t.Fatalf("invalid declared phase %q", phase)
		}
	}
	for _, state := range TaskStates {
		if !ValidTaskState(state) {
			t.Fatalf("invalid declared task state %q", state)
		}
		for _, action := range []TaskAction{TaskPrepare, TaskActivate, TaskSubmit, TaskComplete, TaskBlock, TaskRequestChanges, TaskRetry} {
			for _, actor := range []Actor{ActorHuman, ActorBrain, ActorPinky, ActorSystem} {
				next, ok := ApplyTask(state, action, actor)
				want := containsTaskAction(LegalTaskActions(state), action) && taskActors[action][actor]
				if ok != want {
					t.Fatalf("task legality differs for %s/%s/%s", state, action, actor)
				}
				if !ok && next != state {
					t.Fatalf("illegal task transition changed %s to %s", state, next)
				}
			}
		}
	}
}

func containsTaskAction(actions []TaskAction, want TaskAction) bool {
	for _, action := range actions {
		if action == want {
			return true
		}
	}
	return false
}
