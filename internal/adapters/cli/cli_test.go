package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "orientation", want: "Pathframe\n\nConfigured: false\nActive change: none\nPhase: none\nProgress: 0/0 tasks complete\nBlocked: .pathframe project not found (project_not_found)\nRecommended next action: none\nHuman action required: true\n"},
		{name: "help", args: []string{"--help"}, want: "Usage: pathframe [--json] [--change ID] [status|next|new|template|check|approve|packet|submit-result|verify|accept-task|request-changes|lease-release|edit-check|pause|resume|replan|cancel|codex-install|codex-doctor|claude-install|claude-doctor]\n"},
		{name: "version", args: []string{"--version"}, want: "dev\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := Run(test.args, &stdout, &stderr); code != 0 {
				t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
			}
			if !strings.HasPrefix(stdout.String(), test.want) {
				t.Fatalf("Run() stdout = %q, want prefix %q", stdout.String(), test.want)
			}
			if stderr.Len() != 0 {
				t.Fatalf("Run() stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestTextAndJSONGoldenShareCanonicalResult(t *testing.T) {
	fixtureRoot, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata", "workflow"))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	changeDir := filepath.Join(root, ".pathframe", "changes", "demo")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	event := store.Event{Schema: store.EventSchema, ID: "initial", Timestamp: time.Unix(1, 0).UTC(), Actor: workflow.ActorSystem, Action: workflow.ActionInitialize, Change: "demo", Target: workflow.PhaseExploring}
	if err := store.Append(filepath.Join(changeDir, "history.jsonl"), event); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	var textOut, jsonOut, stderr bytes.Buffer
	if code := Run([]string{"status"}, &textOut, &stderr); code != 0 {
		t.Fatalf("text code = %d, stderr = %s", code, stderr.String())
	}
	stderr.Reset()
	if code := Run([]string{"status", "--json"}, &jsonOut, &stderr); code != 0 {
		t.Fatalf("json code = %d, stderr = %s", code, stderr.String())
	}
	wantText, err := os.ReadFile(filepath.Join(fixtureRoot, "status.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if textOut.String() != string(wantText) {
		t.Fatalf("text = %q, want %q", textOut.String(), wantText)
	}
	var got map[string]any
	if err := json.Unmarshal(jsonOut.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	got["project_root"] = "ROOT"
	normalized, _ := json.MarshalIndent(got, "", "  ")
	normalized = append(normalized, '\n')
	wantJSON, err := os.ReadFile(filepath.Join(fixtureRoot, "status.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(normalized, wantJSON) {
		t.Fatalf("json = %s, want %s", normalized, wantJSON)
	}
	if !strings.Contains(textOut.String(), got["phase"].(string)) || !strings.Contains(textOut.String(), got["recommended"].(map[string]any)["action"].(string)) {
		t.Fatal("text does not project canonical JSON phase/action")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"start"}, &stdout, &stderr); code != 2 {
		t.Fatalf("Run() code = %d, want 2", code)
	}
	if got := stderr.String(); got != "pathframe: unknown command \"start\"\n" {
		t.Fatalf("Run() stderr = %q", got)
	}
}
