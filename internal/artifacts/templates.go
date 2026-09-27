package artifacts

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed templates/*.md
var templateFiles embed.FS

func Template(mode Mode, kind string) (string, error) {
	_, err := ArtifactInstructions(mode, kind)
	if err != nil {
		return "", err
	}
	data, err := templateFiles.ReadFile("templates/" + kind + ".md")
	if err != nil {
		return "", fmt.Errorf("template %s: %w", kind, err)
	}
	policy := "<brain-or-delegated>"
	if mode == Quick {
		policy = "brain"
	}
	return strings.ReplaceAll(string(data), "{{execution_policy}}", policy), nil
}
