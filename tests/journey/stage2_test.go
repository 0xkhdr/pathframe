package journey_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQuickStandardHighRiskApprovalJourney(t *testing.T) {
	binary := buildPathframe(t)
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	for _, mode := range []string{"quick", "standard", "high-risk"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			if got := run(t, binary, root, "new", "--change", "demo", "--mode", mode); !strings.Contains(got, "Valid: false") {
				t.Fatalf("new = %s", got)
			}
			source := filepath.Join(repository, "testdata", "artifacts", mode)
			target := filepath.Join(root, ".pathframe", "changes", "demo")
			copyAuthored(t, source, target)
			if mode != "quick" {
				roleDir := filepath.Join(root, ".pathframe", "roles")
				if err := os.MkdirAll(roleDir, 0o755); err != nil {
					t.Fatal(err)
				}
				role := "schema: pathframe.role/v1\nid: backend-engineer\nmission: Implement bounded work.\nreads: []\noptional_reads: []\nallowed_actions: [\"edit\",\"test\"]\nreturns: [\"summary\",\"changed_files\",\"verification\",\"discoveries\",\"questions\"]\n"
				if err := os.WriteFile(filepath.Join(roleDir, "backend-engineer.yaml"), []byte(role), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if got := run(t, binary, root, "check", "--change", "demo", "--json"); !strings.Contains(got, `"valid": true`) {
				t.Fatalf("check = %s", got)
			}
			if got := run(t, binary, root, "approve", "--change", "demo", "--json"); !strings.Contains(got, `"phase": "ready"`) {
				t.Fatalf("approve = %s", got)
			}
			if got := run(t, binary, root, "status", "--change", "demo", "--json"); !strings.Contains(got, `"phase": "ready"`) {
				t.Fatalf("status = %s", got)
			}
		})
	}
}

func TestApprovalInvalidatesAfterPolicyEditJourney(t *testing.T) {
	binary := buildPathframe(t)
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	root := t.TempDir()
	run(t, binary, root, "new", "--change", "demo", "--mode", "quick")
	target := filepath.Join(root, ".pathframe", "changes", "demo")
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "quick"), target)
	run(t, binary, root, "approve", "--change", "demo")
	task := filepath.Join(target, "tasks", "T1.md")
	data, _ := os.ReadFile(task)
	if err := os.WriteFile(task, []byte(strings.Replace(string(data), "execution_policy: brain", "execution_policy: delegated", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := run(t, binary, root, "status", "--change", "demo", "--json"); !strings.Contains(got, `"phase": "replanning"`) {
		t.Fatalf("status = %s", got)
	}
}

func copyAuthored(t *testing.T, source, target string) {
	t.Helper()
	if err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(source, path)
		if rel == "change.yaml" {
			return nil
		}
		out := filepath.Join(target, rel)
		if entry.IsDir() {
			return os.MkdirAll(out, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(out, data, 0o600)
	}); err != nil {
		t.Fatal(err)
	}
}
