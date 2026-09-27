package artifacts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Parse(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	return ParseBytes(filepath.ToSlash(path), data)
}

func ParseBytes(path string, data []byte) (Document, error) {
	s := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(string(data), "\r\n", "\n")))
	if !s.Scan() || s.Text() != "---" {
		return Document{}, fmt.Errorf("%s: missing opening front matter delimiter", path)
	}
	doc := Document{Path: path, Fields: map[string]string{}, Sections: map[string]string{}}
	closed := false
	for s.Scan() {
		line := s.Text()
		if line == "---" {
			closed = true
			break
		}
		key, value, ok := strings.Cut(line, ":")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" || value == "" || doc.Fields[key] != "" {
			return Document{}, fmt.Errorf("%s: invalid front matter line %q", path, line)
		}
		doc.Fields[key] = value
	}
	if !closed {
		return Document{}, fmt.Errorf("%s: unclosed front matter", path)
	}
	var heading string
	var body []string
	flush := func() {
		if heading != "" {
			doc.Sections[heading] = strings.TrimSpace(strings.Join(body, "\n"))
		}
		body = nil
	}
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "## ") {
			flush()
			heading = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			if heading == "" || doc.Sections[heading] != "" {
				return Document{}, fmt.Errorf("%s: invalid or duplicate section", path)
			}
			continue
		}
		if heading == "" {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "# ") {
				return Document{}, fmt.Errorf("%s: content outside a section", path)
			}
			continue
		}
		body = append(body, line)
	}
	flush()
	if err := s.Err(); err != nil {
		return Document{}, err
	}
	return doc, nil
}

func list(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" || value == "[]" || value == "none" {
		return nil
	}
	if strings.HasPrefix(value, "[") {
		var values []string
		if json.Unmarshal([]byte(value), &values) == nil {
			return values
		}
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func ParseTask(doc Document) (Task, error) {
	task := Task{Document: doc, ID: doc.Fields["id"], Title: doc.Fields["title"], ExecutionPolicy: doc.Fields["execution_policy"], Role: doc.Fields["role"], References: list(doc.Fields["references"]), Dependencies: list(doc.Fields["dependencies"])}
	for _, line := range strings.Split(doc.Sections["Verification"], "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-"))
		if line == "" || line == "none" {
			continue
		}
		var argv []string
		if err := json.Unmarshal([]byte(line), &argv); err != nil || len(argv) == 0 {
			return Task{}, fmt.Errorf("%s: verification must contain JSON argv arrays", doc.Path)
		}
		task.Verification = append(task.Verification, argv)
	}
	return task, nil
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
