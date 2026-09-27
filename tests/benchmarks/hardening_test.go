package benchmarks_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	pfcontext "github.com/0xkhdr/pathframe/internal/context"
	"github.com/0xkhdr/pathframe/internal/store"
)

func TestPerformanceBudgets(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, ".pathframe", "changes", "demo", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "change")
	copyTree(t, filepath.Join("..", "..", "testdata", "artifacts", "standard"), dir)
	task := artifacts.Task{Document: artifacts.Document{Sections: map[string]string{"Objective": "measure", "Acceptance": "assembled"}}, ID: "T1", Title: "Budget", ExecutionPolicy: "brain"}
	role := pfcontext.Role{Schema: pfcontext.RoleSchema, ID: "none", Mission: "measure", Reads: []string{}, OptionalReads: []string{}, AllowedActions: []string{"edit"}, Returns: []string{"summary"}}
	entries := []pfcontext.Entry{{Layer: pfcontext.Change, Path: "intent.md", Required: true, Content: "intent", Bytes: 6}}

	budget(t, "Discovery", 10*time.Millisecond, func() { _, _ = store.Discover(nested) })
	budget(t, "Validation", 25*time.Millisecond, func() { artifacts.Validate(dir) })
	budget(t, "Packet", 2*time.Millisecond, func() {
		_, _ = pfcontext.Assemble("demo", task, role, entries, nil, "phase=ready", pfcontext.DefaultBudget)
	})
}

func budget(t *testing.T, name string, limit time.Duration, operation func()) {
	t.Helper()
	const runs = 100
	start := time.Now()
	for range runs {
		operation()
	}
	if average := time.Since(start) / runs; average > limit {
		t.Errorf("%s average %s exceeds %s budget", name, average, limit)
	}
}

func BenchmarkDiscovery(b *testing.B) {
	root := b.TempDir()
	nested := filepath.Join(root, ".pathframe", "changes", "demo", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for range b.N {
		if _, err := store.Discover(nested); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidation(b *testing.B) {
	dir := filepath.Join(b.TempDir(), "change")
	copyTree(b, filepath.Join("..", "..", "testdata", "artifacts", "standard"), dir)
	b.ResetTimer()
	for range b.N {
		if _, issues := artifacts.Validate(dir); len(issues) != 0 {
			b.Fatal(issues)
		}
	}
}

func BenchmarkPacketAssembly(b *testing.B) {
	task := artifacts.Task{Document: artifacts.Document{Sections: map[string]string{"Objective": "measure packet assembly", "Acceptance": "packet assembled"}}, ID: "T1", Title: "Benchmark", ExecutionPolicy: "brain"}
	role := pfcontext.Role{Schema: pfcontext.RoleSchema, ID: "none", Mission: "benchmark", Reads: []string{}, OptionalReads: []string{}, AllowedActions: []string{"edit"}, Returns: []string{"summary"}}
	entries := []pfcontext.Entry{{Layer: pfcontext.Change, Path: "intent.md", Required: true, Content: "intent", Bytes: 6}}
	b.ResetTimer()
	for range b.N {
		if _, err := pfcontext.Assemble("demo", task, role, entries, nil, "phase=ready", pfcontext.DefaultBudget); err != nil {
			b.Fatal(err)
		}
	}
}

func copyTree(tb testing.TB, source, target string) {
	tb.Helper()
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
		tb.Fatal(err)
	}
}
