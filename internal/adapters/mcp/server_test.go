package mcpadapter

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/0xkhdr/pathframe/internal/app"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTypedPlanningToolsAndCanonicalOrientation(t *testing.T) {
	ctx := context.Background()
	service := app.Service{Dir: t.TempDir()}
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := New(service).Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()

	listed, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"pathframe_orient": true, "pathframe_assess_request": true, "pathframe_create_change": true, "pathframe_get_template": true, "pathframe_validate_plan": true, "pathframe_get_next": true, "pathframe_prepare_delegation": true, "pathframe_recover": true}
	for _, tool := range listed.Tools {
		delete(want, tool.Name)
		if tool.Name == "run_pathframe_command" {
			t.Fatal("generic raw-command tool exposed")
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing tools: %v", want)
	}

	call, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: "pathframe_orient", Arguments: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(call.StructuredContent)
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got["schema"] != "pathframe.workflow/v1" {
		t.Fatalf("orientation schema = %v", got["schema"])
	}
	canonical, err := service.Orient("")
	if err != nil {
		t.Fatal(err)
	}
	wantData, _ := json.Marshal(canonical)
	var wantOrientation map[string]any
	_ = json.Unmarshal(wantData, &wantOrientation)
	if !reflect.DeepEqual(got, wantOrientation) {
		t.Fatalf("MCP orientation = %s, app orientation = %s (%v)", data, wantData, wantOrientation)
	}
	assessment, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: "pathframe_assess_request", Arguments: map[string]any{"dependent_steps": true}})
	if err != nil {
		t.Fatal(err)
	}
	data, _ = json.Marshal(assessment.StructuredContent)
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got["activation"] != "offer" {
		t.Fatalf("activation = %v", got["activation"])
	}
}
