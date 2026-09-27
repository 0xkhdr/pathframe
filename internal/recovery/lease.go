package recovery

import (
	"path/filepath"
	"time"

	"github.com/0xkhdr/pathframe/internal/delegation"
)

func DiagnoseLease(changeDir string, now time.Time, staleAfter time.Duration) ([]Diagnosis, error) {
	lease, active, err := delegation.ActiveLease(filepath.Join(changeDir, "runs"))
	if err != nil {
		return []Diagnosis{{Code: "lease_invalid", Severity: SeverityError, Subject: "runs/active.json", Evidence: err.Error(), HumanRequired: true}}, nil
	}
	if !active {
		return nil, nil
	}
	if staleAfter > 0 && now.Sub(lease.CreatedAt) >= staleAfter {
		return []Diagnosis{{Code: "abandoned_lease", Severity: SeverityWarning, Subject: lease.Task, Evidence: "active lease has exceeded its recovery window", Repair: "release_abandoned_lease"}}, nil
	}
	return []Diagnosis{{Code: "active_lease", Severity: SeverityInfo, Subject: lease.Task, Evidence: "work is actively leased", HumanRequired: true}}, nil
}
