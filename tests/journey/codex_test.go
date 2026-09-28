package journey_test

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCodexPlanningJourney(t *testing.T) {
	binary := buildPathframe(t)
	root := t.TempDir()
	if got := run(t, binary, root, "codex-install", "--session-orientation", "--json"); !containsAll(got, `"schema": "pathframe.integration/v1"`, `"version": "1.0.0"`) {
		t.Fatalf("codex install = %s", got)
	}
	if got := run(t, binary, root, "codex-doctor", "--json"); !containsAll(got, `"healthy": true`, `"diagnostics": []`) {
		t.Fatalf("codex doctor = %s", got)
	}
	ctx := context.Background()
	session := connectCodex(t, ctx, binary, root)

	assessment := callTool(t, ctx, session, "pathframe_assess_request", map[string]any{"explicit_pathframe": true})
	if assessment["activation"] != "must_use" {
		t.Fatalf("assessment = %v", assessment)
	}
	callTool(t, ctx, session, "pathframe_create_change", map[string]any{"change": "demo", "mode": "quick"})
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "quick"), filepath.Join(root, ".pathframe", "changes", "demo"))
	validated := callTool(t, ctx, session, "pathframe_validate_plan", map[string]any{"change": "demo", "human_approved": false})
	if validated["approval_required"] != true {
		t.Fatalf("validation = %v", validated)
	}
	approved := callTool(t, ctx, session, "pathframe_validate_plan", map[string]any{"change": "demo", "human_approved": true})
	if approved["approval"] == nil {
		t.Fatalf("approval = %v", approved)
	}
	doctor := callTool(t, ctx, session, "pathframe_doctor", map[string]any{"change": "demo"})
	if doctor["healthy"] != true {
		t.Fatalf("MCP doctor = %v", doctor)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}

	resumed := connectCodex(t, ctx, binary, root)
	defer resumed.Close()
	orientation := callTool(t, ctx, resumed, "pathframe_orient", map[string]any{"change": "demo"})
	if orientation["phase"] != "ready" || orientation["schema"] != "pathframe.workflow/v1" {
		t.Fatalf("fresh-session orientation = %v", orientation)
	}
}

func containsAll(value string, values ...string) bool {
	for _, candidate := range values {
		if !strings.Contains(value, candidate) {
			return false
		}
	}
	return true
}

func connectCodex(t *testing.T, ctx context.Context, binary, root string) *mcp.ClientSession {
	t.Helper()
	command := exec.Command(binary, "mcp")
	command.Dir = root
	client := mcp.NewClient(&mcp.Implementation{Name: "codex-journey", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: command}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callTool(t *testing.T, ctx context.Context, session *mcp.ClientSession, name string, arguments map[string]any) map[string]any {
	t.Helper()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("%s failed: %v", name, result.Content)
	}
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var output map[string]any
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatal(err)
	}
	return output
}
