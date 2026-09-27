package delegation

import (
	"path/filepath"
	"testing"
	"time"

	pfcontext "github.com/0xkhdr/pathframe/internal/context"
)

func TestPreflightLeaseResultAndNoFallback(t *testing.T) {
	packet := pfcontext.Packet{ExecutionPolicy: "delegated", Role: pfcontext.Role{ID: "backend"}, WriteScope: []string{"internal/x.go"}, Verification: [][]string{{"go", "test", "./..."}}}
	capabilities := CapabilityDeclaration{Host: "codex", Capabilities: []string{CapabilitySequentialSubagent, CapabilitySharedWorkspace, CapabilityResultReturn}}
	if got := Check(packet, capabilities, false); !got.Ready || got.WriteScopeAssurance != "advisory" {
		t.Fatalf("preflight = %#v", got)
	}
	runs := t.TempDir()
	lease, err := Acquire(runs, "demo", "T1", "codex", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Acquire(runs, "demo", "T2", "codex", time.Now()); err == nil {
		t.Fatal("conflicting lease succeeded")
	}
	result := Result{Schema: ResultSchema, LeaseID: lease.ID, Change: "demo", Task: "T1", Status: ResultCompleted, Summary: "done", ChangedFiles: []string{filepath.Join("internal", "x.go")}, Verification: []VerificationReport{}, Discoveries: []string{}, Questions: []string{}, Risks: []string{}}
	if err := ValidateResult(result, lease); err != nil {
		t.Fatal(err)
	}
	if got := Reconcile(result.Status); got.TaskState != "submitted" || got.BrainMayEdit {
		t.Fatalf("reconcile = %#v", got)
	}
	if BrainEditAllowed("delegated", false) {
		t.Fatal("delegation failure authorized Brain fallback")
	}
}

func TestPreflightRejectsMissingCapabilityAndEnforcementMustBeDeclared(t *testing.T) {
	packet := pfcontext.Packet{ExecutionPolicy: "delegated", Role: pfcontext.Role{ID: "backend"}, WriteScope: []string{"x"}, Verification: [][]string{{"test"}}}
	got := Check(packet, CapabilityDeclaration{Host: "codex", Capabilities: []string{CapabilitySequentialSubagent}}, false)
	if got.Ready || got.WriteScopeAssurance != "advisory" || len(got.Issues) != 2 {
		t.Fatalf("preflight = %#v", got)
	}
	capabilities := CapabilityDeclaration{Host: "codex", Capabilities: []string{CapabilitySequentialSubagent, CapabilitySharedWorkspace, CapabilityResultReturn, CapabilityWriteScopeEnforcement}}
	if got := Check(packet, capabilities, false); !got.Ready || got.WriteScopeAssurance != "host_enforced" {
		t.Fatalf("enforced = %#v", got)
	}
}

func TestAllWorkerStatusesAndMalformedResult(t *testing.T) {
	want := map[Status]string{ResultCompleted: "submitted", ResultFailed: "blocked", ResultBlocked: "blocked", ResultNeedsReplan: "blocked"}
	for status, state := range want {
		if got := Reconcile(status); got.TaskState != state || (status == ResultNeedsReplan) != got.RequiresReplan {
			t.Fatalf("%s reconciliation = %#v", status, got)
		}
	}
	lease := Lease{ID: "lease", Change: "demo", Task: "T1"}
	bad := Result{Schema: ResultSchema, LeaseID: "wrong", Change: "demo", Task: "T1", Status: ResultCompleted, Summary: "done"}
	if err := ValidateResult(bad, lease); err == nil {
		t.Fatal("malformed result was accepted")
	}
	bad.LeaseID, bad.ChangedFiles = "lease", []string{"../escape"}
	if err := ValidateResult(bad, lease); err == nil {
		t.Fatal("escaping changed file was accepted")
	}
}
