package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	pfcontext "github.com/0xkhdr/pathframe/internal/context"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/store"
	"github.com/0xkhdr/pathframe/internal/workflow"
)

type PrepareInput struct {
	Change      string `json:"change,omitempty"`
	Task        string `json:"task,omitempty"`
	BudgetBytes int    `json:"budget_bytes,omitempty"`
	Host        string `json:"host,omitempty"`
}

type PrepareResult struct {
	Packet             *pfcontext.Packet     `json:"packet,omitempty"`
	Projection         delegation.Projection `json:"projection"`
	Preflight          *delegation.Preflight `json:"preflight,omitempty"`
	Lease              *delegation.Lease     `json:"lease,omitempty"`
	WorkerInstructions string                `json:"worker_instructions,omitempty"`
	Issues             []string              `json:"issues"`
	Recovery           []string              `json:"recovery"`
}

func (s Service) PrepareDelegation(input PrepareInput) (PrepareResult, error) {
	view, err := s.Orient(input.Change)
	if err != nil {
		return PrepareResult{}, err
	}
	result := PrepareResult{Issues: []string{}, Recovery: []string{}}
	if view.Change == "" {
		return result, fmt.Errorf("cannot prepare a packet without a selected change")
	}
	if view.Phase != workflow.PhaseReady && view.Phase != workflow.PhaseBlocked {
		result.Issues = []string{"change must be ready or blocked for delegation preparation"}
		result.Recovery = []string{"validate and approve the plan, or recover the active execution blocker"}
		return result, nil
	}
	root, err := store.Discover(s.Dir)
	if err != nil {
		return PrepareResult{}, err
	}
	dir, err := store.ChangeDir(root, view.Change)
	if err != nil {
		return PrepareResult{}, err
	}
	plan, issues := artifacts.Validate(dir)
	if len(issues) > 0 {
		result.Issues = []string{"approved plan is no longer valid"}
		result.Recovery = []string{"repair the named artifacts, revalidate, and request human reapproval"}
		return result, nil
	}
	result.Projection = delegation.Project(plan.Tasks, nil)
	if len(result.Projection.Diagnostics) > 0 {
		result.Issues = append(result.Issues, result.Projection.Diagnostics...)
		result.Recovery = append(result.Recovery, result.Projection.Recovery...)
		return result, nil
	}
	taskID := input.Task
	if taskID == "" {
		taskID = result.Projection.Next
	}
	var task *artifacts.Task
	for i := range plan.Tasks {
		if plan.Tasks[i].ID == taskID {
			task = &plan.Tasks[i]
			break
		}
	}
	if task == nil {
		result.Issues = []string{"task is not present in the approved plan: " + taskID}
		result.Recovery = []string{"select a task identifier from the projected frontier or waves"}
		return result, nil
	}
	ready := false
	for _, id := range result.Projection.Frontier {
		ready = ready || id == task.ID
	}
	if !ready {
		result.Issues = []string{"task is not in the current ready frontier: " + task.ID}
		result.Recovery = []string{"complete its projected prerequisites or select a frontier task"}
		return result, nil
	}
	role := pfcontext.Role{Schema: pfcontext.RoleSchema, ID: "none", Mission: "Execute the approved task contract.", Reads: []string{}, OptionalReads: []string{}, AllowedActions: []string{"edit", "test"}, Returns: []string{"summary", "changed_files", "verification", "discoveries", "questions"}}
	if task.Role != "" && task.Role != "none" {
		if filepath.Base(task.Role) != task.Role || strings.ContainsAny(task.Role, `/\\`) || task.Role == "." || task.Role == ".." {
			result.Issues = []string{"invalid role identifier: " + task.Role}
			result.Recovery = []string{"assign a simple role identifier and request human reapproval"}
			return result, nil
		}
		role, err = pfcontext.LoadRole(filepath.Join(root, ".pathframe", "roles", task.Role+".yaml"))
		if err != nil {
			result.Issues = []string{"role cannot be loaded: " + err.Error()}
			result.Recovery = []string{"create or repair .pathframe/roles/" + task.Role + ".yaml and preview again"}
			return result, nil
		}
		if role.ID != task.Role {
			result.Issues = []string{"role id does not match task role: " + role.ID}
			result.Recovery = []string{"make the role id match " + task.Role}
			return result, nil
		}
	} else if task.ExecutionPolicy == "delegated" {
		result.Issues = []string{"delegated task requires an explicit role"}
		result.Recovery = []string{"assign an existing role and request human reapproval"}
		return result, nil
	}
	references := make([]pfcontext.Reference, 0, len(role.Reads)+len(role.OptionalReads)+len(task.Sections["Required Reads"]))
	for _, path := range role.Reads {
		references = append(references, pfcontext.Reference{Layer: pfcontext.Foundation, Path: path, Required: true})
	}
	for _, path := range role.OptionalReads {
		references = append(references, pfcontext.Reference{Layer: pfcontext.Foundation, Path: path})
	}
	for _, path := range sectionBullets(task.Sections["Required Reads"]) {
		references = append(references, pfcontext.Reference{Layer: pfcontext.Change, Path: path, Required: true})
	}
	resolved, unresolved, err := pfcontext.Resolve(root, dir, references)
	if err != nil {
		result.Issues = []string{err.Error()}
		result.Recovery = []string{"repair the context reference or amend and reapprove the task"}
		return result, nil
	}
	limit := input.BudgetBytes
	if limit == 0 {
		limit = pfcontext.DefaultBudget
	}
	runtime := fmt.Sprintf("change=%s\nphase=%s\nfrontier=%s", view.Change, view.Phase, strings.Join(result.Projection.Frontier, ","))
	packet, err := pfcontext.Assemble(view.Change, *task, role, resolved, unresolved, runtime, limit)
	if err != nil {
		result.Issues = []string{err.Error()}
		result.Recovery = []string{"split the task or explicitly raise budget_bytes and preview again"}
		return result, nil
	}
	result.Packet = &packet
	if input.Host == "" {
		return result, nil
	}
	hostDir, manifestHost := input.Host, input.Host
	switch input.Host {
	case "codex":
	case "claude-code":
		hostDir = "claude"
	default:
		result.Issues = []string{"unsupported delegation host: " + input.Host}
		result.Recovery = []string{"select codex or claude-code and install its integration"}
		return result, nil
	}
	capabilities, err := delegation.LoadCapabilities(filepath.Join(root, ".pathframe", "integrations", hostDir, "manifest.json"), manifestHost)
	if err != nil {
		result.Issues = []string{"host capabilities cannot be loaded: " + err.Error()}
		result.Recovery = []string{"install or repair the " + input.Host + " integration, then retry"}
		return result, nil
	}
	runsDir := filepath.Join(dir, "runs")
	_, active, err := delegation.ActiveLease(runsDir)
	if err != nil {
		return PrepareResult{}, err
	}
	preflight := delegation.Check(packet, capabilities, active)
	result.Preflight = &preflight
	if !preflight.Ready {
		result.Issues, result.Recovery = preflight.Issues, preflight.Recovery
		return result, nil
	}
	packet.WriteScopeAssurance, packet.HostAssurance = preflight.WriteScopeAssurance, "declared"
	result.Packet = &packet
	lease, err := delegation.Acquire(runsDir, view.Change, task.ID, input.Host, now())
	if err != nil {
		return PrepareResult{}, err
	}
	if err := s.appendDelegationTransition(dir, workflow.ActionExecute, workflow.ActorBrain, "delegation lease acquired for "+lease.ID); err != nil {
		_ = delegation.Release(runsDir, lease)
		return PrepareResult{}, err
	}
	result.Lease = &lease
	result.WorkerInstructions = delegation.WorkerInstructions(packet)
	return result, nil
}

func sectionBullets(body string) []string {
	values := []string{}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if value != "" && value != "none" {
				values = append(values, value)
			}
		}
	}
	return values
}
