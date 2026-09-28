package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func (s Service) CreateChange(change string, mode artifacts.Mode) (artifacts.CheckResult, error) {
	if !mode.Valid() {
		return artifacts.CheckResult{}, fmt.Errorf("invalid mode %q", mode)
	}
	if change == "" || filepath.Base(change) != change || strings.ContainsAny(change, `/\\`) || change == "." || change == ".." {
		return artifacts.CheckResult{}, fmt.Errorf("invalid change identifier %q", change)
	}
	start := s.Dir
	if start == "" {
		start = "."
	}
	root, err := store.Discover(start)
	if errors.Is(err, store.ErrProjectNotFound) {
		root, err = filepath.Abs(start)
		if err != nil {
			return artifacts.CheckResult{}, err
		}
		managed := filepath.Join(root, ".pathframe")
		if info, statErr := os.Lstat(managed); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return artifacts.CheckResult{}, errors.New("managed directory must not be a symlink")
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return artifacts.CheckResult{}, statErr
		}
		if err = os.MkdirAll(filepath.Join(root, ".pathframe", "changes"), 0o755); err != nil {
			return artifacts.CheckResult{}, err
		}
	} else if err != nil {
		return artifacts.CheckResult{}, err
	}
	dir := filepath.Join(root, ".pathframe", "changes", change)
	existing, loadErr := artifacts.LoadChange(dir)
	if loadErr == nil && modeRank(mode) < modeRank(existing.Mode) {
		return artifacts.CheckResult{}, fmt.Errorf("mode cannot be lowered from %s to %s", existing.Mode, mode)
	}
	if loadErr != nil && !errors.Is(loadErr, os.ErrNotExist) {
		return artifacts.CheckResult{}, loadErr
	}
	if err := os.MkdirAll(filepath.Join(dir, "tasks"), 0o755); err != nil {
		return artifacts.CheckResult{}, err
	}
	if loadErr != nil || existing.Mode != mode {
		changeFile := fmt.Sprintf("schema: %s\nprofile: %s\nid: %s\nmode: %s\n", artifacts.ChangeSchema, artifacts.Profile, change, mode)
		if err := atomicWrite(filepath.Join(dir, "change.yaml"), []byte(changeFile)); err != nil {
			return artifacts.CheckResult{}, err
		}
	}
	for _, kind := range artifacts.RequiredArtifacts(mode) {
		path := filepath.Join(dir, kind+".md")
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return artifacts.CheckResult{}, err
		}
		body, err := artifacts.Template(mode, kind)
		if err != nil {
			return artifacts.CheckResult{}, err
		}
		if err := atomicWrite(path, []byte(body)); err != nil {
			return artifacts.CheckResult{}, err
		}
	}
	taskPath := filepath.Join(dir, "tasks", "T1.md")
	if _, err := os.Stat(taskPath); errors.Is(err, os.ErrNotExist) {
		body, templateErr := artifacts.Template(mode, "task")
		if templateErr != nil {
			return artifacts.CheckResult{}, templateErr
		}
		if err := atomicWrite(taskPath, []byte(body)); err != nil {
			return artifacts.CheckResult{}, err
		}
	}
	journal := filepath.Join(dir, "history.jsonl")
	if _, err := os.Stat(journal); errors.Is(err, os.ErrNotExist) {
		initial := store.Event{Schema: store.EventSchema, ID: change + "-initial", Timestamp: now(), Actor: workflow.ActorSystem, Action: workflow.ActionInitialize, Change: change, Target: workflow.PhaseExploring}
		if err := store.Append(journal, initial); err != nil {
			return artifacts.CheckResult{}, err
		}
		before := workflow.State{Change: change, Phase: workflow.PhaseExploring}
		after, _ := workflow.Apply(before, workflow.ActionPlan, workflow.ActorBrain)
		event, err := store.NewEvent(before, after, workflow.ActionPlan, workflow.ActorBrain, "planning artifacts created")
		if err != nil {
			return artifacts.CheckResult{}, err
		}
		if err := store.Append(journal, event); err != nil {
			return artifacts.CheckResult{}, err
		}
	}
	if _, err := store.RefreshProjection(dir); err != nil {
		return artifacts.CheckResult{}, err
	}
	return s.Check(change)
}

func modeRank(mode artifacts.Mode) int {
	if mode == artifacts.Quick {
		return 1
	}
	if mode == artifacts.Standard {
		return 2
	}
	return 3
}

func atomicWrite(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".artifact-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = tmp.Chmod(0o600); err == nil {
		_, err = tmp.Write(data)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(name, path)
}
