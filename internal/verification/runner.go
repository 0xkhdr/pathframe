package verification

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

type Result struct {
	Argv            []string `json:"argv"`
	Workdir         string   `json:"workdir"`
	ExitCode        int      `json:"exit_code"`
	DurationMS      int64    `json:"duration_ms"`
	TimedOut        bool     `json:"timed_out"`
	Interrupted     bool     `json:"interrupted"`
	Stdout          string   `json:"stdout"`
	Stderr          string   `json:"stderr"`
	StdoutTruncated bool     `json:"stdout_truncated"`
	StderrTruncated bool     `json:"stderr_truncated"`
}

func Run(ctx context.Context, root string, command Command) Result {
	result := Result{Argv: append([]string(nil), command.Argv...), Workdir: command.Workdir, ExitCode: -1}
	dir, err := command.Validate(root)
	if err != nil {
		result.Stderr = err.Error()
		return result
	}
	if result.Workdir == "" {
		result.Workdir = "."
	}
	ctx, cancel := context.WithTimeout(ctx, command.Timeout)
	defer cancel()
	stdout, stderr := &boundedWriter{limit: command.MaxBytes}, &boundedWriter{limit: command.MaxBytes}
	cmd := exec.CommandContext(ctx, command.Argv[0], command.Argv[1:]...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, stdout, stderr
	started := time.Now()
	err = cmd.Run()
	result.DurationMS = time.Since(started).Milliseconds()
	result.Stdout, result.Stderr = stdout.String(), stderr.String()
	result.StdoutTruncated, result.StderrTruncated = stdout.truncated, stderr.truncated
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}
	result.TimedOut = errors.Is(ctx.Err(), context.DeadlineExceeded)
	result.Interrupted = errors.Is(ctx.Err(), context.Canceled)
	if err != nil && result.Stderr == "" && !result.TimedOut && !result.Interrupted {
		result.Stderr = err.Error()
	}
	return result
}

type boundedWriter struct {
	data      []byte
	limit     int
	truncated bool
}

func (w *boundedWriter) Write(p []byte) (int, error) {
	remaining := w.limit - len(w.data)
	if remaining > 0 {
		if remaining > len(p) {
			remaining = len(p)
		}
		w.data = append(w.data, p[:remaining]...)
	}
	if remaining < len(p) {
		w.truncated = true
	}
	return len(p), nil
}

func (w *boundedWriter) String() string { return string(w.data) }
