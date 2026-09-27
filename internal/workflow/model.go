package workflow

const Schema = "pathframe.workflow/v1"

type Phase string

const (
	PhaseExploring  Phase = "exploring"
	PhasePlanning   Phase = "planning"
	PhaseReady      Phase = "ready"
	PhaseExecuting  Phase = "executing"
	PhaseReviewing  Phase = "reviewing"
	PhaseBlocked    Phase = "blocked"
	PhaseReplanning Phase = "replanning"
	PhaseDone       Phase = "done"
	PhasePaused     Phase = "paused"
	PhaseCancelled  Phase = "cancelled"
)

var TaskStates = []TaskState{TaskPending, TaskReady, TaskActive, TaskSubmitted, TaskCompleted, TaskBlocked, TaskChangesRequested}

func ValidTaskState(state TaskState) bool {
	for _, candidate := range TaskStates {
		if state == candidate {
			return true
		}
	}
	return false
}

var Phases = []Phase{
	PhaseExploring, PhasePlanning, PhaseReady, PhaseExecuting, PhaseReviewing,
	PhaseBlocked, PhaseReplanning, PhaseDone, PhasePaused, PhaseCancelled,
}

func (p Phase) Terminal() bool { return p == PhaseDone || p == PhaseCancelled }

func ValidPhase(p Phase) bool {
	for _, candidate := range Phases {
		if p == candidate {
			return true
		}
	}
	return false
}

type TaskState string

const (
	TaskPending          TaskState = "pending"
	TaskReady            TaskState = "ready"
	TaskActive           TaskState = "active"
	TaskSubmitted        TaskState = "submitted"
	TaskCompleted        TaskState = "completed"
	TaskBlocked          TaskState = "blocked"
	TaskChangesRequested TaskState = "changes_requested"
)

type Actor string

const (
	ActorHuman  Actor = "human"
	ActorBrain  Actor = "brain"
	ActorPinky  Actor = "pinky"
	ActorSystem Actor = "system"
)

type Action string

const (
	ActionInitialize Action = "initialize"
	ActionOrient     Action = "orient"
	ActionPlan       Action = "plan"
	ActionApprove    Action = "approve"
	ActionExecute    Action = "execute"
	ActionReview     Action = "review"
	ActionComplete   Action = "complete"
	ActionBlock      Action = "block"
	ActionPause      Action = "pause"
	ActionResume     Action = "resume"
	ActionReplan     Action = "replan"
	ActionCancel     Action = "cancel"
	ActionSelect     Action = "select_change"
	ActionNew        Action = "new"
)

type ReasonCode string

const (
	ReasonIllegalTransition ReasonCode = "illegal_transition"
	ReasonActorNotAllowed   ReasonCode = "actor_not_allowed"
	ReasonTerminalState     ReasonCode = "terminal_state"
	ReasonNotPaused         ReasonCode = "not_paused"
	ReasonNoProject         ReasonCode = "project_not_found"
	ReasonNoChange          ReasonCode = "change_not_found"
	ReasonAmbiguousChange   ReasonCode = "ambiguous_change"
	ReasonCorruptProjection ReasonCode = "corrupt_projection"
	ReasonIncompleteTail    ReasonCode = "incomplete_journal_tail"
	ReasonArtifactsMissing  ReasonCode = "artifacts_missing"
)

type State struct {
	Change       string `json:"change"`
	Phase        Phase  `json:"phase"`
	ResumePhase  Phase  `json:"resume_phase,omitempty"`
	PlanIdentity string `json:"plan_identity,omitempty"`
	PlanMode     string `json:"plan_mode,omitempty"`
	Completed    int    `json:"completed"`
	Total        int    `json:"total"`
}

type ActionResult struct {
	Action Action `json:"action"`
}

type Blocker struct {
	Code     ReasonCode     `json:"code"`
	Message  string         `json:"message"`
	Recovery []ActionResult `json:"recovery"`
}

type Progress struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
}

type Result struct {
	Schema        string         `json:"schema"`
	ProjectRoot   string         `json:"project_root"`
	Configured    bool           `json:"configured"`
	Change        string         `json:"change,omitempty"`
	Phase         Phase          `json:"phase,omitempty"`
	Progress      Progress       `json:"progress"`
	Recommended   ActionResult   `json:"recommended"`
	Alternatives  []ActionResult `json:"alternatives"`
	Blockers      []Blocker      `json:"blockers"`
	Diagnostics   []ReasonCode   `json:"diagnostics"`
	HumanRequired bool           `json:"human_required"`
}
