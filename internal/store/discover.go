package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ErrProjectNotFound = errors.New(".pathframe project not found")

func Discover(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	if dir, err = filepath.EvalSymlinks(dir); err != nil {
		return "", err
	}
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		dir = filepath.Dir(dir)
	}
	for {
		managed := filepath.Join(dir, ".pathframe")
		info, err := os.Lstat(managed)
		if err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrProjectNotFound
		}
		dir = parent
	}
}

func Changes(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, ".pathframe", "changes"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var changes []string
	for _, entry := range entries {
		if entry.IsDir() {
			changes = append(changes, entry.Name())
		}
	}
	sort.Strings(changes)
	return changes, nil
}

func ChangeDir(root, change string) (string, error) {
	if change == "" || change == "." || change == ".." || filepath.Base(change) != change || strings.ContainsAny(change, `/\\`) {
		return "", fmt.Errorf("invalid change identifier %q", change)
	}
	base := filepath.Join(root, ".pathframe", "changes")
	dir := filepath.Join(base, change)
	resolvedBase, err := filepath.EvalSymlinks(base)
	if err != nil {
		return "", err
	}
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	if !contained(resolvedRoot, resolvedBase) {
		return "", errors.New("managed changes directory escapes project")
	}
	rel, err := filepath.Rel(resolvedBase, resolvedDir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("change path escapes project: %q", change)
	}
	return resolvedDir, nil
}

func contained(base, path string) bool {
	rel, err := filepath.Rel(base, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
