package delegation

type Reconciliation struct {
	Status         Status   `json:"status"`
	TaskState      string   `json:"task_state"`
	BrainMayEdit   bool     `json:"brain_may_edit"`
	RequiresReplan bool     `json:"requires_replan"`
	Recovery       []string `json:"recovery"`
}

func Reconcile(status Status) Reconciliation {
	result := Reconciliation{Status: status, BrainMayEdit: false, Recovery: []string{}}
	switch status {
	case ResultCompleted:
		result.TaskState = "submitted"
		result.Recovery = []string{"review the submission; Stage 7 will add verification and acceptance"}
	case ResultNeedsReplan:
		result.TaskState, result.RequiresReplan = "blocked", true
		result.Recovery = []string{"replan the task and obtain human reapproval"}
	default:
		result.TaskState = "blocked"
		result.Recovery = []string{"repair the blocker and explicitly retry delegation", "replan the task"}
	}
	return result
}

func BrainEditAllowed(policy string, activeLease bool) bool { return policy == "brain" && !activeLease }
