package verification

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestRunBoundsAndContainsVerification(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Stage 7 supports Linux amd64 first")
	}
	root := t.TempDir()
	result := Run(context.Background(), root, Command{Argv: []string{"printf", "abcdef"}, Workdir: ".", Timeout: time.Second, MaxBytes: 3})
	if result.ExitCode != 0 || result.Stdout != "abc" || !result.StdoutTruncated {
		t.Fatalf("bounded result = %#v", result)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	result = Run(context.Background(), root, Command{Argv: []string{"true"}, Workdir: "escape", Timeout: time.Second, MaxBytes: 10})
	if result.ExitCode != -1 || result.Stderr == "" {
		t.Fatalf("symlink escape = %#v", result)
	}
}

func TestRunTimeoutAndInterruption(t *testing.T) {
	root := t.TempDir()
	timed := Run(context.Background(), root, Command{Argv: []string{"sleep", "1"}, Workdir: ".", Timeout: 10 * time.Millisecond, MaxBytes: 10})
	if !timed.TimedOut || timed.ExitCode == 0 {
		t.Fatalf("timeout = %#v", timed)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	interrupted := Run(ctx, root, Command{Argv: []string{"sleep", "1"}, Workdir: ".", Timeout: time.Second, MaxBytes: 10})
	if !interrupted.Interrupted {
		t.Fatalf("interruption = %#v", interrupted)
	}
}

func TestIdentityAndScope(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "a.txt")
	if err := os.WriteFile(path, []byte("one"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := ContentIdentity(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("two"), 0o600); err != nil {
		t.Fatal(err)
	}
	after, err := ContentIdentity(root)
	if err != nil || before == after {
		t.Fatalf("identity before=%s after=%s err=%v", before, after, err)
	}
	violations := ScopeViolations([]string{"internal/a.go", "README.md"}, []string{"internal/**"})
	if len(violations) != 1 || violations[0] != "README.md" {
		t.Fatalf("violations = %#v", violations)
	}
}

func TestSecurityRunDoesNotInterpretShellSyntax(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "injected")
	argument := "$(touch " + marker + ")"
	result := Run(context.Background(), root, Command{Argv: []string{"printf", "%s", argument}, Workdir: ".", Timeout: time.Second, MaxBytes: 1024})
	if result.ExitCode != 0 || result.Stdout != argument {
		t.Fatalf("structured argv result = %#v", result)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("shell syntax executed: %v", err)
	}
}
