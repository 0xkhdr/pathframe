package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestAcceptanceRejectsStaleContentAndRequiresScopeDecision(t *testing.T) {
	for _, test := range []struct {
		name    string
		changed []string
		mutate  bool
		keep    bool
		wantErr string
	}{
		{name: "stale", changed: []string{"internal/example.go"}, mutate: true, wantErr: "changed after verification"},
		{name: "scope", changed: []string{"README.md"}, wantErr: "scope violations"},
		{name: "scope kept", changed: []string{"README.md"}, keep: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			service, root := submittedService(t, test.changed, false)
			verified, err := service.RunVerification(context.Background(), RunVerificationInput{Change: "demo", Task: "T1", TimeoutMS: 1000})
			if err != nil || !verified.Verification.Passed {
				t.Fatalf("verify = %#v, %v", verified, err)
			}
			if test.mutate {
				if err := os.WriteFile(filepath.Join(root, "changed.txt"), []byte("after verification"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			accepted, err := service.AcceptTask(AcceptTaskInput{Change: "demo", Task: "T1", Reason: "semantically accepted", KeepScopeViolations: test.keep})
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("accept = %#v, %v", accepted, err)
				}
				return
			}
			if err != nil || accepted.Phase != workflow.PhaseDone {
				t.Fatalf("accept = %#v, %v", accepted, err)
			}
		})
	}
}

func TestAcceptedTaskAdvancesFrontier(t *testing.T) {
	service, _ := submittedService(t, []string{"internal/example.go"}, true)
	if _, err := service.RunVerification(context.Background(), RunVerificationInput{Change: "demo", Task: "T1", TimeoutMS: 1000}); err != nil {
		t.Fatal(err)
	}
	accepted, err := service.AcceptTask(AcceptTaskInput{Change: "demo", Task: "T1", Reason: "accepted"})
	if err != nil || accepted.Phase != workflow.PhaseReady || accepted.Projection.Next != "T2" || accepted.Progress.Completed != 1 || accepted.Progress.Total != 2 {
		t.Fatalf("accepted = %#v, %v", accepted, err)
	}
}

func TestDirectBrainTaskVerifiesAndCompletes(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "quick"), dir)
	taskPath := filepath.Join(dir, "tasks", "T1.md")
	data, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, []byte(strings.Replace(string(data), `["go", "test", "./..."]`, `["true"]`, 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Quick); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Approve("demo"); err != nil {
		t.Fatal(err)
	}
	verified, err := service.RunVerification(context.Background(), RunVerificationInput{Change: "demo", Task: "T1", ChangedFiles: []string{"internal/example.go"}, TimeoutMS: 1000})
	if err != nil || !verified.Verification.Passed {
		t.Fatalf("verified = %#v, %v", verified, err)
	}
	accepted, err := service.AcceptTask(AcceptTaskInput{Change: "demo", Task: "T1", Reason: "accepted"})
	if err != nil || accepted.Phase != workflow.PhaseDone {
		t.Fatalf("accepted = %#v, %v", accepted, err)
	}
}

func submittedService(t *testing.T, changed []string, second bool) (Service, string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "standard"), dir)
	taskPath := filepath.Join(dir, "tasks", "T1.md")
	data, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), `["go", "test", "./..."]`, `["true"]`, 1))
	if err := os.WriteFile(taskPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if second {
		task2 := strings.Replace(string(data), "id: T1", "id: T2", 1)
		task2 = strings.Replace(task2, "title: Implement the standard feature", "title: Follow-up", 1)
		task2 = strings.Replace(task2, "dependencies: []", `dependencies: ["T1"]`, 1)
		if err := os.WriteFile(filepath.Join(dir, "tasks", "T2.md"), []byte(task2), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Standard); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Approve("demo"); err != nil {
		t.Fatal(err)
	}
	lease, err := delegation.Acquire(filepath.Join(dir, "runs"), "demo", "T1", "codex", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := service.appendDelegationTransition(dir, workflow.ActionExecute, workflow.ActorBrain, "test lease"); err != nil {
		t.Fatal(err)
	}
	result := delegation.Result{Schema: delegation.ResultSchema, LeaseID: lease.ID, Change: "demo", Task: "T1", Status: delegation.ResultCompleted, Summary: "implemented", ChangedFiles: changed, Verification: []delegation.VerificationReport{}, Discoveries: []string{}, Questions: []string{}, Risks: []string{}}
	if _, err := service.SubmitResult(SubmitResultInput{Result: result}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Discover(root); err != nil {
		t.Fatal(err)
	}
	return service, root
}
