package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pfcontext "github.com/0xkhdr/pathframe/internal/context"
)

func TestInstallDoctorAndSafeUpdate(t *testing.T) {
	root := t.TempDir()
	result, err := Install(root, "/usr/local/bin/pathframe", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changed) != 5 || !Doctor(root).Healthy {
		t.Fatalf("install = %#v, doctor = %#v", result, Doctor(root))
	}
	for _, name := range []string{".claude/skills/pathframe/SKILL.md", ".claude/commands/pathframe.md", ".mcp.json", ".claude/settings.json"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(name))); err != nil {
			t.Fatal(err)
		}
	}
	skillPath := filepath.Join(root, ".claude", "skills", "pathframe", "SKILL.md")
	skill, _ := os.ReadFile(skillPath)
	for _, rule := range []string{"must_use", "offer", "must_not_use", "Never construct Pathframe CLI commands", "pathframe_prepare_delegation", "never launches", "explicit"} {
		if !strings.Contains(string(skill), rule) {
			t.Fatalf("skill missing %q", rule)
		}
	}
	if err := os.WriteFile(skillPath, append(skill, []byte("custom\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(root, "/usr/local/bin/pathframe", true); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("modified install error = %v", err)
	}
	if Doctor(root).Healthy {
		t.Fatal("doctor accepted modified generated file")
	}
}

func TestWorkerInstructionsKeepPinkyBounded(t *testing.T) {
	got := WorkerInstructions(pfcontext.Packet{Schema: pfcontext.PacketSchema, Change: "demo", Task: "T1", ExecutionPolicy: "delegated"})
	for _, want := range []string{"You are Pinky", "pathframe.task-result/v1", "Do not change Pathframe plan state", `"task": "T1"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("instructions missing %q", want)
		}
	}
}

func TestManifestConfigurationAndOptionalHook(t *testing.T) {
	root := t.TempDir()
	if _, err := Install(root, "/bin/pathframe", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".claude", "settings.json")); !os.IsNotExist(err) {
		t.Fatalf("optional hook exists: %v", err)
	}
	for _, name := range []string{".mcp.json", ".pathframe/integrations/claude/manifest.json"} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Fatal(err)
		}
		var value map[string]any
		if err := json.Unmarshal(data, &value); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	data, _ := os.ReadFile(filepath.Join(root, ".pathframe", "integrations", "claude", "manifest.json"))
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil || manifest.Host != "claude-code" || manifest.Transport != "stdio" {
		t.Fatalf("manifest = %#v, err = %v", manifest, err)
	}
}

func TestInstallRefusesUserConfigAndRemovesOwnedOptionalHook(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".mcp.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(root, "/bin/pathframe", false); err == nil || !strings.Contains(err.Error(), "user-owned") {
		t.Fatalf("user config error = %v", err)
	}
	root = t.TempDir()
	if _, err := Install(root, "/bin/pathframe", true); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(root, "/bin/pathframe", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".claude", "settings.json")); !os.IsNotExist(err) {
		t.Fatalf("owned optional hook was not removed: %v", err)
	}
	if !Doctor(root).Healthy {
		t.Fatalf("doctor = %#v", Doctor(root))
	}
}

func TestSharedActivationFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testdata", "evals")
	for _, name := range []string{"must-use.yaml", "offer.yaml", "must-not-use.yaml"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		want := "expected: " + strings.TrimSuffix(strings.ReplaceAll(name, "-", "_"), ".yaml")
		if !strings.Contains(string(data), want) {
			t.Fatalf("%s has no matching expected class", name)
		}
	}
}
