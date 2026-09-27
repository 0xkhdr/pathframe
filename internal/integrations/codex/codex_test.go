package codex

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallDoctorAndNonDestructiveUpdate(t *testing.T) {
	root := t.TempDir()
	result, err := Install(root, "/usr/local/bin/pathframe", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changed) != 4 || !Doctor(root).Healthy {
		t.Fatalf("install = %#v, doctor = %#v", result, Doctor(root))
	}
	skillPath := filepath.Join(root, ".agents", "skills", "pathframe", "SKILL.md")
	skill, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, rule := range []string{"must_use", "offer", "must_not_use", "Never construct Pathframe CLI commands", "explicitly"} {
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

func TestManifestSchemaAndOptionalHook(t *testing.T) {
	schemaData, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "integration-v1.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(schemaData, &schema); err != nil || schema["$id"] != IntegrationSchema {
		t.Fatalf("integration schema: %v, id=%v", err, schema["$id"])
	}
	root := t.TempDir()
	if _, err := Install(root, "/bin/pathframe", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "hooks.json")); !os.IsNotExist(err) {
		t.Fatalf("optional hook exists: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".pathframe", "integrations", "codex", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Schema != IntegrationSchema || manifest.Version != IntegrationVersion || manifest.Transport != "stdio" {
		t.Fatalf("manifest = %#v", manifest)
	}
}

func TestActivationEvaluationFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testdata", "evals")
	for _, name := range []string{"must-use.yaml", "offer.yaml", "must-not-use.yaml"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "expected: "+strings.TrimSuffix(strings.ReplaceAll(name, "-", "_"), ".yaml")) {
			t.Fatalf("%s has no matching expected class", name)
		}
	}
}
