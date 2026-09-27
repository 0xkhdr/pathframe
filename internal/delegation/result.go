package delegation

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

const ResultSchema = "pathframe.task-result/v1"
const MaxResultBytes = 256 * 1024

type Status string

const (
	ResultCompleted   Status = "completed"
	ResultFailed      Status = "failed"
	ResultBlocked     Status = "blocked"
	ResultNeedsReplan Status = "needs_replan"
)

type VerificationReport struct {
	Argv    []string `json:"argv"`
	Outcome string   `json:"outcome"`
	Summary string   `json:"summary,omitempty"`
}

type Result struct {
	Schema       string               `json:"schema"`
	LeaseID      string               `json:"lease_id"`
	Change       string               `json:"change"`
	Task         string               `json:"task"`
	Status       Status               `json:"status"`
	Summary      string               `json:"summary"`
	ChangedFiles []string             `json:"changed_files"`
	Verification []VerificationReport `json:"verification"`
	Discoveries  []string             `json:"discoveries"`
	Questions    []string             `json:"questions"`
	Risks        []string             `json:"risks"`
}

func ValidateResult(result Result, lease Lease) error {
	data, err := json.Marshal(result)
	if err != nil || len(data) > MaxResultBytes {
		return fmt.Errorf("result exceeds the %d byte run-record limit", MaxResultBytes)
	}
	if result.Schema != ResultSchema || result.LeaseID != lease.ID || result.Change != lease.Change || result.Task != lease.Task || strings.TrimSpace(result.Summary) == "" {
		return fmt.Errorf("result identity, schema, or summary does not match the active lease")
	}
	switch result.Status {
	case ResultCompleted, ResultFailed, ResultBlocked, ResultNeedsReplan:
	default:
		return fmt.Errorf("unsupported result status %q", result.Status)
	}
	for _, name := range result.ChangedFiles {
		clean := filepath.Clean(name)
		if name == "" || strings.Contains(name, `\`) || filepath.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return fmt.Errorf("changed file escapes the project: %q", name)
		}
	}
	return nil
}
