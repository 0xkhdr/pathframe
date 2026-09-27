package journey_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type event struct {
	Schema      string `json:"schema"`
	ID          string `json:"id"`
	Timestamp   string `json:"timestamp"`
	Actor       string `json:"actor"`
	Action      string `json:"action"`
	Change      string `json:"change"`
	Source      string `json:"source,omitempty"`
	Target      string `json:"target"`
	ResumePhase string `json:"resume_phase,omitempty"`
}

func TestNavigatorRecoveryJourneys(t *testing.T) {
	binary := buildPathframe(t)
	chains := map[string][]event{
		"exploring":  nil,
		"planning":   {{Actor: "brain", Action: "plan", Source: "exploring", Target: "planning"}},
		"ready":      {{Actor: "brain", Action: "plan", Source: "exploring", Target: "planning"}, {Actor: "human", Action: "approve", Source: "planning", Target: "ready"}},
		"executing":  {{Actor: "brain", Action: "plan", Source: "exploring", Target: "planning"}, {Actor: "human", Action: "approve", Source: "planning", Target: "ready"}, {Actor: "brain", Action: "execute", Source: "ready", Target: "executing"}},
		"reviewing":  {{Actor: "brain", Action: "plan", Source: "exploring", Target: "planning"}, {Actor: "human", Action: "approve", Source: "planning", Target: "ready"}, {Actor: "brain", Action: "execute", Source: "ready", Target: "executing"}, {Actor: "brain", Action: "review", Source: "executing", Target: "reviewing"}},
		"blocked":    {{Actor: "brain", Action: "plan", Source: "exploring", Target: "planning"}, {Actor: "human", Action: "approve", Source: "planning", Target: "ready"}, {Actor: "brain", Action: "execute", Source: "ready", Target: "executing"}, {Actor: "system", Action: "block", Source: "executing", Target: "blocked"}},
		"replanning": {{Actor: "human", Action: "replan", Source: "exploring", Target: "replanning"}},
		"paused":     {{Actor: "human", Action: "pause", Source: "exploring", Target: "paused", ResumePhase: "exploring"}},
	}
	for phase, chain := range chains {
		t.Run(phase, func(t *testing.T) {
			root := makeProject(t, chain)
			command, want := "pause", `"phase": "paused"`
			if phase == "paused" {
				command, want = "resume", `"phase": "exploring"`
			}
			output := run(t, binary, root, "--change", "demo", command, "--json")
			if !strings.Contains(output, want) {
				t.Fatalf("%s output = %s", command, output)
			}
		})
	}
}

func TestNavigatorPauseResumeReplanCancelAndAmbiguity(t *testing.T) {
	binary := buildPathframe(t)
	root := makeProject(t, nil)
	if got := run(t, binary, root); !strings.Contains(got, "Phase: exploring") {
		t.Fatalf("orientation = %s", got)
	}
	run(t, binary, root, "pause")
	run(t, binary, root, "resume")
	run(t, binary, root, "replan")
	if got := run(t, binary, root, "cancel", "--json"); !strings.Contains(got, `"phase": "cancelled"`) {
		t.Fatalf("cancel = %s", got)
	}

	second := filepath.Join(root, ".pathframe", "changes", "other")
	if err := os.MkdirAll(second, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJournal(t, filepath.Join(second, "history.jsonl"), "other", nil)
	command := exec.Command(binary, "status", "--json")
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "ambiguous_change") || !strings.Contains(string(output), "select_change") {
		t.Fatalf("ambiguous status = %s, %v", output, err)
	}
}

func buildPathframe(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "pathframe")
	command := exec.Command("go", "build", "-o", binary, "./cmd/pathframe")
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build: %s: %v", output, err)
	}
	return binary
}

func makeProject(t *testing.T, chain []event) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJournal(t, filepath.Join(dir, "history.jsonl"), "demo", chain)
	return root
}

func writeJournal(t *testing.T, path, change string, chain []event) {
	t.Helper()
	events := append([]event{{Actor: "system", Action: "initialize", Change: change, Target: "exploring"}}, chain...)
	var output bytes.Buffer
	for i := range events {
		events[i].Schema = "pathframe.transition/v1"
		events[i].ID = change + "-" + string(rune('a'+i))
		events[i].Timestamp = "2026-01-01T00:00:00Z"
		events[i].Change = change
		data, err := json.Marshal(events[i])
		if err != nil {
			t.Fatal(err)
		}
		output.Write(data)
		output.WriteByte('\n')
	}
	if err := os.WriteFile(path, output.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, binary, root string, args ...string) string {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("pathframe %v: %s: %v", args, output, err)
	}
	return string(output)
}
