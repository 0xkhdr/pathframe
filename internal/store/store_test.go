package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestDiscoverAndRejectEscapes(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".pathframe", "changes", "good", "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := Discover(filepath.Join(root, ".pathframe", "changes", "good", "nested"))
	if err != nil || got != root {
		t.Fatalf("Discover() = %q, %v", got, err)
	}
	if _, err := ChangeDir(root, "../escape"); err == nil {
		t.Fatal("ChangeDir accepted path escape")
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".pathframe", "changes", "link")); err != nil {
		t.Fatal(err)
	}
	if _, err := ChangeDir(root, "link"); err == nil {
		t.Fatal("ChangeDir accepted symlink escape")
	}
}

func TestReplayRepairsProjectionAndCrashTail(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "change")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	journal := filepath.Join(dir, "history.jsonl")
	initial := Event{Schema: EventSchema, ID: "one", Timestamp: time.Unix(1, 0).UTC(), Actor: workflow.ActorSystem, Action: workflow.ActionInitialize, Change: "change", Target: workflow.PhaseExploring}
	if err := Append(journal, initial); err != nil {
		t.Fatal(err)
	}
	before := workflow.State{Change: "change", Phase: workflow.PhaseExploring}
	after, _ := workflow.Apply(before, workflow.ActionPlan, workflow.ActorBrain)
	event := Event{Schema: EventSchema, ID: "two", Timestamp: time.Unix(2, 0).UTC(), Actor: workflow.ActorBrain, Action: workflow.ActionPlan, Change: "change", Source: before.Phase, Target: after.Phase}
	if err := Append(journal, event); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "state.json"), []byte("bad"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(journal, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString(`{"schema":"pathframe.transition/v1"`); err != nil {
		t.Fatal(err)
	}
	file.Close()

	replay, err := ReplayAndRepair(dir)
	if err != nil {
		t.Fatal(err)
	}
	if replay.State.Phase != workflow.PhasePlanning || replay.Events != 2 {
		t.Fatalf("replay = %#v", replay)
	}
	if len(replay.Diagnostics) != 2 {
		t.Fatalf("diagnostics = %v, want corrupt projection and incomplete tail", replay.Diagnostics)
	}
	if projection, err := LoadProjection(filepath.Join(dir, "state.json")); err != nil || projection.State != replay.State {
		t.Fatalf("projection = %#v, %v", projection, err)
	}
	data, _ := os.ReadFile(journal)
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatal("incomplete crash tail was not safely removed")
	}
}
