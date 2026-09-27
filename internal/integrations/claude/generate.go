package claude

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	IntegrationSchema  = "pathframe.integration/v1"
	IntegrationVersion = "1.0.0"
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
	assets, manifest, err := generated(executable, sessionOrientation)
	if err != nil {
		return InstallResult{}, err
	}
	old, _ := readManifest(root)
	for name, body := range assets {
		current, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
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
		if _, keep := assets[name]; keep {
			continue
		}
		current, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if readErr == nil && digest(current) != want {
			return InstallResult{}, fmt.Errorf("refusing to remove modified generated file %s; move it aside first", name)
		}
		if readErr != nil && !os.IsNotExist(readErr) {
			return InstallResult{}, readErr
		}
	}
	changed := []string{}
	for name := range old.Files {
		if _, keep := assets[name]; !keep {
			if err := os.Remove(filepath.Join(root, filepath.FromSlash(name))); err == nil {
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
	manifestName := ".pathframe/integrations/claude/manifest.json"
	manifestPath := filepath.Join(root, filepath.FromSlash(manifestName))
	current, _ := os.ReadFile(manifestPath)
	if string(current) != string(data) {
		if err := write(manifestPath, data); err != nil {
			return InstallResult{}, err
		}
		changed = append(changed, manifestName)
	}
	sort.Strings(changed)
	return InstallResult{Schema: IntegrationSchema, Version: IntegrationVersion, Changed: changed}, nil
}

func generated(executable string, sessionOrientation bool) (map[string]string, Manifest, error) {
	executable, err := filepath.Abs(executable)
	if err != nil {
		return nil, Manifest{}, err
	}
	mcp := map[string]any{"mcpServers": map[string]any{"pathframe": map[string]any{"type": "stdio", "command": executable, "args": []string{"mcp"}}}}
	mcpData, _ := json.MarshalIndent(mcp, "", "  ")
	assets := map[string]string{
		".claude/skills/pathframe/SKILL.md": skill,
		".claude/commands/pathframe.md":     command,
		".mcp.json":                         string(mcpData) + "\n",
	}
	capabilities := []string{"project_skill", "slash_command", "stdio_mcp", "typed_planning_tools", "task_packet_preview"}
	if sessionOrientation {
		settings := map[string]any{
			"hooks": map[string]any{
				"SessionStart": []any{map[string]any{
					"matcher": "startup|resume|clear|compact",
					"hooks": []any{map[string]any{
						"type":          "command",
						"command":       fmt.Sprintf("%q status --json", executable),
						"timeout":       10,
						"statusMessage": "Loading Pathframe orientation",
					}},
				}},
			},
		}
		data, _ := json.MarshalIndent(settings, "", "  ")
		assets[".claude/settings.json"] = string(data) + "\n"
		capabilities = append(capabilities, "session_start_orientation")
	}
	return assets, Manifest{Schema: IntegrationSchema, Version: IntegrationVersion, Host: "claude-code", Transport: "stdio", Capabilities: capabilities}, nil
}

func readManifest(root string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(root, ".pathframe", "integrations", "claude", "manifest.json"))
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

Use only the typed MCP tools from the `+"`pathframe`"+` server. Never construct Pathframe CLI commands.

1. Call `+"`pathframe_orient`"+` when Pathframe is explicitly requested or a workflow may already be active.
2. Brain classifies request facts, then calls `+"`pathframe_assess_request`"+`. Pathframe returns `+"`must_use`"+`, `+"`offer`"+`, or `+"`must_not_use`"+`.
3. Continue for `+"`must_use`"+`; ask once and wait for `+"`offer`"+`; proceed directly for `+"`must_not_use`"+`.
4. Create a change and request every required template through typed tools. Brain authors content; Pathframe validates it.
5. Validate with `+"`human_approved: false`"+` and fix every issue.
6. Ask for explicit human approval, then validate with `+"`human_approved: true`"+`.
7. Use `+"`pathframe_get_next`"+` to resume. For a ready change, use `+"`pathframe_prepare_delegation`"+` to inspect the bounded task packet and sequential frontier; it never launches a worker. Use `+"`pathframe_recover`"+` only for typed recovery actions.

Do not implement while the plan is not ready. Do not treat packet preview as launch authority. Do not invent worker launch, result submission, verification execution, Pinky, parallel, or compatibility operations.
`) + "\n"

const command = `---
description: Start or resume Pathframe planning through typed MCP operations
---

Use the project Pathframe skill for this request. Begin with ` + "`pathframe_orient`" + ` and follow its activation and approval rules. Use only typed MCP operations; never construct Pathframe CLI commands.
`
