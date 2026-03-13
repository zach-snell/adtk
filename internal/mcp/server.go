package mcp

import (
	"context"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
	"github.com/zach-snell/adtk/internal/version"
)

// New creates and configures the Azure DevOps MCP server.
func New(organization, pat string) *mcp.Server {
	client := devops.NewClient(organization, pat)
	return newServer(client)
}

func newServer(client *devops.Client) *mcp.Server {
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "adtk",
			Version: version.Version,
		},
		nil,
	)

	registerTools(s, client)
	return s
}

// addTool is a helper function to conditionally register a tool handler.
func addTool[In any](s *mcp.Server, disabled map[string]bool, tool mcp.Tool, handler func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, any, error)) {
	if disabled[tool.Name] {
		return
	}
	mcp.AddTool(s, &tool, handler)
}

func registerTools(s *mcp.Server, c *devops.Client) {
	disabledToolsEnv := os.Getenv("AZURE_DEVOPS_DISABLED_TOOLS")
	disabled := make(map[string]bool)
	if disabledToolsEnv != "" {
		for _, t := range strings.Split(disabledToolsEnv, ",") {
			disabled[strings.TrimSpace(t)] = true
		}
	}

	enableWrites := os.Getenv("ADTK_ENABLE_WRITES") == "true"

	// ─── Projects ────────────────────────────────────────────────────
	projectActions := "'list', 'get', 'list_teams', 'get_team'"
	if enableWrites {
		projectActions += ", 'create'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_projects",
		Description: "Manage Azure DevOps projects and teams. Actions: " + projectActions,
	}, ManageProjectsHandler(c, enableWrites))

	// ─── Users ──────────────────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_users",
		Description: "Search and get Azure DevOps users. Actions: 'get_current', 'search'",
	}, ManageUsersHandler(c))
}
