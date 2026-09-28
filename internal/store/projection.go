package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/workflow"
)

type Projection struct {
	Schema string         `json:"schema"`
	Events int            `json:"events"`
	State  workflow.State `json:"state"`
}

func LoadProjection(path string) (Projection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Projection{}, err
	}
	var projection Projection
	if err := json.Unmarshal(data, &projection); err != nil {
		return Projection{}, err
	}
	if projection.Schema != workflow.Schema || projection.Events < 1 || !workflow.ValidPhase(projection.State.Phase) {
		return Projection{}, errors.New("invalid workflow projection")
	}
	return projection, nil
}

func WriteProjection(path string, projection Projection) error {
	data, err := json.MarshalIndent(projection, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".state-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err = tmp.Chmod(0o600); err == nil {
		_, err = tmp.Write(data)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func ReplayAndRepair(changeDir string) (Replay, error) {
	return replayAndRepair(changeDir, true)
}

func RefreshProjection(changeDir string) (Replay, error) {
	return replayAndRepair(changeDir, false)
}

func replayAndRepair(changeDir string, diagnoseMismatch bool) (Replay, error) {
	journalPath := filepath.Join(changeDir, "history.jsonl")
	replay, err := ReadJournal(journalPath)
	if err != nil {
		return Replay{}, err
	}
	for _, diagnostic := range replay.Diagnostics {
		if diagnostic == workflow.ReasonIncompleteTail {
			if err := os.Truncate(journalPath, replay.CompleteBytes); err != nil {
				return Replay{}, err
			}
		}
	}
	want := Projection{Schema: workflow.Schema, Events: replay.Events, State: replay.State}
	got, loadErr := LoadProjection(filepath.Join(changeDir, "state.json"))
	if loadErr != nil || got != want {
		if diagnoseMismatch && ((loadErr != nil && !errors.Is(loadErr, os.ErrNotExist)) || (loadErr == nil && got != want)) {
			replay.Diagnostics = append(replay.Diagnostics, workflow.ReasonCorruptProjection)
		}
		if err := WriteProjection(filepath.Join(changeDir, "state.json"), want); err != nil {
			return Replay{}, err
		}
	}
	return replay, nil
}
