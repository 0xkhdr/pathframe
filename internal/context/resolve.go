package context

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadRole(path string) (Role, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Role{}, err
	}
	role := Role{}
	fields := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(string(data), "\r\n", "\n")))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" || value == "" || fields[key] != "" {
			return Role{}, fmt.Errorf("%s: invalid role line %q", path, line)
		}
		fields[key] = value
	}
	if err := scanner.Err(); err != nil {
		return Role{}, err
	}
	allowed := map[string]bool{"schema": true, "id": true, "mission": true, "reads": true, "optional_reads": true, "allowed_actions": true, "returns": true}
	for key := range fields {
		if !allowed[key] {
			return Role{}, fmt.Errorf("%s: unexpected role field %s", path, key)
		}
	}
	role.Schema, role.ID, role.Mission = fields["schema"], fields["id"], fields["mission"]
	for key, target := range map[string]*[]string{"reads": &role.Reads, "optional_reads": &role.OptionalReads, "allowed_actions": &role.AllowedActions, "returns": &role.Returns} {
		if fields[key] == "" {
			continue
		}
		if err := json.Unmarshal([]byte(fields[key]), target); err != nil {
			return Role{}, fmt.Errorf("%s: %s must be a JSON string array", path, key)
		}
	}
	if role.Schema != RoleSchema || role.ID == "" || role.Mission == "" || fields["reads"] == "" || fields["optional_reads"] == "" || len(role.AllowedActions) == 0 || len(role.Returns) == 0 {
		return Role{}, fmt.Errorf("%s: role requires schema, id, mission, reads, optional_reads, allowed_actions, and returns", path)
	}
	return role, nil
}

func Resolve(root, changeDir string, references []Reference) ([]Entry, []Omission, error) {
	entries := make([]Entry, 0, len(references))
	omissions := []Omission{}
	for _, ref := range references {
		base := root
		if ref.Layer == Change {
			base = changeDir
		}
		if ref.Layer != Foundation && ref.Layer != Change {
			return nil, nil, fmt.Errorf("file reference has unsupported layer %q", ref.Layer)
		}
		path, err := containedFile(base, ref.Path)
		if err != nil {
			if !ref.Required && errors.Is(err, os.ErrNotExist) {
				omissions = append(omissions, Omission{Layer: ref.Layer, Path: filepath.ToSlash(ref.Path), Reason: "optional context is unavailable"})
				continue
			}
			return nil, nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			if !ref.Required && errors.Is(err, os.ErrNotExist) {
				omissions = append(omissions, Omission{Layer: ref.Layer, Path: filepath.ToSlash(ref.Path), Reason: "optional context is unavailable"})
				continue
			}
			return nil, nil, fmt.Errorf("resolve %s context %q: %w", ref.Layer, ref.Path, err)
		}
		entries = append(entries, Entry{Layer: ref.Layer, Path: filepath.ToSlash(ref.Path), Required: ref.Required, Content: string(data), Bytes: len(data)})
	}
	return entries, omissions, nil
}

func containedFile(base, reference string) (string, error) {
	clean := filepath.Clean(reference)
	if reference == "" || filepath.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("context reference escapes its %s root: %q", filepath.Base(base), reference)
	}
	resolvedBase, err := filepath.EvalSymlinks(base)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(filepath.Join(resolvedBase, clean))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(resolvedBase, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("context reference escapes its %s root: %q", filepath.Base(base), reference)
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("context reference is not a regular file: %q", reference)
	}
	return resolved, nil
}
