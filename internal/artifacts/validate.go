package artifacts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var requirementID = regexp.MustCompile(`(?m)^-\s*(REQ-[A-Za-z0-9._-]+):\s*\S`)

func LoadChange(dir string) (Change, error) {
	data, err := os.ReadFile(filepath.Join(dir, "change.yaml"))
	if err != nil {
		return Change{}, err
	}
	doc, err := ParseBytes("change.yaml", append([]byte("---\n"), append(data, []byte("---\n")...)...))
	if err != nil {
		return Change{}, err
	}
	c := Change{Schema: doc.Fields["schema"], Profile: doc.Fields["profile"], ID: doc.Fields["id"], Mode: Mode(doc.Fields["mode"])}
	if c.Schema != ChangeSchema || c.Profile != Profile || c.ID == "" || !c.Mode.Valid() {
		return Change{}, fmt.Errorf("invalid change.yaml")
	}
	return c, nil
}

func Validate(dir string) (Plan, []Issue) {
	change, err := LoadChange(dir)
	if err != nil {
		return Plan{}, []Issue{{Path: "change.yaml", Code: "invalid_change", Message: err.Error()}}
	}
	plan := Plan{Change: change}
	var issues []Issue
	for _, kind := range RequiredArtifacts(change.Mode) {
		path := filepath.Join(dir, kind+".md")
		doc, err := Parse(path)
		if err != nil {
			issues = append(issues, Issue{Path: kind + ".md", Code: "invalid_artifact", Message: err.Error()})
			continue
		}
		issues = append(issues, validateDocument(doc, kind)...)
		plan.Documents = append(plan.Documents, doc)
	}
	taskPaths, _ := filepath.Glob(filepath.Join(dir, "tasks", "*.md"))
	sort.Strings(taskPaths)
	if len(taskPaths) == 0 {
		issues = append(issues, Issue{Path: "tasks", Code: "missing_task", Message: "at least one task file is required"})
	}
	seenTasks := map[string]bool{}
	for _, path := range taskPaths {
		doc, err := Parse(path)
		if err != nil {
			issues = append(issues, Issue{Path: filepath.ToSlash(path), Code: "invalid_task", Message: err.Error()})
			continue
		}
		issues = append(issues, validateDocument(doc, "task")...)
		task, err := ParseTask(doc)
		if err != nil {
			issues = append(issues, Issue{Path: doc.Path, Code: "invalid_verification", Message: err.Error()})
			continue
		}
		if task.ID == "" || task.Title == "" || (task.ExecutionPolicy != "brain" && task.ExecutionPolicy != "delegated") {
			issues = append(issues, Issue{Path: doc.Path, Code: "invalid_task_fields", Message: "task requires id, title, and brain or delegated execution_policy"})
		}
		if seenTasks[task.ID] {
			issues = append(issues, Issue{Path: doc.Path, Code: "duplicate_task", Message: "task ID is not unique"})
		}
		seenTasks[task.ID] = true
		if len(task.Verification) == 0 {
			issues = append(issues, Issue{Path: doc.Path, Code: "missing_verification", Message: "task requires at least one verification argv"})
		}
		plan.Tasks = append(plan.Tasks, task)
		plan.Documents = append(plan.Documents, doc)
	}
	validReqs := map[string]bool{}
	for _, doc := range plan.Documents {
		if filepath.Base(doc.Path) == "requirements.md" {
			for _, match := range requirementID.FindAllStringSubmatch(doc.Sections["Requirements"], -1) {
				validReqs[match[1]] = true
			}
		}
	}
	for _, task := range plan.Tasks {
		for _, ref := range task.References {
			if !validReqs[ref] {
				issues = append(issues, Issue{Path: task.Path, Code: "broken_reference", Message: "unknown reference " + ref})
			}
		}
		for _, dependency := range task.Dependencies {
			if !seenTasks[dependency] || dependency == task.ID {
				issues = append(issues, Issue{Path: task.Path, Code: "broken_reference", Message: "unknown or self dependency " + dependency})
			}
		}
		for _, read := range bulletList(task.Sections["Required Reads"]) {
			if read == "none" {
				continue
			}
			clean := filepath.Clean(read)
			if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
				issues = append(issues, Issue{Path: task.Path, Code: "broken_reference", Message: "required read escapes change: " + read})
				continue
			}
			if _, err := os.Stat(filepath.Join(dir, clean)); err != nil {
				issues = append(issues, Issue{Path: task.Path, Code: "broken_reference", Message: "required read not found: " + read})
			}
		}
	}
	if len(issues) == 0 {
		plan.Identity = identity(plan)
	}
	return plan, issues
}

