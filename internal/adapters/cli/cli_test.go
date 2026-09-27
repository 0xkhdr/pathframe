package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "orientation", want: "Pathframe\n\nStatus: not initialized\nNext: workflow commands are not available in Stage 0.\n"},
		{name: "help", args: []string{"--help"}, want: "Usage: pathframe [--help] [--version]\n"},
		{name: "version", args: []string{"--version"}, want: "dev\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := Run(test.args, &stdout, &stderr); code != 0 {
				t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
			}
			if !strings.HasPrefix(stdout.String(), test.want) {
				t.Fatalf("Run() stdout = %q, want prefix %q", stdout.String(), test.want)
			}
			if stderr.Len() != 0 {
				t.Fatalf("Run() stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"start"}, &stdout, &stderr); code != 2 {
		t.Fatalf("Run() code = %d, want 2", code)
	}
	if got := stderr.String(); got != "pathframe: unknown command \"start\"\n" {
		t.Fatalf("Run() stderr = %q", got)
	}
}
