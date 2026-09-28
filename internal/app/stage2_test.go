package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func TestSecurityCreateRejectsManagedDirectorySymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".pathframe")); err != nil {
		t.Fatal(err)
	}
	if _, err := (Service{Dir: root}).CreateChange("escape", artifacts.Quick); err == nil {
		t.Fatal("CreateChange accepted a symlinked managed directory")
	}
	if _, err := os.Stat(filepath.Join(outside, "changes")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("outside directory changed: %v", err)
	}
}

func TestApprovalBindsIdentityAndMaterialEditInvalidates(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join("..", "..", "testdata", "artifacts", "quick")
	dir := filepath.Join(root, ".pathframe", "changes", "quick")
	copyFixture(t, fixture, dir)
	service := Service{Dir: root}
	if _, err := service.CreateChange("quick", artifacts.Quick); err != nil {
		t.Fatal(err)
	}
	approved, err := service.Approve("quick")
	if err != nil || approved.Phase != string(workflow.PhaseReady) || approved.Identity == "" {
		t.Fatalf("Approve() = %#v, %v", approved, err)
	}

	intent := filepath.Join(dir, "intent.md")
	data, _ := os.ReadFile(intent)
	if err := os.WriteFile(intent, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	status, err := service.Orient("quick")
	if err != nil || status.Phase != workflow.PhaseReady {
		t.Fatalf("format-only status = %#v, %v", status, err)
	}

	data, _ = os.ReadFile(intent)
	if err := os.WriteFile(intent, []byte(strings.Replace(string(data), "small behavior", "material behavior", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	status, err = service.Orient("quick")
	if err != nil || status.Phase != workflow.PhaseReplanning {
		t.Fatalf("material status = %#v, %v", status, err)
	}
}

func TestModeIncreaseAddsOnlyMissingArtifactsAndLoweringFails(t *testing.T) {
	root := t.TempDir()
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Quick); err != nil {
		t.Fatal(err)
	}
	intent := filepath.Join(root, ".pathframe", "changes", "demo", "intent.md")
	custom := []byte("authored content")
	if err := os.WriteFile(intent, custom, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreateChange("demo", artifacts.Standard); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(intent)
	if string(got) != string(custom) {
		t.Fatal("mode increase overwrote authored intent")
	}
	for _, name := range []string{"requirements.md", "design.md"} {
		if _, err := os.Stat(filepath.Join(filepath.Dir(intent), name)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := service.CreateChange("demo", artifacts.Quick); err == nil {
		t.Fatal("mode lowering succeeded")
	}
}

func TestCheckRejectsRoleMismatchBeforeApproval(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "quick")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "quick"), dir)
	service := Service{Dir: root}
	if _, err := service.CreateChange("quick", artifacts.Quick); err != nil {
		t.Fatal(err)
	}
	taskPath := filepath.Join(dir, "tasks", "T1.md")
	data, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, []byte(strings.Replace(string(data), "role: none", "role: laravel-implementer", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	check, err := service.Check("quick")
	if err != nil || check.Valid || len(check.Issues) == 0 || check.Issues[len(check.Issues)-1].Code != "invalid_role" {
		t.Fatalf("Check() = %#v, %v", check, err)
	}
}

func copyFixture(t *testing.T, source, target string) {
	t.Helper()
	if err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(source, path)
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

func writeBackendRole(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, ".pathframe", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	role := "schema: pathframe.role/v1\nid: backend-engineer\nmission: Implement backend work.\nreads: []\noptional_reads: []\nallowed_actions: [\"edit\",\"test\"]\nreturns: [\"summary\",\"changed_files\",\"verification\",\"discoveries\",\"questions\"]\n"
	if err := os.WriteFile(filepath.Join(dir, "backend-engineer.yaml"), []byte(role), 0o600); err != nil {
		t.Fatal(err)
	}
}
