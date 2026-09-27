package codex

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	IntegrationSchema  = "pathframe.integration/v1"
	IntegrationVersion = "1.0.0"
	generatedMarker    = "pathframe-generated:pathframe.integration/v1"
)

type Manifest struct {
	Schema       string            `json:"schema"`
	Version      string            `json:"version"`
	Host         string            `json:"host"`
	Transport    string            `json:"transport"`
	Capabilities []string          `json:"capabilities"`
	Files        map[string]string `json:"files"`
}

type InstallResult struct {
	Schema  string   `json:"schema"`
	Version string   `json:"version"`
	Changed []string `json:"changed"`
}

func Install(root, executable string, sessionOrientation bool) (InstallResult, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return InstallResult{}, err
	}
	assets, manifest, err := generated(root, executable, sessionOrientation)
	if err != nil {
		return InstallResult{}, err
	}
	old, _ := readManifest(root)
	for name, body := range assets {
		path := filepath.Join(root, filepath.FromSlash(name))
		current, readErr := os.ReadFile(path)
		if os.IsNotExist(readErr) || string(current) == body {
			continue
		}
		if readErr != nil {
			return InstallResult{}, readErr
		}
		if old.Files[name] == "" || old.Files[name] != digest(current) {
			return InstallResult{}, fmt.Errorf("refusing to overwrite user-owned or modified %s; move it aside or restore the generated version", name)
		}
	}
	for name, want := range old.Files {
		if _, stillGenerated := assets[name]; stillGenerated {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(name))
		current, readErr := os.ReadFile(path)
		if readErr == nil && digest(current) != want {
			return InstallResult{}, fmt.Errorf("refusing to remove modified generated file %s; move it aside first", name)
		}
		if readErr != nil && !os.IsNotExist(readErr) {
			return InstallResult{}, readErr
		}
	}
	changed := []string{}
	for name := range old.Files {
		if _, stillGenerated := assets[name]; !stillGenerated {
			path := filepath.Join(root, filepath.FromSlash(name))
			if err := os.Remove(path); err == nil {
				changed = append(changed, name)
			} else if !os.IsNotExist(err) {
				return InstallResult{}, err
			}
		}
	}
	for name, body := range assets {
		path := filepath.Join(root, filepath.FromSlash(name))
		current, _ := os.ReadFile(path)
		if string(current) == body {
			continue
		}
		if err := write(path, []byte(body)); err != nil {
			return InstallResult{}, err
		}
		changed = append(changed, name)
	}
	manifest.Files = map[string]string{}
	for name, body := range assets {
		manifest.Files[name] = digest([]byte(body))
	}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	data = append(data, '\n')
	manifestName := ".pathframe/integrations/codex/manifest.json"
	manifestPath := filepath.Join(root, filepath.FromSlash(manifestName))
	currentManifest, _ := os.ReadFile(manifestPath)
	if string(currentManifest) != string(data) {
		if err := write(manifestPath, data); err != nil {
			return InstallResult{}, err
		}
		changed = append(changed, manifestName)
	}
	return InstallResult{Schema: IntegrationSchema, Version: IntegrationVersion, Changed: changed}, nil
}

func generated(root, executable string, sessionOrientation bool) (map[string]string, Manifest, error) {
	executable, err := filepath.Abs(executable)
	if err != nil {
		return nil, Manifest{}, err
	}
	q := func(value string) string { return fmt.Sprintf("%q", value) }
	tools := `[
  "pathframe_orient",
  "pathframe_assess_request",
  "pathframe_create_change",
  "pathframe_get_template",
  "pathframe_validate_plan",
  "pathframe_get_next",
  "pathframe_recover",
]`
	config := "# " + generatedMarker + "\n[mcp_servers.pathframe]\ncommand = " + q(executable) + "\nargs = [\"mcp\"]\ncwd = " + q(root) + "\nenabled_tools = " + tools + "\n"
	assets := map[string]string{
		".agents/skills/pathframe/SKILL.md": skill,
		".codex/config.toml":                config,
	}
	capabilities := []string{"repository_skill", "stdio_mcp", "typed_planning_tools"}
	if sessionOrientation {
		assets[".codex/hooks.json"] = hooks
		capabilities = append(capabilities, "session_start_orientation")
	}
	return assets, Manifest{Schema: IntegrationSchema, Version: IntegrationVersion, Host: "codex", Transport: "stdio", Capabilities: capabilities}, nil
}

func readManifest(root string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(root, ".pathframe", "integrations", "codex", "manifest.json"))
	if err != nil {
		return Manifest{}, err
	}
	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	return manifest, err
}

func write(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".pathframe-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err = tmp.Write(data); err == nil {
		err = tmp.Close()
	} else {
		_ = tmp.Close()
	}
	if err != nil {
		return err
	}
	return os.Rename(name, path)
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

var skill = strings.TrimSpace(`---
name: pathframe
description: Use Pathframe when explicitly requested or when continuing an active Pathframe workflow; offer it once for multi-step, durable, agreement-sensitive, or risky development work; never activate it for explanations, read-only review, brainstorming, trivial isolated edits, non-development operations, or explicit direct execution.
---

# Pathframe planning

Use only the typed `+"`pathframe_*`"+` MCP tools. Never construct Pathframe CLI commands.

1. Call `+"`pathframe_orient`"+` when Pathframe is explicitly requested or a workflow may already be active.
2. Brain classifies request facts, then calls `+"`pathframe_assess_request`"+`. Pathframe applies the deterministic activation rules.
3. For `+"`must_use`"+`, continue. For `+"`offer`"+`, ask once and wait. For `+"`must_not_use`"+`, proceed directly without Pathframe.
4. Create a change with `+"`pathframe_create_change`"+` and use `+"`pathframe_get_template`"+` for each required artifact. Brain authors the content; Pathframe validates it.
5. Call `+"`pathframe_validate_plan`"+` with `+"`human_approved: false`"+`. Fix all reported issues.
6. Ask the human to approve the validated plan explicitly. Only after approval, call the same tool with `+"`human_approved: true`"+`.
7. Call `+"`pathframe_get_next`"+` to resume or identify the next legal action. Use `+"`pathframe_recover`"+` only for its typed recovery actions.

Do not begin implementation while the plan is not ready. Do not invent unavailable delegation, verification, Claude, or compatibility operations.
`) + "\n"

const hooks = `{
  "description": "pathframe-generated:pathframe.integration/v1",
  "hooks": {
    "SessionStart": [
      {
        "matcher": "startup|resume",
        "hooks": [
          {
            "type": "mcp_tool",
            "server": "pathframe",
            "tool": "pathframe_orient",
            "input": {},
            "timeout": 10,
            "statusMessage": "Loading Pathframe orientation"
          }
        ]
      }
    ]
  }
}
`
