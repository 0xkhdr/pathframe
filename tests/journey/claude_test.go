package journey_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestClaudePlanningJourney(t *testing.T) {
	binary := buildPathframe(t)
	root := t.TempDir()
	if got := run(t, binary, root, "claude-install", "--session-orientation", "--json"); !containsAll(got, `"schema": "pathframe.integration/v1"`, `"version": "1.0.0"`) {
		t.Fatalf("claude install = %s", got)
	}
	if got := run(t, binary, root, "claude-doctor", "--json"); !containsAll(got, `"healthy": true`, `"diagnostics": []`) {
		t.Fatalf("claude doctor = %s", got)
	}
	settings, err := os.ReadFile(filepath.Join(root, ".claude", "settings.json"))
	if err != nil || !containsAll(string(settings), `"SessionStart"`, `status --json`) {
		t.Fatalf("session hook = %s, %v", settings, err)
	}
	ctx := context.Background()
	session := connectCodex(t, ctx, binary, root)
	if got := callTool(t, ctx, session, "pathframe_assess_request", map[string]any{"active_workflow": true}); got["activation"] != "must_use" {
		t.Fatalf("assessment = %v", got)
	}
	callTool(t, ctx, session, "pathframe_create_change", map[string]any{"change": "demo", "mode": "quick"})
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "quick"), filepath.Join(root, ".pathframe", "changes", "demo"))
	if got := callTool(t, ctx, session, "pathframe_validate_plan", map[string]any{"change": "demo", "human_approved": false}); got["approval_required"] != true {
		t.Fatalf("validation = %v", got)
	}
	if got := callTool(t, ctx, session, "pathframe_validate_plan", map[string]any{"change": "demo", "human_approved": true}); got["approval"] == nil {
		t.Fatalf("approval = %v", got)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	resumed := connectCodex(t, ctx, binary, root)
	defer resumed.Close()
	if got := callTool(t, ctx, resumed, "pathframe_orient", map[string]any{"change": "demo"}); got["phase"] != "ready" || got["schema"] != "pathframe.workflow/v1" {
		t.Fatalf("fresh-session orientation = %v", got)
	}
	if err := os.WriteFile(filepath.Join(root, ".mcp.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := runFailure(t, binary, root, "claude-doctor", "--json"); !containsAll(got, `"healthy": false`, `"code": "missing_or_modified_asset"`, `pathframe claude-install`) {
		t.Fatalf("broken config diagnosis = %s", got)
	}
}

func runFailure(t *testing.T, binary, root string, args ...string) string {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("pathframe %v unexpectedly succeeded: %s", args, output)
	}
	return string(output)
}
