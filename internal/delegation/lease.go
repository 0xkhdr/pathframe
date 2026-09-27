package delegation

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/0xkhdr/pathframe/internal/verification"
)

const RunSchema = "pathframe.run/v1"

type Lease struct {
	Schema    string    `json:"schema"`
	ID        string    `json:"id"`
	Change    string    `json:"change"`
	Task      string    `json:"task"`
	Host      string    `json:"host"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type RunEvent struct {
	Schema       string        `json:"schema"`
	LeaseID      string        `json:"lease_id"`
	State        string        `json:"state"`
	Timestamp    time.Time     `json:"timestamp"`
	Result       *Result       `json:"result,omitempty"`
	Verification *Verification `json:"verification,omitempty"`
	Review       *Review       `json:"review,omitempty"`
}

type Verification struct {
	Commands        []verification.Result `json:"commands"`
	ContentIdentity string                `json:"content_identity"`
	Passed          bool                  `json:"passed"`
	ScopeViolations []string              `json:"scope_violations"`
}

type Review struct {
	Decision     string   `json:"decision"`
	Reason       string   `json:"reason"`
	Recovery     []string `json:"recovery"`
	TaskIdentity string   `json:"task_identity,omitempty"`
}

func ActiveLease(runsDir string) (Lease, bool, error) {
	data, err := os.ReadFile(filepath.Join(runsDir, "active.json"))
	if errors.Is(err, os.ErrNotExist) {
		return Lease{}, false, nil
	}
	if err != nil {
		return Lease{}, false, err
	}
	var lease Lease
	if err := json.Unmarshal(data, &lease); err != nil || lease.Schema != RunSchema || lease.ID == "" || lease.State != "active" {
		return Lease{}, false, fmt.Errorf("invalid active task lease")
	}
	return lease, true, nil
}

func Acquire(runsDir, change, task, host string, at time.Time) (Lease, error) {
	if _, active, err := ActiveLease(runsDir); err != nil || active {
		if err != nil {
			return Lease{}, err
		}
		return Lease{}, fmt.Errorf("a task lease is already active")
	}
	if err := os.MkdirAll(runsDir, 0o700); err != nil {
		return Lease{}, err
	}
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return Lease{}, err
	}
	lease := Lease{Schema: RunSchema, ID: hex.EncodeToString(id), Change: change, Task: task, Host: host, State: "active", CreatedAt: at.UTC()}
	data, _ := json.MarshalIndent(lease, "", "  ")
	data = append(data, '\n')
	file, err := os.OpenFile(filepath.Join(runsDir, "active.json"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return Lease{}, err
	}
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(filepath.Join(runsDir, "active.json"))
		return Lease{}, err
	}
	if err := AppendRun(runsDir, lease.ID, RunEvent{Schema: RunSchema, LeaseID: lease.ID, State: "leased", Timestamp: at.UTC()}); err != nil {
		_ = os.Remove(filepath.Join(runsDir, "active.json"))
		return Lease{}, err
	}
	return lease, nil
}

func AppendRun(runsDir, id string, event RunEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(runsDir, id+".jsonl"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err = file.Write(append(data, '\n')); err == nil {
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}

func TaskRun(runsDir, task string) (string, []RunEvent, error) {
	paths, err := filepath.Glob(filepath.Join(runsDir, "*.jsonl"))
	if err != nil {
		return "", nil, err
	}
	sort.Strings(paths)
	var selected []RunEvent
	var selectedID string
	var selectedAt time.Time
	for _, path := range paths {
		events, err := readRun(path)
		if err != nil {
			return "", nil, err
		}
		for _, event := range events {
			if event.Result != nil && event.Result.Task == task && (selectedAt.IsZero() || event.Timestamp.After(selectedAt)) {
				selectedID, selected, selectedAt = event.LeaseID, events, event.Timestamp
			}
		}
	}
	if selectedID != "" {
		return selectedID, selected, nil
	}
	return "", nil, fmt.Errorf("no submitted run exists for task %s", task)
}

func CompletedTasks(runsDir string) ([]string, error) {
	return CompletedTasksMatching(runsDir, nil)
}

// CompletedTasksMatching excludes completions whose recorded task contract no
// longer matches. A nil identity map retains the Stage 7 behavior.
func CompletedTasksMatching(runsDir string, identities map[string]string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(runsDir, "*.jsonl"))
	if err != nil {
		return nil, err
	}
	done := map[string]bool{}
	for _, path := range paths {
		events, err := readRun(path)
		if err != nil {
			return nil, err
		}
		var task, identity string
		for _, event := range events {
			if event.Result != nil {
				task = event.Result.Task
			}
			if event.Review != nil {
				identity = event.Review.TaskIdentity
			}
			if task != "" && event.State == "completed" && (identities == nil || identity != "" && identities[task] == identity) {
				done[task] = true
			}
		}
	}
	result := make([]string, 0, len(done))
	for task := range done {
		result = append(result, task)
	}
	sort.Strings(result)
	return result, nil
}

func RunStates(runsDir string) (map[string]string, error) {
	paths, err := filepath.Glob(filepath.Join(runsDir, "*.jsonl"))
	if err != nil {
		return nil, err
	}
	states := make(map[string]string, len(paths))
	for _, path := range paths {
		events, err := readRun(path)
		if err != nil {
			return nil, err
		}
		if len(events) > 0 {
			states[events[len(events)-1].LeaseID] = events[len(events)-1].State
		}
	}
	return states, nil
}

func readRun(path string) ([]RunEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var events []RunEvent
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event RunEvent
		if err := json.Unmarshal(line, &event); err != nil || event.Schema != RunSchema || event.LeaseID == "" {
			return nil, fmt.Errorf("invalid run record %s", filepath.Base(path))
		}
		events = append(events, event)
	}
	return events, nil
}

func Release(runsDir string, lease Lease) error {
	active, ok, err := ActiveLease(runsDir)
	if err != nil || !ok || active.ID != lease.ID {
		return fmt.Errorf("active lease does not match submission")
	}
	return os.Remove(filepath.Join(runsDir, "active.json"))
}
