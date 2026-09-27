package recovery

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

func DiagnoseProjection(changeDir string) ([]Diagnosis, error) {
	replay, err := store.ReadJournal(filepath.Join(changeDir, "history.jsonl"))
	if err != nil {
		return []Diagnosis{{Code: "journal_invalid", Severity: SeverityError, Subject: "history.jsonl", Evidence: err.Error(), HumanRequired: true}}, nil
	}
	var out []Diagnosis
	for _, code := range replay.Diagnostics {
		if code == workflow.ReasonIncompleteTail {
			out = append(out, Diagnosis{Code: string(code), Severity: SeverityWarning, Subject: "history.jsonl", Evidence: "trailing bytes do not form a complete event", Repair: "truncate_incomplete_journal_tail"})
		}
	}
	want := store.Projection{Schema: workflow.Schema, Events: replay.Events, State: replay.State}
	got, loadErr := store.LoadProjection(filepath.Join(changeDir, "state.json"))
	if loadErr != nil || got != want {
		evidence := "projection differs from deterministic replay"
		if loadErr != nil && !errors.Is(loadErr, os.ErrNotExist) {
			evidence = loadErr.Error()
		}
		out = append(out, Diagnosis{Code: "corrupt_projection", Severity: SeverityWarning, Subject: "state.json", Evidence: evidence, Repair: "rebuild_projection"})
	}
	return out, nil
}

func RepairProjection(changeDir string) error {
	_, err := store.ReplayAndRepair(changeDir)
	return err
}
