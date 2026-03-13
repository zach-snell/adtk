package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManagePullRequestsHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"pullRequestId":1,"title":"Test PR","status":"active","sourceRefName":"refs/heads/feature","targetRefName":"refs/heads/main"}]}`))
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "list",
		RepoID: "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Test PR")
}

func TestManagePullRequestsHandler_List_MissingRepo(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "list",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "repo_id is required")
}

func TestManagePullRequestsHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"pullRequestId":42,"title":"My PR","status":"active","description":"A good PR"}`))
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "get",
		RepoID: "repo1",
		PRID:   42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "My PR")
}

func TestManagePullRequestsHandler_Get_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "get",
		RepoID: "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "pr_id are required")
}

func TestManagePullRequestsHandler_Create_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:       "create",
		RepoID:       "repo1",
		Title:        "New PR",
		SourceBranch: "feature",
		TargetBranch: "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePullRequestsHandler_Create(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"pullRequestId":99,"title":"New PR","status":"active"}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:       "create",
		RepoID:       "repo1",
		Title:        "New PR",
		SourceBranch: "feature",
		TargetBranch: "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "New PR")
}

func TestManagePullRequestsHandler_ListComments(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":1,"comments":[{"id":1,"content":"Great work!"}],"status":"active"}]}`))
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "list_comments",
		RepoID: "repo1",
		PRID:   42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Great work!")
}

func TestManagePullRequestsHandler_ListReviewers(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":"r1","displayName":"Reviewer One","vote":10}]}`))
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "list_reviewers",
		RepoID: "repo1",
		PRID:   42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Reviewer One")
}

func TestManagePullRequestsHandler_Vote_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:     "vote",
		RepoID:     "repo1",
		PRID:       42,
		ReviewerID: "r1",
		Vote:       10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePullRequestsHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
