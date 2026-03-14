package mcp

import (
	"context"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManagePullRequestsInput defines the input schema for the manage_pull_requests tool.
type ManagePullRequestsInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'list', 'get', 'create', 'update', 'add_comment', 'list_comments', 'vote', 'list_reviewers', 'update_reviewers', 'create_thread', 'update_thread', 'reply_to_comment'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name"`
	RepoID     string `json:"repo_id,omitempty" jsonschema:"Repository name or ID (required for most actions)"`
	PRID       int    `json:"pr_id,omitempty" jsonschema:"Pull request ID (required for get, update, add_comment, list_comments, vote, list_reviewers)"`

	// List filters
	Status string `json:"status,omitempty" jsonschema:"Filter by status: active, completed, abandoned, all (for list)"`
	Top    int    `json:"top,omitempty" jsonschema:"Max results to return"`

	// Create fields
	Title        string `json:"title,omitempty" jsonschema:"PR title (required for create)"`
	Description  string `json:"description,omitempty" jsonschema:"PR description (for create, update)"`
	SourceBranch string `json:"source_branch,omitempty" jsonschema:"Source branch name (required for create)"`
	TargetBranch string `json:"target_branch,omitempty" jsonschema:"Target branch name (required for create)"`
	IsDraft      bool   `json:"is_draft,omitempty" jsonschema:"Create as draft PR (for create)"`

	// Comment
	Comment string `json:"comment,omitempty" jsonschema:"Comment content (for add_comment, create_thread, reply_to_comment)"`

	// Vote
	ReviewerID string `json:"reviewer_id,omitempty" jsonschema:"Reviewer ID (for vote)"`
	Vote       int    `json:"vote,omitempty" jsonschema:"Vote: 10=approved, 5=approved with suggestions, 0=no vote, -5=waiting, -10=rejected"`

	// Reviewer management
	ReviewerIDs string `json:"reviewer_ids,omitempty" jsonschema:"Comma-separated reviewer IDs (for update_reviewers)"`

	// Thread management
	ThreadID int    `json:"thread_id,omitempty" jsonschema:"Thread ID (for update_thread, reply_to_comment)"`
	FilePath string `json:"file_path,omitempty" jsonschema:"File path for inline comment (for create_thread)"`
	Line     int    `json:"line,omitempty" jsonschema:"Line number for inline comment (for create_thread)"`
}

// ManagePullRequestsHandler returns the handler for the manage_pull_requests tool.
func ManagePullRequestsHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "list":
			return handlePRList(c, input)
		case "get":
			return handlePRGet(c, input)
		case "create":
			return handlePRCreate(c, input, enableWrites)
		case "update":
			return handlePRUpdate(c, input, enableWrites)
		case "add_comment":
			return handlePRAddComment(c, input, enableWrites)
		case "list_comments":
			return handlePRListComments(c, input)
		case "vote":
			return handlePRVote(c, input, enableWrites)
		case "list_reviewers":
			return handlePRListReviewers(c, input)
		case "update_reviewers":
			return handlePRUpdateReviewers(c, input, enableWrites)
		case "create_thread":
			return handlePRCreateThread(c, input, enableWrites)
		case "update_thread":
			return handlePRUpdateThread(c, input, enableWrites)
		case "reply_to_comment":
			return handlePRReplyToComment(c, input, enableWrites)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handlePRList(c *devops.Client, input ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'list' action")
	}
	prs, err := c.ListPullRequests(input.ProjectKey, input.RepoID, input.Status, input.Top)
	if err != nil {
		return resultError(fmt.Sprintf("listing PRs: %v", err))
	}
	return resultJSON(prs)
}

func handlePRGet(c *devops.Client, input ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" || input.PRID == 0 {
		return resultError("repo_id and pr_id are required for 'get' action")
	}
	pr, err := c.GetPullRequest(input.ProjectKey, input.RepoID, input.PRID)
	if err != nil {
		return resultError(fmt.Sprintf("getting PR: %v", err))
	}
	return resultJSON(pr)
}

func handlePRCreate(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
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
}

func handlePRUpdate(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
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
}

func handlePRAddComment(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
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
}

func handlePRListComments(c *devops.Client, input ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" || input.PRID == 0 {
		return resultError("repo_id and pr_id are required for 'list_comments' action")
	}
	threads, err := c.ListPRThreads(input.ProjectKey, input.RepoID, input.PRID)
	if err != nil {
		return resultError(fmt.Sprintf("listing PR comments: %v", err))
	}
	return resultJSON(threads)
}

func handlePRVote(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
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
}

func handlePRListReviewers(c *devops.Client, input ManagePullRequestsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" || input.PRID == 0 {
		return resultError("repo_id and pr_id are required for 'list_reviewers' action")
	}
	reviewers, err := c.ListPRReviewers(input.ProjectKey, input.RepoID, input.PRID)
	if err != nil {
		return resultError(fmt.Sprintf("listing reviewers: %v", err))
	}
	return resultJSON(reviewers)
}

func handlePRUpdateReviewers(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("update_reviewers action requires ADTK_ENABLE_WRITES=true")
	}
	if input.RepoID == "" || input.PRID == 0 || input.ReviewerIDs == "" {
		return resultError("repo_id, pr_id, and reviewer_ids are required for 'update_reviewers' action")
	}
	action := input.Status // reuse status field: "add" or "remove"
	if action == "" {
		action = "add"
	}
	ids := parseCSV(input.ReviewerIDs)
	if err := c.UpdatePRReviewers(input.ProjectKey, input.RepoID, input.PRID, ids, action); err != nil {
		return resultError(fmt.Sprintf("updating reviewers: %v", err))
	}
	return resultText(fmt.Sprintf("Reviewers updated (%s) for PR %d", action, input.PRID))
}

func handlePRCreateThread(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("create_thread action requires ADTK_ENABLE_WRITES=true")
	}
	if input.RepoID == "" || input.PRID == 0 || input.Comment == "" {
		return resultError("repo_id, pr_id, and comment are required for 'create_thread' action")
	}
	thread, err := c.CreatePRThread(input.ProjectKey, input.RepoID, input.PRID, input.Comment, input.FilePath, input.Line)
	if err != nil {
		return resultError(fmt.Sprintf("creating PR thread: %v", err))
	}
	return resultJSON(thread)
}

func handlePRUpdateThread(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("update_thread action requires ADTK_ENABLE_WRITES=true")
	}
	if input.RepoID == "" || input.PRID == 0 || input.ThreadID == 0 || input.Status == "" {
		return resultError("repo_id, pr_id, thread_id, and status are required for 'update_thread' action")
	}
	if err := c.UpdatePRThread(input.ProjectKey, input.RepoID, input.PRID, input.ThreadID, input.Status); err != nil {
		return resultError(fmt.Sprintf("updating PR thread: %v", err))
	}
	return resultText(fmt.Sprintf("Thread %d status updated to %q", input.ThreadID, input.Status))
}

func handlePRReplyToComment(c *devops.Client, input ManagePullRequestsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("reply_to_comment action requires ADTK_ENABLE_WRITES=true")
	}
	if input.RepoID == "" || input.PRID == 0 || input.ThreadID == 0 || input.Comment == "" {
		return resultError("repo_id, pr_id, thread_id, and comment are required for 'reply_to_comment' action")
	}
	comment, err := c.ReplyToComment(input.ProjectKey, input.RepoID, input.PRID, input.ThreadID, input.Comment)
	if err != nil {
		return resultError(fmt.Sprintf("replying to comment: %v", err))
	}
	return resultJSON(comment)
}

// parseCSV splits a comma-separated string into trimmed non-empty parts.
func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
