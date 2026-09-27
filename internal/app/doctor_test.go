package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestDoctorRepairsAbandonedLeaseWithoutPolicyFallback(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "standard"), dir)
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Standard); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Approve("demo"); err != nil {
		t.Fatal(err)
	}
	created := time.Unix(1, 0).UTC()
	lease, err := delegation.Acquire(filepath.Join(dir, "runs"), "demo", "T1", "codex", created)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.appendDelegationTransition(dir, workflow.ActionExecute, workflow.ActorBrain, "test abandoned lease"); err != nil {
		t.Fatal(err)
	}
	oldNow := now
	now = func() time.Time { return created.Add(25 * time.Hour) }
	t.Cleanup(func() { now = oldNow })
	report, err := service.Doctor(DoctorInput{Change: "demo", Repair: true})
	if err != nil || !report.Repaired {
		t.Fatalf("Doctor() = %#v, %v", report, err)
	}
	if _, active, err := delegation.ActiveLease(filepath.Join(dir, "runs")); err != nil || active {
		t.Fatalf("active = %t, %v", active, err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "runs", lease.ID+".jsonl"))
	if err != nil || !strings.Contains(string(data), `"state":"released"`) {
		t.Fatalf("run = %s, %v", data, err)
	}
	allowed, _, err := service.BrainEditGuard("demo", "T1")
	if err != nil || allowed {
		t.Fatalf("delegated policy changed after recovery: allowed=%t err=%v", allowed, err)
	}
}
