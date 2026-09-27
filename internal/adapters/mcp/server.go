package mcpadapter

import (
	"context"

	"github.com/0xkhdr/pathframe/internal/app"
	"github.com/0xkhdr/pathframe/internal/integrations/codex"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func New(service app.Service) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "pathframe", Title: "Pathframe", Version: codex.IntegrationVersion}, &mcp.ServerOptions{
		Instructions: "Use typed Pathframe tools for planning and sequential delegation. Ask for explicit human approval before approving a valid plan. Never construct Pathframe CLI commands or treat delegation failure as permission for Brain to edit delegated work.",
	})
	addTools(server, service)
	return server
}

func Run(ctx context.Context, service app.Service) error {
	return New(service).Run(ctx, &mcp.StdioTransport{})
}
