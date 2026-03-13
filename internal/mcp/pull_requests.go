package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManagePullRequestsInput defines the input schema for the manage_pull_requests tool.
type ManagePullRequestsInput struct {
	Action     string `json:"action" jsonschema:"required,description=Action to perform: 'list'\\, 'get'\\, 'create'\\, 'update'\\, 'add_comment'\\, 'list_comments'\\, 'vote'\\, 'list_reviewers'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"description=Project name"`
	RepoID     string `json:"repo_id,omitempty" jsonschema:"description=Repository name or ID (required for most actions)"`
	PRID       int    `json:"pr_id,omitempty" jsonschema:"description=Pull request ID (required for get\\, update\\, add_comment\\, list_comments\\, vote\\, list_reviewers)"`

	// List filters
	Status string `json:"status,omitempty" jsonschema:"description=Filter by status: active\\, completed\\, abandoned\\, all (for list)"`
	Top    int    `json:"top,omitempty" jsonschema:"description=Max results to return"`

	// Create fields
	Title        string `json:"title,omitempty" jsonschema:"description=PR title (required for create)"`
	Description  string `json:"description,omitempty" jsonschema:"description=PR description (for create\\, update)"`
	SourceBranch string `json:"source_branch,omitempty" jsonschema:"description=Source branch name (required for create)"`
	TargetBranch string `json:"target_branch,omitempty" jsonschema:"description=Target branch name (required for create)"`
	IsDraft      bool   `json:"is_draft,omitempty" jsonschema:"description=Create as draft PR (for create)"`

	// Comment
	Comment string `json:"comment,omitempty" jsonschema:"description=Comment content (for add_comment)"`

	// Vote
	ReviewerID string `json:"reviewer_id,omitempty" jsonschema:"description=Reviewer ID (for vote)"`
	Vote       int    `json:"vote,omitempty" jsonschema:"description=Vote: 10=approved\\, 5=approved with suggestions\\, 0=no vote\\, -5=waiting\\, -10=rejected"`
}

// ManagePullRequestsHandler returns the handler for the manage_pull_requests tool.
func ManagePullRequestsHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "list":
			if input.RepoID == "" {
				return resultError("repo_id is required for 'list' action")
			}
			prs, err := c.ListPullRequests(input.ProjectKey, input.RepoID, input.Status, input.Top)
			if err != nil {
				return resultError(fmt.Sprintf("listing PRs: %v", err))
			}
			return resultJSON(prs)
		case "get":
			if input.RepoID == "" || input.PRID == 0 {
				return resultError("repo_id and pr_id are required for 'get' action")
			}
			pr, err := c.GetPullRequest(input.ProjectKey, input.RepoID, input.PRID)
			if err != nil {
				return resultError(fmt.Sprintf("getting PR: %v", err))
			}
			return resultJSON(pr)
		case "create":
			if !enableWrites {
				return resultError("create action requires ADTK_ENABLE_WRITES=true")
			}
			if input.RepoID == "" || input.Title == "" || input.SourceBranch == "" || input.TargetBranch == "" {
				return resultError("repo_id, title, source_branch, and target_branch are required for 'create' action")
			}
			pr := &devops.PullRequest{
				Title:         input.Title,
				Description:   input.Description,
				SourceRefName: "refs/heads/" + input.SourceBranch,
				TargetRefName: "refs/heads/" + input.TargetBranch,
				IsDraft:       input.IsDraft,
			}
			result, err := c.CreatePullRequest(input.ProjectKey, input.RepoID, pr)
			if err != nil {
				return resultError(fmt.Sprintf("creating PR: %v", err))
			}
			return resultJSON(result)
		case "update":
			if !enableWrites {
				return resultError("update action requires ADTK_ENABLE_WRITES=true")
			}
			if input.RepoID == "" || input.PRID == 0 {
				return resultError("repo_id and pr_id are required for 'update' action")
			}
			update := make(map[string]interface{})
			if input.Title != "" {
				update["title"] = input.Title
			}
			if input.Description != "" {
				update["description"] = input.Description
			}
			if input.Status != "" {
				update["status"] = input.Status
			}
			result, err := c.UpdatePullRequest(input.ProjectKey, input.RepoID, input.PRID, update)
			if err != nil {
				return resultError(fmt.Sprintf("updating PR: %v", err))
			}
			return resultJSON(result)
		case "add_comment":
			if !enableWrites {
				return resultError("add_comment action requires ADTK_ENABLE_WRITES=true")
			}
			if input.RepoID == "" || input.PRID == 0 || input.Comment == "" {
				return resultError("repo_id, pr_id, and comment are required for 'add_comment' action")
			}
			thread, err := c.AddPRComment(input.ProjectKey, input.RepoID, input.PRID, input.Comment)
			if err != nil {
				return resultError(fmt.Sprintf("adding PR comment: %v", err))
			}
			return resultJSON(thread)
		case "list_comments":
			if input.RepoID == "" || input.PRID == 0 {
				return resultError("repo_id and pr_id are required for 'list_comments' action")
			}
			threads, err := c.ListPRThreads(input.ProjectKey, input.RepoID, input.PRID)
			if err != nil {
				return resultError(fmt.Sprintf("listing PR comments: %v", err))
			}
			return resultJSON(threads)
		case "vote":
			if !enableWrites {
				return resultError("vote action requires ADTK_ENABLE_WRITES=true")
			}
			if input.RepoID == "" || input.PRID == 0 || input.ReviewerID == "" {
				return resultError("repo_id, pr_id, and reviewer_id are required for 'vote' action")
			}
			if err := c.VotePR(input.ProjectKey, input.RepoID, input.PRID, input.ReviewerID, input.Vote); err != nil {
				return resultError(fmt.Sprintf("voting on PR: %v", err))
			}
			return resultText(fmt.Sprintf("Vote %d submitted for PR %d", input.Vote, input.PRID))
		case "list_reviewers":
			if input.RepoID == "" || input.PRID == 0 {
				return resultError("repo_id and pr_id are required for 'list_reviewers' action")
			}
			reviewers, err := c.ListPRReviewers(input.ProjectKey, input.RepoID, input.PRID)
			if err != nil {
				return resultError(fmt.Sprintf("listing reviewers: %v", err))
			}
			return resultJSON(reviewers)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
