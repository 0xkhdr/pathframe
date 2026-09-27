package codex

import (
	"fmt"
	"os"
	"path/filepath"
)

type Diagnostic struct {
	Path     string `json:"path"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Recovery string `json:"recovery"`
}

type DoctorResult struct {
	Schema      string       `json:"schema"`
	Healthy     bool         `json:"healthy"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

func Doctor(root string) DoctorResult {
	result := DoctorResult{Schema: IntegrationSchema, Healthy: true, Diagnostics: []Diagnostic{}}
	manifest, err := readManifest(root)
	if err != nil {
		return DoctorResult{Schema: IntegrationSchema, Healthy: false, Diagnostics: []Diagnostic{{Path: ".pathframe/integrations/codex/manifest.json", Code: "missing_or_invalid_manifest", Message: err.Error(), Recovery: "run pathframe codex-install"}}}
	}
	if manifest.Schema != IntegrationSchema || manifest.Version != IntegrationVersion || manifest.Host != "codex" || manifest.Transport != "stdio" {
		result.Healthy = false
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Path: ".pathframe/integrations/codex/manifest.json", Code: "incompatible_manifest", Message: fmt.Sprintf("found schema=%q version=%q host=%q transport=%q", manifest.Schema, manifest.Version, manifest.Host, manifest.Transport), Recovery: "run pathframe codex-install to update generated assets"})
	}
	for _, name := range []string{".agents/skills/pathframe/SKILL.md", ".codex/config.toml"} {
		if manifest.Files[name] == "" {
			result.Healthy = false
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Path: name, Code: "incomplete_manifest", Message: "required generated asset is not declared", Recovery: "run pathframe codex-install"})
		}
	}
	for name, want := range manifest.Files {
		data, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if readErr != nil || digest(data) != want {
			result.Healthy = false
			message := "generated file differs from its manifest"
			if readErr != nil {
				message = readErr.Error()
			}
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Path: name, Code: "missing_or_modified_asset", Message: message, Recovery: "restore the generated file or move custom content aside, then run pathframe codex-install"})
		}
	}
	return result
}
