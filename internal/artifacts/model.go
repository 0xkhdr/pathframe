package artifacts

import "time"

const (
	Profile            = "okf-markdown/v1"
	InstructionsSchema = "pathframe.artifact-instructions/v1"
	ChangeSchema       = "pathframe.change/v1"
	CheckSchema        = "pathframe.plan-check/v1"
	ApprovalSchema     = "pathframe.plan-approval/v1"
)

type Mode string

const (
	Quick    Mode = "quick"
	Standard Mode = "standard"
	HighRisk Mode = "high-risk"
)

func (m Mode) Valid() bool { return m == Quick || m == Standard || m == HighRisk }

func (m Mode) rank() int {
	switch m {
	case Quick:
		return 1
	case Standard:
		return 2
	case HighRisk:
		return 3
	}
	return 0
}

type Change struct {
	Schema  string
	Profile string
	ID      string
	Mode    Mode
}

type Document struct {
	Path     string
	Fields   map[string]string
	Sections map[string]string
}

type Task struct {
	Document
	ID              string
	Title           string
	ExecutionPolicy string
	Role            string
	References      []string
	Dependencies    []string
	Verification    [][]string
}

type Issue struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Plan struct {
	Change    Change
	Documents []Document
	Tasks     []Task
	Identity  string
}

type CheckResult struct {
	Schema   string   `json:"schema"`
	Change   string   `json:"change"`
	Mode     Mode     `json:"mode"`
	Valid    bool     `json:"valid"`
	Identity string   `json:"identity,omitempty"`
	Issues   []Issue  `json:"issues"`
	Recovery []string `json:"recovery"`
}

type Instructions struct {
	Schema           string   `json:"schema"`
	Profile          string   `json:"profile"`
	Mode             Mode     `json:"mode"`
	Artifact         string   `json:"artifact"`
	RequiredFields   []string `json:"required_fields"`
	RequiredSections []string `json:"required_sections"`
	Rules            []string `json:"rules"`
}

type ApprovalResult struct {
	Schema   string    `json:"schema"`
	Change   string    `json:"change"`
	Mode     Mode      `json:"mode"`
	Identity string    `json:"identity"`
	Approved time.Time `json:"approved_at"`
	Phase    string    `json:"phase"`
}
