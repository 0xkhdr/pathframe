package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/0xkhdr/pathframe/internal/artifacts"
)

func TestPrepareDelegationPacketAndVisibleOmission(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "standard"), dir)
	taskData, err := os.ReadFile(filepath.Join(dir, "tasks", "T1.md"))
	if err != nil {
		t.Fatal(err)
	}
	unrelated := strings.Replace(string(taskData), "id: T1", "id: T2", 1)
	unrelated = strings.Replace(unrelated, "title: Implement the standard feature", "title: Unrelated task", 1)
	unrelated = strings.Replace(unrelated, "Implement the deterministic result.", "UNRELATED TASK CONTENT", 1)
	if err := os.WriteFile(filepath.Join(dir, "tasks", "T2.md"), []byte(unrelated), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".pathframe", "roles"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "foundation.md"), []byte("required foundation"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "optional.md"), []byte(strings.Repeat("optional", 100)), 0o600); err != nil {
		t.Fatal(err)
	}
	role := "schema: pathframe.role/v1\nid: backend-engineer\nmission: Implement backend work.\nreads: [\"foundation.md\"]\noptional_reads: [\"optional.md\"]\nallowed_actions: [\"edit\",\"test\"]\nreturns: [\"summary\",\"changed_files\",\"verification\",\"discoveries\",\"questions\"]\n"
	if err := os.WriteFile(filepath.Join(root, ".pathframe", "roles", "backend-engineer.yaml"), []byte(role), 0o600); err != nil {
		t.Fatal(err)
	}
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Standard); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Approve("demo"); err != nil {
		t.Fatal(err)
	}
	large, err := service.PrepareDelegation(PrepareInput{Change: "demo", Task: "T1"})
	if err != nil || large.Packet == nil {
		t.Fatalf("large preview = %#v, %v", large, err)
	}
	got, err := service.PrepareDelegation(PrepareInput{Change: "demo", Task: "T1", BudgetBytes: large.Packet.Budget.RequiredBytes})
	if err != nil || got.Packet == nil || len(got.Packet.Omissions) != 1 {
		t.Fatalf("preview = %#v, %v", got, err)
	}
	for _, entry := range got.Packet.Context {
		if strings.Contains(entry.Content, "optionaloptional") {
			t.Fatal("omitted optional context was included")
		}
		if strings.Contains(entry.Path, "T2") || strings.Contains(entry.Content, "UNRELATED") {
			t.Fatal("unrelated task context was included")
		}
	}
	if got.Packet.WriteScopeAssurance != "advisory" || got.Packet.HostAssurance != "not_evaluated" {
		t.Fatalf("assurance = %#v", got.Packet)
	}
}

func TestPrepareDelegationRequiredBudgetBlocks(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "quick"), dir)
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Quick); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Approve("demo"); err != nil {
		t.Fatal(err)
	}
	got, err := service.PrepareDelegation(PrepareInput{Change: "demo", BudgetBytes: 1})
	if err != nil || got.Packet != nil || len(got.Issues) != 1 || !strings.Contains(got.Issues[0], "split the task") {
		t.Fatalf("preview = %#v, %v", got, err)
	}
}

func TestCheckMissingRoleBlocksBeforeApproval(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	copyFixture(t, filepath.Join("..", "..", "testdata", "artifacts", "standard"), dir)
	service := Service{Dir: root}
	if _, err := service.CreateChange("demo", artifacts.Standard); err != nil {
		t.Fatal(err)
	}
	got, err := service.Check("demo")
	if err != nil || got.Valid || len(got.Issues) == 0 || got.Issues[len(got.Issues)-1].Code != "unknown_role" {
		t.Fatalf("check = %#v, %v", got, err)
	}
	view, err := service.Orient("demo")
	if err != nil || view.Phase != "planning" {
		t.Fatalf("check mutated state: %#v, %v", view, err)
	}
}
