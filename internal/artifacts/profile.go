package artifacts

import "fmt"

var artifactSections = map[string][]string{
	"intent":       {"Summary", "Outcomes", "Non-goals", "Questions"},
	"requirements": {"Requirements", "Acceptance", "Questions"},
	"design":       {"Approach", "Decisions", "Questions"},
	"risks":        {"Risks", "Mitigations", "Questions"},
	"rollout":      {"Steps", "Rollback", "Questions"},
	"recovery":     {"Failures", "Recovery", "Questions"},
	"task":         {"Objective", "Required Reads", "Write Scope", "Constraints", "Verification", "Acceptance", "Questions", "Assumptions"},
}

var artifactSchemas = map[string]string{
	"intent": "pathframe.intent/v1", "requirements": "pathframe.requirements/v1",
	"design": "pathframe.design/v1", "risks": "pathframe.risks/v1",
	"rollout": "pathframe.rollout/v1", "recovery": "pathframe.recovery/v1",
	"task": "pathframe.task/v1",
}

func RequiredArtifacts(mode Mode) []string {
	result := []string{"intent"}
	if mode.rank() >= Standard.rank() {
		result = append(result, "requirements", "design")
	}
	if mode == HighRisk {
		result = append(result, "risks", "rollout", "recovery")
	}
	return result
}

func ArtifactInstructions(mode Mode, kind string) (Instructions, error) {
	sections, ok := artifactSections[kind]
	if !mode.Valid() || !ok {
		return Instructions{}, fmt.Errorf("unknown mode or artifact: %s/%s", mode, kind)
	}
	fields := []string{"schema", "profile"}
	if kind == "task" {
		fields = append(fields, "id", "title", "execution_policy", "role", "references", "dependencies")
	}
	return Instructions{
		Schema: InstructionsSchema, Profile: Profile, Mode: mode, Artifact: kind,
		RequiredFields: fields, RequiredSections: append([]string(nil), sections...),
		Rules: []string{"replace every angle-bracket placeholder", "use 'none' for resolved questions", "keep stable IDs and references", "write verification as JSON argv arrays; no shell strings"},
	}, nil
}
