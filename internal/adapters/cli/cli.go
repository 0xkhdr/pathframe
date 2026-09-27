package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/0xkhdr/pathframe/internal/app"
	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/integrations/claude"
	"github.com/0xkhdr/pathframe/internal/integrations/codex"
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
	commandTask := commandFlags.String("task", "", "task identifier")
	commandHost := commandFlags.String("host", "", "delegation host: codex or claude-code")
	resultFile := commandFlags.String("result-file", "", "pathframe.task-result/v1 JSON file")
	leaseID := commandFlags.String("lease-id", "", "exact active lease identifier")
	budgetBytes := commandFlags.Int("budget-bytes", 0, "context byte budget")
	timeoutMS := commandFlags.Int64("timeout-ms", 0, "required verification timeout in milliseconds")
	maxBytes := commandFlags.Int("max-bytes", 0, "per-stream verification output limit")
	workdir := commandFlags.String("workdir", ".", "project-relative verification working directory")
	changedFiles := commandFlags.String("changed-files", "", "comma-separated changed files for direct Brain verification")
	reason := commandFlags.String("reason", "", "semantic review reason")
	keepScope := commandFlags.Bool("keep-scope-violations", false, "explicitly accept advisory scope violations")
	sessionOrientation := commandFlags.Bool("session-orientation", false, "install optional host session-start orientation")
	if err := commandFlags.Parse(remaining); err != nil || commandFlags.NArg() != 0 {
		return 2
	}

	service := app.Service{Dir: "."}
	if command == "codex-install" {
		executable, err := os.Executable()
		if err != nil {
			fmt.Fprintf(stderr, "pathframe: %v\n", err)
			return 1
		}
		result, err := codex.Install(".", executable, *sessionOrientation)
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "codex-doctor" {
		result := codex.Doctor(".")
		if *commandJSON {
			if code := renderJSON(stdout, stderr, result); code != 0 {
				return code
			}
			if !result.Healthy {
				return 1
			}
			return 0
		}
		fmt.Fprintf(stdout, "Codex integration healthy: %t\n", result.Healthy)
		for _, diagnostic := range result.Diagnostics {
			fmt.Fprintf(stdout, "%s: %s (%s); recovery: %s\n", diagnostic.Path, diagnostic.Message, diagnostic.Code, diagnostic.Recovery)
		}
		if !result.Healthy {
			return 1
		}
		return 0
	}
	if command == "claude-install" {
		executable, err := os.Executable()
		if err != nil {
			fmt.Fprintf(stderr, "pathframe: %v\n", err)
			return 1
		}
		result, err := claude.Install(".", executable, *sessionOrientation)
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "claude-doctor" {
		result := claude.Doctor(".")
		if *commandJSON {
			if code := renderJSON(stdout, stderr, result); code != 0 {
				return code
			}
			if !result.Healthy {
				return 1
			}
			return 0
		}
		fmt.Fprintf(stdout, "Claude Code integration healthy: %t\n", result.Healthy)
		for _, diagnostic := range result.Diagnostics {
			fmt.Fprintf(stdout, "%s: %s (%s); recovery: %s\n", diagnostic.Path, diagnostic.Message, diagnostic.Code, diagnostic.Recovery)
		}
		if !result.Healthy {
			return 1
		}
		return 0
	}
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
	if command == "packet" {
		result, err := service.PrepareDelegation(app.PrepareInput{Change: *commandChange, Task: *commandTask, BudgetBytes: *budgetBytes, Host: *commandHost})
		if err != nil {
			fmt.Fprintf(stderr, "pathframe: %v\n", err)
			return 1
		}
		if *commandJSON {
			if code := renderJSON(stdout, stderr, result); code != 0 {
				return code
			}
			if result.Packet == nil || (*commandHost != "" && result.Lease == nil) {
				return 1
			}
			return 0
		}
		renderPacket(stdout, result)
		if result.Packet == nil || (*commandHost != "" && result.Lease == nil) {
			return 1
		}
		return 0
	}
	if command == "submit-result" {
		data, err := os.ReadFile(*resultFile)
		if err != nil {
			fmt.Fprintf(stderr, "pathframe: %v\n", err)
			return 1
		}
		var input app.SubmitResultInput
		if err := json.Unmarshal(data, &input.Result); err != nil {
			fmt.Fprintf(stderr, "pathframe: %v\n", err)
			return 1
		}
		result, err := service.SubmitResult(input)
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "verify" {
		files := []string{}
		if *changedFiles != "" {
			for _, name := range strings.Split(*changedFiles, ",") {
				if name = strings.TrimSpace(name); name != "" {
					files = append(files, name)
				}
			}
		}
		result, err := service.RunVerification(context.Background(), app.RunVerificationInput{Change: *commandChange, Task: *commandTask, ChangedFiles: files, Workdir: *workdir, TimeoutMS: *timeoutMS, MaxBytes: *maxBytes})
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "accept-task" {
		result, err := service.AcceptTask(app.AcceptTaskInput{Change: *commandChange, Task: *commandTask, Reason: *reason, KeepScopeViolations: *keepScope})
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "request-changes" {
		result, err := service.RequestChanges(app.RequestChangesInput{Change: *commandChange, Task: *commandTask, Reason: *reason})
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "lease-release" {
		result, err := service.ReleaseLease(app.ReleaseLeaseInput{Change: *commandChange, LeaseID: *leaseID})
		return renderArtifactResult(stdout, stderr, result, *commandJSON, true, err)
	}
	if command == "edit-check" {
		result, err := service.CheckBrainEdit(app.EditGuardInput{Change: *commandChange, Task: *commandTask})
		if *commandJSON {
			return renderArtifactResult(stdout, stderr, result, true, true, err)
		}
		fmt.Fprintf(stdout, "Brain edit allowed: %t\nReason: %s\n", result.Allowed, result.Reason)
		if err != nil || !result.Allowed {
			return 1
		}
		return 0
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
		case codex.InstallResult:
			fmt.Fprintf(w, "Codex integration %s installed\n", result.Version)
			for _, path := range result.Changed {
				fmt.Fprintf(w, "Changed: %s\n", path)
			}
		case claude.InstallResult:
			fmt.Fprintf(w, "Claude Code integration %s installed\n", result.Version)
			for _, path := range result.Changed {
				fmt.Fprintf(w, "Changed: %s\n", path)
			}
		case app.SubmitResultOutput:
			fmt.Fprintf(w, "Result %s for %s/%s reconciled to %s; phase: %s\n", result.LeaseID, result.Change, result.Task, result.Reconciliation.TaskState, result.Phase)
		case app.ReleaseLeaseOutput:
			fmt.Fprintf(w, "Lease %s released for %s/%s; phase: %s\n", result.LeaseID, result.Change, result.Task, result.Phase)
		case app.RunVerificationOutput:
			fmt.Fprintf(w, "Verification for %s/%s passed: %t; identity: %s\n", result.Change, result.Task, result.Verification.Passed, result.Verification.ContentIdentity)
		case app.ReviewOutput:
			fmt.Fprintf(w, "Task %s/%s %s; phase: %s; progress: %d/%d\n", result.Change, result.Task, result.Decision, result.Phase, result.Progress.Completed, result.Progress.Total)
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

func renderPacket(w io.Writer, result app.PrepareResult) {
	if result.Packet == nil {
		fmt.Fprintln(w, "Task packet unavailable")
		for _, issue := range result.Issues {
			fmt.Fprintf(w, "Issue: %s\n", issue)
		}
		for _, recovery := range result.Recovery {
			fmt.Fprintf(w, "Recovery: %s\n", recovery)
		}
		return
	}
	packet := result.Packet
	fmt.Fprintf(w, "Task packet %s\n\nChange: %s\nTask: %s - %s\nExecution policy: %s\nRole: %s\n", packet.Schema, packet.Change, packet.Task, packet.Title, packet.ExecutionPolicy, packet.Role.ID)
	fmt.Fprintf(w, "Role mission: %s\nObjective:\n%s\nAcceptance:\n%s\n", packet.Role.Mission, packet.Objective, packet.Acceptance)
	for _, field := range []struct {
		name  string
		value any
	}{
		{"References", packet.References}, {"Dependencies", packet.Dependencies}, {"Completed prerequisites", packet.CompletedPrerequisites},
		{"Write scope", packet.WriteScope}, {"Constraints", packet.Constraints}, {"Verification", packet.Verification},
		{"Role actions", packet.Role.AllowedActions}, {"Role returns", packet.Role.Returns},
	} {
		data, _ := json.Marshal(field.value)
		fmt.Fprintf(w, "%s: %s\n", field.name, data)
	}
	fmt.Fprintf(w, "Context budget: %d/%d bytes (%d required, %d omitted)\nWrite scope assurance: %s\n", packet.Budget.IncludedBytes, packet.Budget.LimitBytes, packet.Budget.RequiredBytes, packet.Budget.OmittedBytes, packet.WriteScopeAssurance)
	for _, entry := range packet.Context {
		fmt.Fprintf(w, "Context: %s %s (%d bytes, required=%t)\n%s\n", entry.Layer, entry.Path, entry.Bytes, entry.Required, entry.Content)
	}
	for _, omission := range packet.Omissions {
		fmt.Fprintf(w, "Omitted: %s %s (%d bytes): %s\n", omission.Layer, omission.Path, omission.Bytes, omission.Reason)
	}
	fmt.Fprintf(w, "Frontier: %s\n", strings.Join(result.Projection.Frontier, ", "))
	fmt.Fprintf(w, "Sequential next task: %s\n", result.Projection.Next)
	for i, wave := range result.Projection.Waves {
		fmt.Fprintf(w, "Wave %d: %s\n", i+1, strings.Join(wave, ", "))
	}
	if result.Preflight != nil {
		fmt.Fprintf(w, "Preflight ready: %t\nHost: %s\n", result.Preflight.Ready, result.Preflight.Host)
	}
	if result.Lease != nil {
		fmt.Fprintf(w, "Lease: %s (%s)\n", result.Lease.ID, result.Lease.State)
	}
	for _, issue := range result.Issues {
		fmt.Fprintf(w, "Issue: %s\n", issue)
	}
	for _, recovery := range result.Recovery {
		fmt.Fprintf(w, "Recovery: %s\n", recovery)
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: pathframe [--json] [--change ID] [status|next|new|template|check|approve|packet|submit-result|verify|accept-task|request-changes|lease-release|edit-check|pause|resume|replan|cancel|codex-install|codex-doctor|claude-install|claude-doctor]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "No command and status both show project orientation.")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  --change ID  select a change when more than one exists")
	fmt.Fprintln(w, "  --json       render pathframe.workflow/v1 JSON")
	fmt.Fprintln(w, "  --mode MODE  quick, standard, or high-risk (new/template)")
	fmt.Fprintln(w, "  --artifact K artifact kind (template)")
	fmt.Fprintln(w, "  --task ID     task to preview (packet; defaults to first frontier task)")
	fmt.Fprintln(w, "  --budget-bytes N explicit context byte budget (packet)")
	fmt.Fprintln(w, "  --host HOST   preflight and lease for codex or claude-code (packet)")
	fmt.Fprintln(w, "  --result-file PATH structured Pinky result (submit-result)")
	fmt.Fprintln(w, "  --timeout-ms N required verification timeout (verify)")
	fmt.Fprintln(w, "  --workdir PATH project-relative verification workdir (verify)")
	fmt.Fprintln(w, "  --max-bytes N per-stream output bound (verify)")
	fmt.Fprintln(w, "  --changed-files LIST comma-separated changed files for direct Brain verification")
	fmt.Fprintln(w, "  --reason TEXT semantic review reason (accept-task/request-changes)")
	fmt.Fprintln(w, "  --keep-scope-violations explicitly keep advisory out-of-scope files (accept-task)")
	fmt.Fprintln(w, "  --lease-id ID exact lost lease to release (lease-release)")
	fmt.Fprintln(w, "  --task ID     task whose Brain edit authority is checked (edit-check)")
	fmt.Fprintln(w, "  --session-orientation install the optional SessionStart hook (host install commands)")
	fmt.Fprintln(w, "  --help       show help")
	fmt.Fprintln(w, "  --version    show version")
}
