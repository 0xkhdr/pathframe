package recovery

import "sort"

const Schema = "pathframe.diagnosis/v1"

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Diagnosis struct {
	Code          string   `json:"code"`
	Severity      Severity `json:"severity"`
	Subject       string   `json:"subject"`
	Evidence      string   `json:"evidence"`
	Repair        string   `json:"repair,omitempty"`
	HumanRequired bool     `json:"human_required"`
}

type Report struct {
	Schema      string      `json:"schema"`
	ProjectRoot string      `json:"project_root,omitempty"`
	Change      string      `json:"change,omitempty"`
	Healthy     bool        `json:"healthy"`
	Repaired    bool        `json:"repaired"`
	Diagnoses   []Diagnosis `json:"diagnoses"`
}

func New(root, change string, diagnoses []Diagnosis) Report {
	sort.Slice(diagnoses, func(i, j int) bool {
		if diagnoses[i].Code == diagnoses[j].Code {
			return diagnoses[i].Subject < diagnoses[j].Subject
		}
		return diagnoses[i].Code < diagnoses[j].Code
	})
	healthy := true
	for _, diagnosis := range diagnoses {
		healthy = healthy && diagnosis.Severity == SeverityInfo
	}
	return Report{Schema: Schema, ProjectRoot: root, Change: change, Healthy: healthy, Diagnoses: diagnoses}
}
