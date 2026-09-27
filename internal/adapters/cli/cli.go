package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/0xkhdr/pathframe/internal/app"
	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

// Run executes the Pathframe command and returns its process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	root := flag.NewFlagSet("pathframe", flag.ContinueOnError)
	root.SetOutput(stderr)
	help := root.Bool("help", false, "show help")
	version := root.Bool("version", false, "show version")
	jsonView := root.Bool("json", false, "render canonical JSON")
	change := root.String("change", "", "select a change")
	root.Usage = func() { usage(stderr) }
	if err := root.Parse(args); err != nil {
		return 2
	}
	if *help {
		usage(stdout)
		return 0
	}
	if *version {
		fmt.Fprintln(stdout, app.Version)
		return 0
	}

	command := "status"
	remaining := root.Args()
	if len(remaining) > 0 {
		command, remaining = remaining[0], remaining[1:]
	}
	if command == "help" {
		usage(stdout)
		return 0
	}
	commandFlags := flag.NewFlagSet(command, flag.ContinueOnError)
	commandFlags.SetOutput(stderr)
	commandJSON := commandFlags.Bool("json", *jsonView, "render canonical JSON")
	commandChange := commandFlags.String("change", *change, "select a change")
	commandMode := commandFlags.String("mode", "quick", "planning mode: quick, standard, or high-risk")
	commandArtifact := commandFlags.String("artifact", "task", "artifact kind")
	if err := commandFlags.Parse(remaining); err != nil || commandFlags.NArg() != 0 {
		return 2
	}

	service := app.Service{Dir: "."}
	if command == "new" {
		result, err := service.CreateChange(*commandChange, artifacts.Mode(*commandMode))
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "template" {
		instructions, template, err := service.Template(artifacts.Mode(*commandMode), *commandArtifact)
		if err != nil {
			fmt.Fprintf(stderr, "pathframe: %v\n", err)
			return 1
		}
		if *commandJSON {
			return renderJSON(stdout, stderr, instructions)
		}
		fmt.Fprint(stdout, template)
		return 0
	}
	if command == "check" {
		result, err := service.Check(*commandChange)
		return renderArtifactResult(stdout, stderr, result, *commandJSON, false, err)
	}
	if command == "approve" {
		result, err := service.Approve(*commandChange)
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	var result workflow.Result
	var err error
	switch command {
	case "status":
		result, err = service.Orient(*commandChange)
	case "next":
		result, err = service.GetNext(*commandChange)
	case "pause":
		result, err = service.Pause(*commandChange)
	case "resume":
		result, err = service.Resume(*commandChange)
	case "replan":
		result, err = service.Replan(*commandChange)
	case "cancel":
		result, err = service.Cancel(*commandChange)
	default:
		fmt.Fprintf(stderr, "pathframe: unknown command %q\n", command)
		return 2
	}
	if renderErr := render(stdout, result, *commandJSON); renderErr != nil {
		fmt.Fprintf(stderr, "pathframe: %v\n", renderErr)
		return 1
	}
	if err != nil {
		fmt.Fprintf(stderr, "pathframe: %v\n", err)
		return 1
	}
	return 0
}

func renderArtifactResult(w, errw io.Writer, value any, jsonView, allowInvalid bool, operationErr error) int {
	if jsonView {
		if code := renderJSON(w, errw, value); code != 0 {
			return code
		}
	} else {
		switch result := value.(type) {
		case artifacts.CheckResult:
			fmt.Fprintf(w, "Change: %s\nMode: %s\nValid: %t\n", result.Change, result.Mode, result.Valid)
			if result.Identity != "" {
				fmt.Fprintf(w, "Plan identity: %s\n", result.Identity)
			}
			for _, issue := range result.Issues {
				fmt.Fprintf(w, "Issue: %s: %s (%s)\n", issue.Path, issue.Message, issue.Code)
			}
		case artifacts.ApprovalResult:
			fmt.Fprintf(w, "Change: %s\nMode: %s\nPhase: %s\nPlan identity: %s\n", result.Change, result.Mode, result.Phase, result.Identity)
		}
	}
	if operationErr != nil {
		fmt.Fprintf(errw, "pathframe: %v\n", operationErr)
		return 1
	}
	if result, ok := value.(artifacts.CheckResult); ok && !result.Valid && !allowInvalid {
		return 1
	}
	return 0
}

func renderJSON(w, errw io.Writer, value any) int {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		fmt.Fprintf(errw, "pathframe: %v\n", err)
		return 1
	}
	return 0
}

func render(w io.Writer, result workflow.Result, jsonView bool) error {
	if jsonView {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	fmt.Fprintln(w, "Pathframe")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Configured: %t\n", result.Configured)
	if result.Change == "" {
		fmt.Fprintln(w, "Active change: none")
	} else {
		fmt.Fprintf(w, "Active change: %s\n", result.Change)
	}
	if result.Phase == "" {
		fmt.Fprintln(w, "Phase: none")
	} else {
		fmt.Fprintf(w, "Phase: %s\n", result.Phase)
	}
	fmt.Fprintf(w, "Progress: %d/%d tasks complete\n", result.Progress.Completed, result.Progress.Total)
	if len(result.Blockers) == 0 {
		fmt.Fprintln(w, "Blocked: no")
	} else {
		fmt.Fprintf(w, "Blocked: %s (%s)\n", result.Blockers[0].Message, result.Blockers[0].Code)
	}
	if result.Recommended.Action == "" {
		fmt.Fprintln(w, "Recommended next action: none")
	} else {
		fmt.Fprintf(w, "Recommended next action: %s\n", result.Recommended.Action)
	}
	if len(result.Alternatives) > 0 {
		values := make([]string, len(result.Alternatives))
		for i, action := range result.Alternatives {
			values[i] = string(action.Action)
		}
		fmt.Fprintf(w, "Other legal actions: %s\n", strings.Join(values, ", "))
	}
	fmt.Fprintf(w, "Human action required: %t\n", result.HumanRequired)
	for _, diagnostic := range result.Diagnostics {
		fmt.Fprintf(w, "Diagnostic: %s\n", diagnostic)
	}
	return nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: pathframe [--json] [--change ID] [status|next|new|template|check|approve|pause|resume|replan|cancel]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "No command and status both show project orientation.")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  --change ID  select a change when more than one exists")
	fmt.Fprintln(w, "  --json       render pathframe.workflow/v1 JSON")
	fmt.Fprintln(w, "  --mode MODE  quick, standard, or high-risk (new/template)")
	fmt.Fprintln(w, "  --artifact K artifact kind (template)")
	fmt.Fprintln(w, "  --help       show help")
	fmt.Fprintln(w, "  --version    show version")
}
