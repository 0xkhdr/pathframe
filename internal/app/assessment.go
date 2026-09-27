package app

const AssessmentSchema = "pathframe.request-assessment/v1"

type RequestAssessmentInput struct {
	ExplicitPathframe bool   `json:"explicit_pathframe,omitempty"`
	ActiveWorkflow    bool   `json:"active_workflow,omitempty"`
	ContinuesWorkflow bool   `json:"continues_workflow,omitempty"`
	DirectExecution   bool   `json:"direct_execution,omitempty"`
	Kind              string `json:"kind,omitempty"`
	DependentSteps    bool   `json:"dependent_steps,omitempty"`
	NeedsPersistence  bool   `json:"needs_persistence,omitempty"`
	NeedsAgreement    bool   `json:"needs_agreement,omitempty"`
	MeaningfulRisk    bool   `json:"meaningful_risk,omitempty"`
}

type RequestAssessment struct {
	Schema        string `json:"schema"`
	Activation    string `json:"activation"`
	Reason        string `json:"reason"`
	HumanRequired bool   `json:"human_required"`
}

// AssessRequest applies the published activation rules to facts classified by Brain.
func (Service) AssessRequest(in RequestAssessmentInput) RequestAssessment {
	result := RequestAssessment{Schema: AssessmentSchema}
	if in.ExplicitPathframe || in.ActiveWorkflow || in.ContinuesWorkflow {
		result.Activation, result.Reason = "must_use", "Pathframe was requested or its workflow is already active"
		return result
	}
	if in.DirectExecution || in.Kind == "explanation" || in.Kind == "read_only" || in.Kind == "brainstorm" || in.Kind == "trivial_edit" || in.Kind == "non_development" {
		result.Activation, result.Reason = "must_not_use", "the request is excluded from automatic activation"
		return result
	}
	if in.DependentSteps || in.NeedsPersistence || in.NeedsAgreement || in.MeaningfulRisk {
		result.Activation, result.Reason, result.HumanRequired = "offer", "the development request benefits from a durable agreed path", true
		return result
	}
	result.Activation, result.Reason = "must_not_use", "no Pathframe activation condition is present"
	return result
}
