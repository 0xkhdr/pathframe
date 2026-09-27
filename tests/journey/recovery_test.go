package journey_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorRecoveryCorruptProjectionJourney(t *testing.T) {
	binary := buildPathframe(t)
	root := t.TempDir()
	run(t, binary, root, "new", "--change", "demo", "--mode", "quick")
	dir := filepath.Join(root, ".pathframe", "changes", "demo")
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	copyAuthored(t, filepath.Join(repository, "testdata", "artifacts", "quick"), dir)
	intent, err := os.ReadFile(filepath.Join(dir, "intent.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "state.json"), []byte("corrupt"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := runFailure(t, binary, root, "doctor", "--change", "demo", "--json"); !strings.Contains(got, `"code": "corrupt_projection"`) {
		t.Fatalf("doctor = %s", got)
	}
	if got := run(t, binary, root, "doctor", "--change", "demo", "--repair", "--json"); !strings.Contains(got, `"repaired": true`) {
		t.Fatalf("repair = %s", got)
	}
	after, _ := os.ReadFile(filepath.Join(dir, "intent.md"))
	if string(after) != string(intent) {
		t.Fatal("Doctor changed authored intent")
	}
}

func TestRecoveryCancellationRetainsHistoryJourney(t *testing.T) {
	binary := buildPathframe(t)
	root := t.TempDir()
	run(t, binary, root, "new", "--change", "demo", "--mode", "quick")
	run(t, binary, root, "cancel", "--change", "demo")
	history, err := os.ReadFile(filepath.Join(root, ".pathframe", "changes", "demo", "history.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(history), `"action":"cancel"`) {
		t.Fatal("cancellation history was not retained")
	}
}
