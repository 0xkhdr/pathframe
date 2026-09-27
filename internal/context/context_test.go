package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBudgetNeverTruncatesRequiredAndReportsOptional(t *testing.T) {
	entries := []Entry{{Layer: Change, Path: "required", Required: true, Content: "1234", Bytes: 4}, {Layer: Foundation, Path: "optional", Content: "5678", Bytes: 4}}
	if _, _, _, err := ApplyBudget(entries, 3); err == nil || !strings.Contains(err.Error(), "split the task") {
		t.Fatalf("required over budget error = %v", err)
	}
	included, omitted, budget, err := ApplyBudget(entries, 4)
	if err != nil || len(included) != 1 || len(omitted) != 1 || budget.OmittedBytes != 4 || included[0].Content != "1234" {
		t.Fatalf("ApplyBudget() = %#v, %#v, %#v, %v", included, omitted, budget, err)
	}
}

func TestBudgetReservesSpaceForLaterRequiredContext(t *testing.T) {
	entries := []Entry{{Layer: Foundation, Path: "optional", Content: "1234", Bytes: 4}, {Layer: Runtime, Path: "required", Required: true, Content: "5678", Bytes: 4}}
	included, omitted, budget, err := ApplyBudget(entries, 4)
	if err != nil || len(included) != 1 || !included[0].Required || len(omitted) != 1 || budget.IncludedBytes != 4 {
		t.Fatalf("ApplyBudget() = %#v, %#v, %#v, %v", included, omitted, budget, err)
	}
}

func TestResolveRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	change := filepath.Join(root, ".pathframe", "changes", "demo")
	if err := os.MkdirAll(change, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(outside, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(change, "escape")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Resolve(root, change, []Reference{{Layer: Change, Path: "escape", Required: true}}); err == nil {
		t.Fatal("Resolve() accepted symlink escape")
	}
}

func TestResolveReportsMissingOptional(t *testing.T) {
	root := t.TempDir()
	change := filepath.Join(root, ".pathframe", "changes", "demo")
	if err := os.MkdirAll(change, 0o755); err != nil {
		t.Fatal(err)
	}
	entries, omissions, err := Resolve(root, change, []Reference{{Layer: Foundation, Path: "missing.md"}})
	if err != nil || len(entries) != 0 || len(omissions) != 1 || omissions[0].Reason == "" {
		t.Fatalf("Resolve() = %#v, %#v, %v", entries, omissions, err)
	}
}
