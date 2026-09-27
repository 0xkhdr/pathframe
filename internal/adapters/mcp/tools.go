package mcpadapter

import (
	"context"
	"fmt"

	"github.com/0xkhdr/pathframe/internal/app"
	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/workflow"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type changeInput struct {
	Change string `json:"change,omitempty" jsonschema:"change identifier; omit only when exactly one change is active"`
}

type createInput struct {
	Change string         `json:"change" jsonschema:"stable change identifier"`
	Mode   artifacts.Mode `json:"mode" jsonschema:"planning mode: quick, standard, or high-risk"`
}

type templateInput struct {
	Mode     artifacts.Mode `json:"mode" jsonschema:"planning mode: quick, standard, or high-risk"`
	Artifact string         `json:"artifact" jsonschema:"artifact kind"`
}

type templateOutput struct {
	Instructions artifacts.Instructions `json:"instructions"`
	Template     string                 `json:"template"`
}

type validateInput struct {
	Change        string `json:"change,omitempty" jsonschema:"change identifier"`
	HumanApproved bool   `json:"human_approved" jsonschema:"true only after the human explicitly approves the currently validated plan"`
}

type validateOutput struct {
	Check            artifacts.CheckResult     `json:"check"`
	ApprovalRequired bool                      `json:"approval_required"`
	Approval         *artifacts.ApprovalResult `json:"approval,omitempty"`
}

type recoverInput struct {
	Change string          `json:"change,omitempty" jsonschema:"change identifier"`
	Action workflow.Action `json:"action" jsonschema:"one of pause, resume, replan, or cancel"`
}

type prepareInput struct {
	Change      string `json:"change,omitempty" jsonschema:"change identifier"`
	Task        string `json:"task,omitempty" jsonschema:"task identifier; omit to select the first deterministic frontier task"`
	BudgetBytes int    `json:"budget_bytes,omitempty" jsonschema:"positive context byte budget; omit for the documented default"`
}

func addTools(server *mcp.Server, service app.Service) {
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_orient", Description: "Read canonical project and active-change orientation. Safe at session start."}, func(ctx context.Context, _ *mcp.CallToolRequest, in changeInput) (*mcp.CallToolResult, workflow.Result, error) {
		out, err := service.Orient(in.Change)
		return nil, out, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_assess_request", Description: "Apply Pathframe must-use, offer-once, and must-not-use rules to request facts classified by Brain."}, func(ctx context.Context, _ *mcp.CallToolRequest, in app.RequestAssessmentInput) (*mcp.CallToolResult, app.RequestAssessment, error) {
		return nil, service.AssessRequest(in), nil
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_create_change", Description: "Create the selected planning mode and return its initial validation result."}, func(ctx context.Context, _ *mcp.CallToolRequest, in createInput) (*mcp.CallToolResult, artifacts.CheckResult, error) {
		out, err := service.CreateChange(in.Change, in.Mode)
		return nil, out, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_get_template", Description: "Return typed instructions and the human-readable template for one planning artifact."}, func(ctx context.Context, _ *mcp.CallToolRequest, in templateInput) (*mcp.CallToolResult, templateOutput, error) {
		instructions, body, err := service.Template(in.Mode, in.Artifact)
		return nil, templateOutput{Instructions: instructions, Template: body}, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_validate_plan", Description: "Validate a plan. Set human_approved only after explicit human approval; then this same typed operation records approval and advances to ready."}, func(ctx context.Context, _ *mcp.CallToolRequest, in validateInput) (*mcp.CallToolResult, validateOutput, error) {
		check, err := service.Check(in.Change)
		out := validateOutput{Check: check, ApprovalRequired: check.Valid && !in.HumanApproved}
		if err != nil || !in.HumanApproved || !check.Valid {
			return nil, out, err
		}
		approval, err := service.Approve(in.Change)
		out.Approval, out.ApprovalRequired = &approval, false
		return nil, out, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_get_next", Description: "Return the canonical recommended next action and legal alternatives."}, func(ctx context.Context, _ *mcp.CallToolRequest, in changeInput) (*mcp.CallToolResult, workflow.Result, error) {
		out, err := service.GetNext(in.Change)
		return nil, out, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_prepare_delegation", Description: "Preview a bounded pathframe.task/v1 packet and deterministic sequential frontier. This does not launch a worker."}, func(ctx context.Context, _ *mcp.CallToolRequest, in prepareInput) (*mcp.CallToolResult, app.PrepareResult, error) {
		out, err := service.PrepareDelegation(app.PrepareInput{Change: in.Change, Task: in.Task, BudgetBytes: in.BudgetBytes})
		return nil, out, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "pathframe_recover", Description: "Apply one supported recovery action: pause, resume, replan, or cancel."}, func(ctx context.Context, _ *mcp.CallToolRequest, in recoverInput) (*mcp.CallToolResult, workflow.Result, error) {
		var out workflow.Result
		var err error
		switch in.Action {
		case workflow.ActionPause:
			out, err = service.Pause(in.Change)
		case workflow.ActionResume:
			out, err = service.Resume(in.Change)
		case workflow.ActionReplan:
			out, err = service.Replan(in.Change)
		case workflow.ActionCancel:
			out, err = service.Cancel(in.Change)
		default:
			err = fmt.Errorf("unsupported recovery action %q; use pause, resume, replan, or cancel", in.Action)
		}
		return nil, out, err
	})
}
