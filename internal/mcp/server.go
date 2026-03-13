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

	// ─── Work Items ─────────────────────────────────────────────────
	workItemActions := "'get', 'batch_get', 'list_types', 'get_links', 'get_history', 'list_comments'"
	if enableWrites {
		workItemActions += ", 'create', 'update', 'delete', 'add_comment'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_work_items",
		Description: "Manage Azure DevOps work items (tasks, bugs, user stories, epics). Actions: " + workItemActions,
	}, ManageWorkItemsHandler(c, enableWrites))

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

	// ─── Search ─────────────────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_search",
		Description: "Search Azure DevOps using WIQL, code search, work item search, or wiki search. Actions: 'wiql', 'code', 'work_items', 'wiki'",
	}, ManageSearchHandler(c))

	// ─── Repositories ───────────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_repos",
		Description: "Manage Azure DevOps Git repositories. Actions: 'list', 'get', 'list_branches', 'get_file', 'get_tree'",
	}, ManageReposHandler(c))

	// ─── Pull Requests ──────────────────────────────────────────────
	prActions := "'list', 'get', 'list_comments', 'list_reviewers'"
	if enableWrites {
		prActions += ", 'create', 'update', 'add_comment', 'vote'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_pull_requests",
		Description: "Manage Azure DevOps pull requests. Actions: " + prActions,
	}, ManagePullRequestsHandler(c, enableWrites))

	// ─── Iterations ─────────────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_iterations",
		Description: "Manage Azure DevOps iterations (sprints). Actions: 'list', 'get', 'get_current'",
	}, ManageIterationsHandler(c))

	// ─── Boards ─────────────────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_boards",
		Description: "Manage Azure DevOps Kanban boards. Actions: 'list', 'get', 'get_columns'",
	}, ManageBoardsHandler(c))

	// ─── Wiki ───────────────────────────────────────────────────────
	wikiActions := "'list', 'get_page'"
	if enableWrites {
		wikiActions += ", 'create_page', 'update_page', 'delete_page'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_wiki",
		Description: "Manage Azure DevOps wiki pages (markdown-native). Actions: " + wikiActions,
	}, ManageWikiHandler(c, enableWrites))

	// ─── Pipelines ──────────────────────────────────────────────────
	pipelineActions := "'list', 'get', 'list_runs', 'get_run', 'get_logs', 'get_log'"
	if enableWrites {
		pipelineActions += ", 'trigger'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_pipelines",
		Description: "Manage Azure DevOps CI/CD pipelines. Actions: " + pipelineActions,
	}, ManagePipelinesHandler(c, enableWrites))

	// ─── Attachments ────────────────────────────────────────────────
	attachmentActions := "'list', 'download'"
	if enableWrites {
		attachmentActions += ", 'upload'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_attachments",
		Description: "Manage Azure DevOps work item attachments. Actions: " + attachmentActions,
	}, ManageAttachmentsHandler(c, enableWrites))
}
