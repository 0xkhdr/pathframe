package journey_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPacketContextFrontierJourney(t *testing.T) {
	binary := buildPathframe(t)
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	root := t.TempDir()
	run(t, binary, root, "new", "--change", "demo", "--mode", "quick")
	target := filepath.Join(root, ".pathframe", "changes", "demo")
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "quick"), target)
	run(t, binary, root, "approve", "--change", "demo")
	jsonPreview := run(t, binary, root, "packet", "--change", "demo", "--json")
	if !containsAll(jsonPreview, `"schema": "pathframe.task/v1"`, `"write_scope_assurance": "advisory"`, `"layer": "change"`, `"layer": "task"`, `"layer": "runtime"`, `"frontier": [`, `"T1"`) {
		t.Fatalf("packet JSON = %s", jsonPreview)
	}
	textPreview := run(t, binary, root, "packet", "--change", "demo")
	if !containsAll(textPreview, "Task packet pathframe.task/v1", "Objective:\nImplement the described behavior.", "Acceptance:\n- The focused check passes.", "Context budget:", "Write scope assurance: advisory", "Frontier: T1", "Sequential next task: T1", "Wave 1: T1") {
		t.Fatalf("packet text = %s", textPreview)
	}
	data, err := os.ReadFile(filepath.Join(target, "tasks", "T1.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(jsonPreview, string(data)) {
		t.Fatal("packet accumulated the raw task document")
	}
	ctx := context.Background()
	session := connectCodex(t, ctx, binary, root)
	defer session.Close()
	mcpPreview := callTool(t, ctx, session, "pathframe_prepare_delegation", map[string]any{"change": "demo", "task": "T1"})
	packet, ok := mcpPreview["packet"].(map[string]any)
	if !ok || packet["schema"] != "pathframe.task/v1" || packet["task"] != "T1" {
		t.Fatalf("MCP packet = %#v", mcpPreview)
	}
}