func validateDocument(doc Document, kind string) []Issue {
	var issues []Issue
	allowedFields := map[string]bool{"schema": true, "profile": true}
	if kind == "task" {
		for _, key := range []string{"id", "title", "execution_policy", "role", "references", "dependencies"} {
			allowedFields[key] = true
		}
	}
	for key := range doc.Fields {
		if !allowedFields[key] {
			issues = append(issues, Issue{Path: doc.Path, Code: "unexpected_field", Message: "unexpected front matter field " + key})
		}
	}
	allowedSections := map[string]bool{}
	for _, section := range artifactSections[kind] {
		allowedSections[section] = true
	}
	for section := range doc.Sections {
		if !allowedSections[section] {
			issues = append(issues, Issue{Path: doc.Path, Code: "unexpected_section", Message: "unexpected section " + section})
		}
	}
	if doc.Fields["schema"] != artifactSchemas[kind] || doc.Fields["profile"] != Profile {
		issues = append(issues, Issue{Path: doc.Path, Code: "wrong_schema", Message: "schema/profile does not match artifact"})
	}
	for _, section := range artifactSections[kind] {
		body, ok := doc.Sections[section]
		if !ok || strings.TrimSpace(body) == "" {
			issues = append(issues, Issue{Path: doc.Path, Code: "missing_section", Message: "missing content for " + section})
			continue
		}
		lower := strings.ToLower(body)
		if strings.Contains(body, "<") || strings.Contains(body, ">") || strings.Contains(lower, "{{") || strings.Contains(lower, "todo") {
			issues = append(issues, Issue{Path: doc.Path, Code: "placeholder", Message: "unresolved placeholder in " + section})
		}
		if section == "Questions" && strings.TrimSpace(strings.TrimPrefix(body, "-")) != "none" {
			issues = append(issues, Issue{Path: doc.Path, Code: "unresolved_question", Message: "questions must be resolved or 'none'"})
		}
		if (section == "Acceptance" || section == "Objective" || section == "Summary" || section == "Outcomes") && strings.TrimSpace(strings.TrimPrefix(body, "-")) == "none" {
			issues = append(issues, Issue{Path: doc.Path, Code: "missing_" + strings.ToLower(section), Message: section + " cannot be 'none'"})
		}
	}
	return issues
}

func bulletList(body string) []string {
	var values []string
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") {
			values = append(values, strings.TrimSpace(strings.TrimPrefix(line, "-")))
		}
	}
	return values
}

func identity(plan Plan) string {
	type canonicalDoc struct {
		Path     string
		Fields   [][2]string
		Sections [][2]string
	}
	canonical := struct {
		Change Change
		Docs   []canonicalDoc
	}{Change: plan.Change}
	for _, doc := range plan.Documents {
		item := canonicalDoc{Path: filepath.Base(doc.Path)}
		for _, key := range sortedKeys(doc.Fields) {
			item.Fields = append(item.Fields, [2]string{key, doc.Fields[key]})
		}
		for _, key := range sortedKeys(doc.Sections) {
			item.Sections = append(item.Sections, [2]string{key, doc.Sections[key]})
		}
		canonical.Docs = append(canonical.Docs, item)
	}
	data, _ := json.Marshal(canonical)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
