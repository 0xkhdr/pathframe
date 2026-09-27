package verification

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultTimeout  = 2 * time.Minute
	DefaultMaxBytes = 64 * 1024
)

type Command struct {
	Argv     []string      `json:"argv"`
	Workdir  string        `json:"workdir"`
	Timeout  time.Duration `json:"-"`
	MaxBytes int           `json:"max_bytes"`
}

func (c Command) Validate(root string) (string, error) {
	if len(c.Argv) == 0 || strings.TrimSpace(c.Argv[0]) == "" {
		return "", fmt.Errorf("verification argv must not be empty")
	}
	if c.Timeout <= 0 {
		return "", fmt.Errorf("verification timeout must be positive")
	}
	if c.MaxBytes <= 0 {
		return "", fmt.Errorf("verification output limit must be positive")
	}
	workdir := c.Workdir
	if workdir == "" {
		workdir = "."
	}
	if filepath.IsAbs(workdir) {
		return "", fmt.Errorf("verification workdir must be project-relative")
	}
	clean := filepath.Clean(workdir)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("verification workdir escapes the project")
	}
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	realWorkdir, err := filepath.EvalSymlinks(filepath.Join(realRoot, clean))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(realRoot, realWorkdir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("verification workdir escapes the project through a symlink")
	}
	return realWorkdir, nil
}
