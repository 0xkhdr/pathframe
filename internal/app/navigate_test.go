package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestOperationsPersistAndRejectWithRecovery(t *testing.T) {
	root := project(t, "demo")
	service := Service{Dir: filepath.Join(root, ".pathframe", "changes", "demo")}
	paused, err := service.Pause("")
	if err != nil || paused.Phase != workflow.PhasePaused || len(paused.Diagnostics) != 0 {
		t.Fatalf("Pause() = %#v, %v", paused, err)
	}
	resumed, err := service.Resume("")
	if err != nil || resumed.Phase != workflow.PhaseExploring {
		t.Fatalf("Resume() = %#v, %v", resumed, err)
	}
	cancelled, err := service.Cancel("")
	if err != nil || cancelled.Phase != workflow.PhaseCancelled {
		t.Fatalf("Cancel() = %#v, %v", cancelled, err)
	}
	rejected, err := service.Pause("")
	if err == nil || rejected.Phase != workflow.PhaseCancelled || len(rejected.Blockers) != 1 || rejected.Blockers[0].Code != workflow.ReasonTerminalState {
		t.Fatalf("terminal Pause() = %#v, %v", rejected, err)
	}
}

func TestOrientRequiresSelectionForMultipleChanges(t *testing.T) {
	root := project(t, "one")
	makeChange(t, root, "two")
	service := Service{Dir: root}
	result, err := service.Orient("")
	if err != nil || len(result.Blockers) != 1 || result.Blockers[0].Code != workflow.ReasonAmbiguousChange || result.Blockers[0].Recovery[0].Action != workflow.ActionSelect {
		t.Fatalf("Orient() = %#v, %v", result, err)
	}
	selected, err := service.Orient("two")
	if err != nil || selected.Change != "two" {
		t.Fatalf("Orient(two) = %#v, %v", selected, err)
	}
}

func TestOrientMissingSelectedChangeReturnsStructuredRecovery(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".pathframe"), 0o755); err != nil {
		t.Fatal(err)
	}
	result, err := (Service{Dir: root}).Orient("demo")
	if err != nil || len(result.Blockers) != 1 || result.Blockers[0].Code != workflow.ReasonNoChange || result.Recommended.Action != workflow.ActionNew {
		t.Fatalf("Orient(missing) = %#v, %v", result, err)
	}
}

func project(t *testing.T, change string) string {
	t.Helper()
	root := t.TempDir()
	makeChange(t, root, change)
	return root
}

func makeChange(t *testing.T, root, change string) {
	t.Helper()
	dir := filepath.Join(root, ".pathframe", "changes", change)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	event := store.Event{Schema: store.EventSchema, ID: change, Timestamp: time.Unix(1, 0).UTC(), Actor: workflow.ActorSystem, Action: workflow.ActionInitialize, Change: change, Target: workflow.PhaseExploring}
	if err := store.Append(filepath.Join(dir, "history.jsonl"), event); err != nil {
		t.Fatal(err)
	}
}
