package recovery

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestProjectionDiagnosisAndRepairPreservesJournal(t *testing.T) {
	dir := t.TempDir()
	history := filepath.Join(dir, "history.jsonl")
	event := store.Event{Schema: store.EventSchema, ID: "initial", Timestamp: time.Unix(1, 0).UTC(), Actor: workflow.ActorSystem, Action: workflow.ActionInitialize, Change: "demo", Target: workflow.PhaseExploring}
	if err := store.Append(history, event); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(history)
	if err := os.WriteFile(filepath.Join(dir, "state.json"), []byte("broken"), 0o600); err != nil {
		t.Fatal(err)
	}
	diagnoses, err := DiagnoseProjection(dir)
	if err != nil || len(diagnoses) != 1 || diagnoses[0].Repair != RepairProjectionCode {
		t.Fatalf("diagnoses = %#v, %v", diagnoses, err)
	}
	if err := RepairProjection(dir); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(history)
	if string(after) != string(before) {
		t.Fatal("repair changed complete journal history")
	}
	if diagnoses, err = DiagnoseProjection(dir); err != nil || len(diagnoses) != 0 {
		t.Fatalf("after repair = %#v, %v", diagnoses, err)
	}
}

func TestIncompleteTailRepairKeepsCompleteEvents(t *testing.T) {
	dir := t.TempDir()
	history := filepath.Join(dir, "history.jsonl")
	event := store.Event{Schema: store.EventSchema, ID: "initial", Timestamp: time.Unix(1, 0).UTC(), Actor: workflow.ActorSystem, Action: workflow.ActionInitialize, Change: "demo", Target: workflow.PhaseExploring}
	if err := store.Append(history, event); err != nil {
		t.Fatal(err)
	}
	complete, _ := os.ReadFile(history)
	file, err := os.OpenFile(history, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = file.WriteString(`{"schema":"pathframe.transition/v1"`); err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
	diagnoses, err := DiagnoseProjection(dir)
	if err != nil || len(diagnoses) < 1 || diagnoses[0].Repair != RepairJournalTailCode {
		t.Fatalf("diagnoses = %#v, %v", diagnoses, err)
	}
	if err := RepairProjection(dir); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(history)
	if string(after) != string(complete) {
		t.Fatalf("history = %q, want %q", after, complete)
	}
}
