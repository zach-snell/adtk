package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerPrompts registers all MCP prompts on the server.
func registerPrompts(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Name:        "sprint_summary",
		Description: "Generate a sprint/iteration summary for an Azure DevOps project",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Azure DevOps project name", Required: true},
			{Name: "team", Description: "Team name (optional, scopes to a specific team)"},
			{Name: "iteration", Description: "Iteration path or name (defaults to current iteration)"},
		},
	}, sprintSummaryPromptHandler)

	s.AddPrompt(&mcp.Prompt{
		Name:        "pr_review_digest",
		Description: "Generate a PR review digest for an Azure DevOps repository",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Azure DevOps project name", Required: true},
			{Name: "repo", Description: "Repository name", Required: true},
		},
	}, prReviewDigestPromptHandler)

	s.AddPrompt(&mcp.Prompt{
		Name:        "pipeline_health",
		Description: "Analyze CI/CD pipeline health for an Azure DevOps project",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Azure DevOps project name", Required: true},
			{Name: "pipeline_id", Description: "Specific pipeline ID to analyze (omit for all pipelines)"},
		},
	}, pipelineHealthPromptHandler)

	s.AddPrompt(&mcp.Prompt{
		Name:        "release_readiness",
		Description: "Assess release readiness for an Azure DevOps project",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Azure DevOps project name", Required: true},
			{Name: "iteration", Description: "Iteration path or name to assess (defaults to current iteration)"},
		},
	}, releaseReadinessPromptHandler)
}

func sprintSummaryPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	if project == "" {
		return nil, fmt.Errorf("project is required")
	}
	team := req.Params.Arguments["team"]
	iteration := req.Params.Arguments["iteration"]

	teamClause := ""
	if team != "" {
		teamClause = fmt.Sprintf(" and team '%s'", team)
	}

	iterationStep := fmt.Sprintf(
		"1. Use manage_iterations with action 'get_current' and project_key '%s'", project)
	if iteration != "" {
		iterationStep = fmt.Sprintf(
			"1. Use manage_iterations with action 'get' for iteration '%s' in project '%s'", iteration, project)
	}
	if team != "" {
		iterationStep += fmt.Sprintf(" and team '%s'", team)
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Sprint summary for %s%s", project, teamClause),
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf(`Generate a sprint/iteration summary for project '%s'%s.

Steps:
%s
2. Use manage_work_items with action 'iteration_items' to get all work items in the iteration
3. Use manage_boards with action 'get_columns' to understand the board workflow
4. Categorize work items by state and type

Present as a sprint status report:
- Sprint/iteration name, dates, and timeline progress
- Completion percentage (done items / total items)
- Work items grouped by state (New, Active, Resolved, Closed)
- Breakdown by work item type (User Story, Bug, Task)
- Items at risk (active but not updated recently)
- Blocked items or items with impediments
- Team velocity summary if available`, project, teamClause, iterationStep),
				},
			},
		},
	}, nil
}

func prReviewDigestPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	repo := req.Params.Arguments["repo"]
	if project == "" || repo == "" {
		return nil, fmt.Errorf("project and repo are required")
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("PR review digest for %s/%s", project, repo),
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf(`Generate a pull request review digest for repository '%s' in project '%s'.

Steps:
1. Use manage_pull_requests with action 'list', project_key '%s', repo_id '%s', and status 'active' to get all open PRs
2. For each PR, use manage_pull_requests with action 'list_reviewers' to check reviewer status
3. Analyze each PR for review status, age, and potential issues

Present as a PR review digest:
- Total open PRs and summary statistics
- PRs needing review (no reviewers assigned or pending votes)
- PRs with requested changes (rejected votes)
- PRs ready to merge (all approved, no conflicts)
- Stale PRs (open for more than 7 days without activity)
- PRs by author showing review workload distribution
- For each PR: title, author, age, reviewer status, and merge conflicts if any`, repo, project, project, repo),
				},
			},
		},
	}, nil
}

func pipelineHealthPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	if project == "" {
		return nil, fmt.Errorf("project is required")
	}
	pipelineID := req.Params.Arguments["pipeline_id"]

	pipelineStep := fmt.Sprintf(
		"1. Use manage_pipelines with action 'list' and project_key '%s' to get all pipelines", project)
	if pipelineID != "" {
		pipelineStep = fmt.Sprintf(
			"1. Use manage_pipelines with action 'get' with project_key '%s' and pipeline_id %s to get the specific pipeline", project, pipelineID)
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Pipeline health for %s", project),
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf(`Analyze CI/CD pipeline health for project '%s'.

Steps:
%s
2. Use manage_pipelines with action 'list_runs' to get recent pipeline runs (use top 20)
3. For any failed runs, use manage_pipelines with action 'get_run' to get failure details
4. Optionally use manage_pipelines with action 'get_logs' to inspect build logs for failures

Present as a pipeline health report:
- Overall health score (pass rate over recent runs)
- Pipeline success/failure rate breakdown
- Average build duration and trends
- Most common failure reasons
- Pipelines with consecutive failures (needs immediate attention)
- Longest running pipelines
- Recently triggered vs idle pipelines
- Recommendations for improving pipeline reliability`, project, pipelineStep),
				},
			},
		},
	}, nil
}

func releaseReadinessPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	if project == "" {
		return nil, fmt.Errorf("project is required")
	}
	iteration := req.Params.Arguments["iteration"]

	iterationStep := fmt.Sprintf(
		"1. Use manage_iterations with action 'get_current' and project_key '%s' to find the current iteration", project)
	if iteration != "" {
		iterationStep = fmt.Sprintf(
			"1. Use manage_iterations with action 'get' for iteration '%s' in project '%s'", iteration, project)
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Release readiness for %s", project),
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf(`Assess release readiness for project '%s'.

Steps:
%s
2. Use manage_work_items with action 'iteration_items' to get all work items in the iteration
3. Use manage_search with action 'wiql' to find any open bugs:
   WIQL: SELECT [System.Id], [System.Title], [System.State] FROM WorkItems WHERE [System.TeamProject] = '%s' AND [System.WorkItemType] = 'Bug' AND [System.State] <> 'Closed' AND [System.State] <> 'Removed'
4. Use manage_pull_requests with action 'list' and status 'active' to find unmerged PRs
5. Use manage_pipelines with action 'list_runs' to check recent build status

Present as a release readiness assessment:
- Overall readiness score (Ready / At Risk / Not Ready)
- Work item completion: completed vs remaining items
- Open bugs: count and severity breakdown
- Unmerged PRs: list with status and blockers
- Build status: latest pipeline results (green/red)
- Test results: if available from recent builds
- Blockers and risks that could delay release
- Go/No-Go recommendation with justification`, project, iterationStep, project),
				},
			},
		},
	}, nil
}
