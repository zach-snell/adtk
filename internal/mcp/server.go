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
	workItemActions := "'get', 'batch_get', 'list_types', 'get_links', 'get_history', 'list_comments', 'my_items', 'iteration_items'"
	if enableWrites {
		workItemActions += ", 'create', 'update', 'delete', 'add_comment', 'batch_update', 'add_children', 'link', 'unlink', 'add_artifact_link', 'update_comment'"
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
	repoActions := "'list', 'get', 'list_branches', 'get_file', 'get_tree', 'search_commits'"
	if enableWrites {
		repoActions += ", 'create_branch'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_repos",
		Description: "Manage Azure DevOps Git repositories. Actions: " + repoActions,
	}, ManageReposHandler(c, enableWrites))

	// ─── Pull Requests ──────────────────────────────────────────────
	prActions := "'list', 'get', 'list_comments', 'list_reviewers'"
	if enableWrites {
		prActions += ", 'create', 'update', 'add_comment', 'vote', 'update_reviewers', 'create_thread', 'update_thread', 'reply_to_comment'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_pull_requests",
		Description: "Manage Azure DevOps pull requests. Actions: " + prActions,
	}, ManagePullRequestsHandler(c, enableWrites))

	// ─── Iterations ─────────────────────────────────────────────────
	iterationActions := "'list', 'get', 'get_current', 'get_team_settings'"
	if enableWrites {
		iterationActions += ", 'create'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_iterations",
		Description: "Manage Azure DevOps iterations (sprints). Actions: " + iterationActions,
	}, ManageIterationsHandler(c, enableWrites))

	// ─── Boards ─────────────────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_boards",
		Description: "Manage Azure DevOps Kanban boards. Actions: 'list', 'get', 'get_columns'",
	}, ManageBoardsHandler(c))

	// ─── Wiki ───────────────────────────────────────────────────────
	wikiActions := "'list', 'get_page', 'list_pages'"
	if enableWrites {
		wikiActions += ", 'create_page', 'update_page', 'delete_page'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_wiki",
		Description: "Manage Azure DevOps wiki pages (markdown-native). Actions: " + wikiActions,
	}, ManageWikiHandler(c, enableWrites))

	// ─── Pipelines ──────────────────────────────────────────────────
	pipelineActions := "'list', 'get', 'list_runs', 'get_run', 'get_logs', 'get_log', 'get_build_changes', 'list_definitions'"
	if enableWrites {
		pipelineActions += ", 'trigger'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_pipelines",
		Description: "Manage Azure DevOps CI/CD pipelines. Actions: " + pipelineActions,
	}, ManagePipelinesHandler(c, enableWrites))

	// ─── Test Plans ─────────────────────────────────────────────────
	testPlanActions := "'list_plans', 'list_suites', 'list_cases', 'get_test_results'"
	if enableWrites {
		testPlanActions += ", 'create_plan', 'create_suite'"
	}
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_test_plans",
		Description: "Manage Azure DevOps test plans, suites, and cases. Actions: " + testPlanActions,
	}, ManageTestPlansHandler(c, enableWrites))

	// ─── Advanced Security ──────────────────────────────────────────
	addTool(s, disabled, mcp.Tool{
		Name:        "manage_advanced_security",
		Description: "Manage Azure DevOps Advanced Security alerts. Actions: 'list_alerts', 'get_alert'",
	}, ManageAdvancedSecurityHandler(c))

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
