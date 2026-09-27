package context

const (
	RoleSchema    = "pathframe.role/v1"
	PacketSchema  = "pathframe.task/v1"
	DefaultBudget = 64 * 1024
)

type Layer string

const (
	Foundation Layer = "foundation"
	Change     Layer = "change"
	Task       Layer = "task"
	Runtime    Layer = "runtime"
)

type Role struct {
	Schema         string   `json:"schema"`
	ID             string   `json:"id"`
	Mission        string   `json:"mission"`
	Reads          []string `json:"reads"`
	OptionalReads  []string `json:"optional_reads"`
	AllowedActions []string `json:"allowed_actions"`
	Returns        []string `json:"returns"`
}

type Reference struct {
	Layer    Layer
	Path     string
	Required bool
}

type Entry struct {
	Layer    Layer  `json:"layer"`
	Path     string `json:"path"`
	Required bool   `json:"required"`
	Content  string `json:"content"`
	Bytes    int    `json:"bytes"`
}

type Omission struct {
	Layer  Layer  `json:"layer"`
	Path   string `json:"path"`
	Reason string `json:"reason"`
	Bytes  int    `json:"bytes"`
}

type Budget struct {
	LimitBytes    int `json:"limit_bytes"`
	RequiredBytes int `json:"required_bytes"`
	IncludedBytes int `json:"included_bytes"`
	OmittedBytes  int `json:"omitted_bytes"`
}

type Packet struct {
	Schema                 string     `json:"schema"`
	Change                 string     `json:"change"`
	Task                   string     `json:"task"`
	Title                  string     `json:"title"`
	ExecutionPolicy        string     `json:"execution_policy"`
	Role                   Role       `json:"role"`
	Objective              string     `json:"objective"`
	Acceptance             string     `json:"acceptance"`
	References             []string   `json:"references"`
	Context                []Entry    `json:"context"`
	Dependencies           []string   `json:"dependencies"`
	CompletedPrerequisites []string   `json:"completed_prerequisites"`
	WriteScope             []string   `json:"write_scope"`
	WriteScopeAssurance    string     `json:"write_scope_assurance"`
	Constraints            []string   `json:"constraints"`
	Verification           [][]string `json:"verification"`
	ResultSchema           string     `json:"result_schema"`
	HostAssurance          string     `json:"host_assurance"`
	Omissions              []Omission `json:"omissions"`
	Budget                 Budget     `json:"budget"`
}
