package journey_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDelegationVerificationCompletionJourneysCodexAndClaudeCode(t *testing.T) {
	for _, host := range []struct{ name, install string }{{"codex", "codex-install"}, {"claude-code", "claude-install"}} {
		t.Run(host.name, func(t *testing.T) {
			binary := buildPathframe(t)
			root := t.TempDir()
			run(t, binary, root, host.install)
			run(t, binary, root, "new", "--change", "demo", "--mode", "standard")
			repository, _ := filepath.Abs(filepath.Join("..", ".."))
			changeDir := filepath.Join(root, ".pathframe", "changes", "demo")
			copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "standard"), changeDir)
			taskPath := filepath.Join(changeDir, "tasks", "T1.md")
			taskData, err := os.ReadFile(taskPath)
			if err != nil {
				t.Fatal(err)
			}
			taskData = []byte(strings.Replace(string(taskData), `["go", "test", "./..."]`, fmt.Sprintf(`[%q, "--version"]`, binary), 1))
			if err := os.WriteFile(taskPath, taskData, 0o600); err != nil {
				t.Fatal(err)
			}
			roleDir := filepath.Join(root, ".pathframe", "roles")
			if err := os.MkdirAll(roleDir, 0o755); err != nil {
				t.Fatal(err)
			}
			role := "schema: pathframe.role/v1\nid: backend-engineer\nmission: Implement bounded work.\nreads: []\noptional_reads: []\nallowed_actions: [\"edit\",\"test\"]\nreturns: [\"summary\",\"changed_files\",\"verification\",\"discoveries\",\"questions\"]\n"
			if err := os.WriteFile(filepath.Join(roleDir, "backend-engineer.yaml"), []byte(role), 0o600); err != nil {
				t.Fatal(err)
			}
			run(t, binary, root, "approve", "--change", "demo")
			ctx := context.Background()
			session := connectCodex(t, ctx, binary, root)
			deployed := callTool(t, ctx, session, "pathframe_prepare_delegation", map[string]any{"change": "demo", "task": "T1", "host": host.name})
			lease, ok := deployed["lease"].(map[string]any)
			if !ok || lease["id"] == "" || deployed["preflight"].(map[string]any)["ready"] != true {
				t.Fatalf("prepare = %#v", deployed)
			}
			result := map[string]any{"schema": "pathframe.task-result/v1", "lease_id": lease["id"], "change": "demo", "task": "T1", "status": "completed", "summary": "implemented", "changed_files": []string{"internal/example.go"}, "verification": []any{}, "discoveries": []any{}, "questions": []any{}, "risks": []any{}}
			submitted := callTool(t, ctx, session, "pathframe_submit_result", map[string]any{"result": result})
			if submitted["phase"] != "reviewing" || submitted["reconciliation"].(map[string]any)["task_state"] != "submitted" {
				t.Fatalf("submit = %#v", submitted)
			}
			verified := callTool(t, ctx, session, "pathframe_run_verification", map[string]any{"change": "demo", "task": "T1", "workdir": ".", "timeout_ms": 5000})
			if verified["verification"].(map[string]any)["passed"] != true {
				t.Fatalf("verification = %#v", verified)
			}
			accepted := callTool(t, ctx, session, "pathframe_accept_task", map[string]any{"change": "demo", "task": "T1", "reason": "acceptance is satisfied"})
			if accepted["phase"] != "done" || accepted["decision"] != "accepted" {
				t.Fatalf("accept = %#v", accepted)
			}
			if err := session.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(changeDir, "runs", fmt.Sprint(lease["id"])+".jsonl")); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestChangesRequestedJourney(t *testing.T) {
	binary := buildPathframe(t)
	root := t.TempDir()
	run(t, binary, root, "codex-install")
	run(t, binary, root, "new", "--change", "demo", "--mode", "standard")
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	changeDir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "standard"), changeDir)
	roleDir := filepath.Join(root, ".pathframe", "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	role := "schema: pathframe.role/v1\nid: backend-engineer\nmission: Implement bounded work.\nreads: []\noptional_reads: []\nallowed_actions: [\"edit\",\"test\"]\nreturns: [\"summary\",\"changed_files\",\"verification\",\"discoveries\",\"questions\"]\n"
	if err := os.WriteFile(filepath.Join(roleDir, "backend-engineer.yaml"), []byte(role), 0o600); err != nil {
		t.Fatal(err)
	}
	run(t, binary, root, "approve", "--change", "demo")
	ctx := context.Background()
	session := connectCodex(t, ctx, binary, root)
	deployed := callTool(t, ctx, session, "pathframe_prepare_delegation", map[string]any{"change": "demo", "task": "T1", "host": "codex"})
	lease := deployed["lease"].(map[string]any)
	result := map[string]any{"schema": "pathframe.task-result/v1", "lease_id": lease["id"], "change": "demo", "task": "T1", "status": "completed", "summary": "implemented", "changed_files": []string{"internal/example.go"}, "verification": []any{}, "discoveries": []any{}, "questions": []any{}, "risks": []any{}}
	callTool(t, ctx, session, "pathframe_submit_result", map[string]any{"result": result})
	review := callTool(t, ctx, session, "pathframe_request_changes", map[string]any{"change": "demo", "task": "T1", "reason": "acceptance is not satisfied"})
	if review["phase"] != "blocked" || review["decision"] != "changes_requested" {
		t.Fatalf("request changes = %#v", review)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDelegationNeedsReplanGuardAndLostLeaseRecovery(t *testing.T) {
	binary := buildPathframe(t)
	root := t.TempDir()
	run(t, binary, root, "codex-install")
	run(t, binary, root, "new", "--change", "demo", "--mode", "standard")
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "standard"), filepath.Join(root, ".pathframe", "changes", "demo"))
	roleDir := filepath.Join(root, ".pathframe", "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	role := "schema: pathframe.role/v1\nid: backend-engineer\nmission: Implement bounded work.\nreads: []\noptional_reads: []\nallowed_actions: [\"edit\",\"test\"]\nreturns: [\"summary\",\"changed_files\",\"verification\",\"discoveries\",\"questions\"]\n"
	if err := os.WriteFile(filepath.Join(roleDir, "backend-engineer.yaml"), []byte(role), 0o600); err != nil {
		t.Fatal(err)
	}
	run(t, binary, root, "approve", "--change", "demo")
	ctx := context.Background()
	session := connectCodex(t, ctx, binary, root)
	deployed := callTool(t, ctx, session, "pathframe_prepare_delegation", map[string]any{"change": "demo", "task": "T1", "host": "codex"})
	lease := deployed["lease"].(map[string]any)
	guard := callTool(t, ctx, session, "pathframe_check_brain_edit", map[string]any{"change": "demo", "task": "T1"})
	if guard["allowed"] != false {
		t.Fatalf("guard = %#v", guard)
	}
	released := callTool(t, ctx, session, "pathframe_release_delegation", map[string]any{"change": "demo", "lease_id": lease["id"]})
	if released["phase"] != "blocked" {
		t.Fatalf("release = %#v", released)
	}
	retried := callTool(t, ctx, session, "pathframe_prepare_delegation", map[string]any{"change": "demo", "task": "T1", "host": "codex"})
	retryLease := retried["lease"].(map[string]any)
	result := map[string]any{"schema": "pathframe.task-result/v1", "lease_id": retryLease["id"], "change": "demo", "task": "T1", "status": "needs_replan", "summary": "plan discovery", "changed_files": []any{}, "verification": []any{}, "discoveries": []string{"scope must change"}, "questions": []any{}, "risks": []any{}}
	submitted := callTool(t, ctx, session, "pathframe_submit_result", map[string]any{"result": result})
	if submitted["phase"] != "blocked" || submitted["reconciliation"].(map[string]any)["requires_replan"] != true {
		t.Fatalf("needs_replan = %#v", submitted)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
}
