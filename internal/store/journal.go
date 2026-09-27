package store

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/0xkhdr/pathframe/internal/workflow"
)

const EventSchema = "pathframe.transition/v1"

type Event struct {
	Schema    string          `json:"schema"`
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Actor     workflow.Actor  `json:"actor"`
	Action    workflow.Action `json:"action"`
	Change    string          `json:"change"`
	Source    workflow.Phase  `json:"source,omitempty"`
	Target    workflow.Phase  `json:"target"`
	Resume    workflow.Phase  `json:"resume_phase,omitempty"`
	Reason    string          `json:"reason,omitempty"`
}

type Replay struct {
	State         workflow.State
	Events        int
	CompleteBytes int64
	Diagnostics   []workflow.ReasonCode
}

func NewEvent(before, after workflow.State, action workflow.Action, actor workflow.Actor, reason string) (Event, error) {
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return Event{}, err
	}
	return Event{Schema: EventSchema, ID: hex.EncodeToString(id), Timestamp: time.Now().UTC(), Actor: actor, Action: action, Change: after.Change, Source: before.Phase, Target: after.Phase, Resume: after.ResumePhase, Reason: reason}, nil
}

func Append(path string, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
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

func ReadJournal(path string) (Replay, error) {
	file, err := os.Open(path)
	if err != nil {
		return Replay{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	result := Replay{}
	seen := map[string]bool{}
	for {
		line, readErr := reader.ReadBytes('\n')
		if errors.Is(readErr, io.EOF) && len(line) > 0 {
			result.Diagnostics = append(result.Diagnostics, workflow.ReasonIncompleteTail)
			break
		}
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return Replay{}, readErr
		}
		if len(line) > 1 {
			var event Event
			if err := json.Unmarshal(line, &event); err != nil {
				return Replay{}, fmt.Errorf("journal record %d: %w", result.Events+1, err)
			}
			if event.Schema != EventSchema || event.ID == "" || seen[event.ID] || event.Timestamp.IsZero() || event.Change == "" {
				return Replay{}, fmt.Errorf("journal record %d is invalid", result.Events+1)
			}
			seen[event.ID] = true
			if result.Events == 0 {
				if event.Action != workflow.ActionInitialize || event.Actor != workflow.ActorSystem || event.Source != "" || event.Target != workflow.PhaseExploring {
					return Replay{}, fmt.Errorf("first journal record must initialize exploring")
				}
				result.State = workflow.State{Change: event.Change, Phase: event.Target}
			} else {
				if event.Change != result.State.Change || event.Source != result.State.Phase {
					return Replay{}, fmt.Errorf("journal record %d does not continue current state", result.Events+1)
				}
				next, err := workflow.Apply(result.State, event.Action, event.Actor)
				if err != nil || next.Phase != event.Target || next.ResumePhase != event.Resume {
					return Replay{}, fmt.Errorf("journal record %d contains an illegal transition", result.Events+1)
				}
				result.State = next
			}
			result.Events++
			result.CompleteBytes += int64(len(line))
		}
		if errors.Is(readErr, io.EOF) {
			break
		}
	}
	if result.Events == 0 {
		return Replay{}, errors.New("journal contains no complete events")
	}
	return result, nil
}
