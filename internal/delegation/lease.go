package delegation

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
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
	Schema    string    `json:"schema"`
	LeaseID   string    `json:"lease_id"`
	State     string    `json:"state"`
	Timestamp time.Time `json:"timestamp"`
	Result    *Result   `json:"result,omitempty"`
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

func Release(runsDir string, lease Lease) error {
	active, ok, err := ActiveLease(runsDir)
	if err != nil || !ok || active.ID != lease.ID {
		return fmt.Errorf("active lease does not match submission")
	}
	return os.Remove(filepath.Join(runsDir, "active.json"))
}
