package workflow

type TaskAction string

const (
	TaskPrepare        TaskAction = "prepare"
	TaskActivate       TaskAction = "activate"
	TaskSubmit         TaskAction = "submit"
	TaskComplete       TaskAction = "complete"
	TaskBlock          TaskAction = "block"
	TaskRequestChanges TaskAction = "request_changes"
	TaskRetry          TaskAction = "retry"
)

var taskTransitions = map[TaskState]map[TaskAction]TaskState{
	TaskPending:          {TaskPrepare: TaskReady},
	TaskReady:            {TaskActivate: TaskActive},
	TaskActive:           {TaskSubmit: TaskSubmitted, TaskBlock: TaskBlocked},
	TaskSubmitted:        {TaskComplete: TaskCompleted, TaskBlock: TaskBlocked, TaskRequestChanges: TaskChangesRequested},
	TaskBlocked:          {TaskRetry: TaskReady},
	TaskChangesRequested: {TaskRetry: TaskReady},
}

var taskActors = map[TaskAction]map[Actor]bool{
	TaskPrepare:        {ActorBrain: true, ActorSystem: true},
	TaskActivate:       {ActorBrain: true},
	TaskSubmit:         {ActorBrain: true, ActorPinky: true},
	TaskComplete:       {ActorBrain: true, ActorHuman: true},
	TaskBlock:          {ActorBrain: true, ActorPinky: true, ActorSystem: true},
	TaskRequestChanges: {ActorBrain: true},
	TaskRetry:          {ActorBrain: true},
}

func ApplyTask(state TaskState, action TaskAction, actor Actor) (TaskState, bool) {
	if !taskActors[action][actor] {
		return state, false
	}
	next, ok := taskTransitions[state][action]
	if !ok {
		return state, false
	}
	return next, true
}

func LegalTaskActions(state TaskState) []TaskAction {
	var actions []TaskAction
	for _, action := range []TaskAction{TaskPrepare, TaskActivate, TaskSubmit, TaskComplete, TaskBlock, TaskRequestChanges, TaskRetry} {
		if _, ok := taskTransitions[state][action]; ok {
			actions = append(actions, action)
		}
	}
	return actions
}
