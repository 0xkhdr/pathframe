package artifacts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStageTwoSchemasAreValidJSON(t *testing.T) {
	for _, name := range []string{"artifact-instructions-v1.schema.json", "change-v1.schema.json", "plan-check-v1.schema.json", "plan-approval-v1.schema.json", "intent-v1.schema.json", "requirements-v1.schema.json", "design-v1.schema.json", "task-v1.schema.json", "role-v1.schema.json", "task-packet-v1.schema.json", "risks-v1.schema.json", "rollout-v1.schema.json", "recovery-v1.schema.json"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "schemas", name))
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil || schema["$id"] == nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
}

func TestTemplatesAndStrictParser(t *testing.T) {
	for _, mode := range []Mode{Quick, Standard, HighRisk} {
		for _, kind := range append(RequiredArtifacts(mode), "task") {
			body, err := Template(mode, kind)
			if err != nil {
				t.Fatalf("Template(%s, %s): %v", mode, kind, err)
			}
			doc, err := ParseBytes(kind+".md", []byte(body))
			if err != nil {
				t.Fatalf("ParseBytes(%s): %v", kind, err)
			}
			if doc.Fields["profile"] != Profile {
				t.Fatalf("profile = %q", doc.Fields["profile"])
			}
			for _, section := range artifactSections[kind] {
				if _, ok := doc.Sections[section]; !ok {
					t.Fatalf("%s lacks %s", kind, section)
				}
			}
		}
	}
}

func TestValidateFixtureModesAndMaterialIdentity(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "artifacts")
	for _, mode := range []Mode{Quick, Standard, HighRisk} {
		t.Run(string(mode), func(t *testing.T) {
			plan, issues := Validate(filepath.Join(root, string(mode)))
			if len(issues) != 0 {
				t.Fatalf("issues = %#v", issues)
			}
			if plan.Identity == "" || plan.Change.Mode != mode {
				t.Fatalf("plan = %#v", plan)
			}
		})
	}
}

func TestValidationFailures(t *testing.T) {
	source := filepath.Join("..", "..", "testdata", "artifacts", "standard")
	dir := t.TempDir()
	copyTree(t, source, dir)
	task := filepath.Join(dir, "tasks", "T1.md")
	data, _ := os.ReadFile(task)
	data = []byte(strings.Replace(string(data), "references: [\"REQ-1\"]", "references: [\"REQ-missing\"]", 1))
	if err := os.WriteFile(task, data, 0o600); err != nil {
		t.Fatal(err)
	}
	_, issues := Validate(dir)
	if !hasIssue(issues, "broken_reference") {
		t.Fatalf("issues = %#v", issues)
	}
	data = []byte(strings.Replace(string(data), "- none\n\n## Assumptions", "- decide API\n\n## Assumptions", 1))
	if err := os.WriteFile(task, data, 0o600); err != nil {
		t.Fatal(err)
	}
	_, issues = Validate(dir)
	if !hasIssue(issues, "unresolved_question") {
		t.Fatalf("issues = %#v", issues)
	}
}

func hasIssue(issues []Issue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func copyTree(t *testing.T, source, target string) {
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
