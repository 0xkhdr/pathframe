package app

import "testing"

func TestAssessRequestActivationRules(t *testing.T) {
	tests := []struct {
		name string
		in   RequestAssessmentInput
		want string
	}{
		{"explicit", RequestAssessmentInput{ExplicitPathframe: true, DirectExecution: true}, "must_use"},
		{"active", RequestAssessmentInput{ActiveWorkflow: true}, "must_use"},
		{"offer", RequestAssessmentInput{DependentSteps: true}, "offer"},
		{"explanation", RequestAssessmentInput{Kind: "explanation", MeaningfulRisk: true}, "must_not_use"},
		{"direct", RequestAssessmentInput{DirectExecution: true, DependentSteps: true}, "must_not_use"},
		{"ordinary", RequestAssessmentInput{Kind: "development"}, "must_not_use"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := (Service{}).AssessRequest(test.in); got.Activation != test.want {
				t.Fatalf("activation = %q, want %q", got.Activation, test.want)
			}
		})
	}
}
