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

func TestManagePullRequestsHandler_CreateThread(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":10,"comments":[{"id":1,"content":"Review comment"}],"status":"active"}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:  "create_thread",
		RepoID:  "repo1",
		PRID:    42,
		Comment: "Review comment",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Review comment")
}

func TestManagePullRequestsHandler_CreateThread_Inline(t *testing.T) {
	t.Parallel()
	// The handler passes filePath and line to CreatePRThread which includes them in the POST body.
	// The mock returns a thread response; we verify the thread was created successfully.
	c := newTestClient(t, jsonHandler(`{"id":11,"comments":[{"id":1,"content":"Line comment"}],"status":"active"}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "create_thread",
		RepoID:   "repo1",
		PRID:     42,
		Comment:  "Line comment",
		FilePath: "/src/main.go",
		Line:     25,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Line comment")
	assertResultSuccess(t, result, `"id": 11`)
}

func TestManagePullRequestsHandler_CreateThread_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:  "create_thread",
		RepoID:  "repo1",
		PRID:    42,
		Comment: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePullRequestsHandler_CreateThread_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "create_thread",
		RepoID: "repo1",
		PRID:   42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "comment are required")
}

func TestManagePullRequestsHandler_UpdateThread(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "update_thread",
		RepoID:   "repo1",
		PRID:     42,
		ThreadID: 10,
		Status:   "closed",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "closed")
}

func TestManagePullRequestsHandler_UpdateThread_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "update_thread",
		RepoID:   "repo1",
		PRID:     42,
		ThreadID: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "status are required")
}

func TestManagePullRequestsHandler_UpdateThread_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "update_thread",
		RepoID:   "repo1",
		PRID:     42,
		ThreadID: 10,
		Status:   "closed",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePullRequestsHandler_ReplyToComment(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":5,"content":"Thanks for the feedback","commentType":"text"}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "reply_to_comment",
		RepoID:   "repo1",
		PRID:     42,
		ThreadID: 10,
		Comment:  "Thanks for the feedback",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Thanks for the feedback")
}

func TestManagePullRequestsHandler_ReplyToComment_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "reply_to_comment",
		RepoID:   "repo1",
		PRID:     42,
		ThreadID: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "comment are required")
}

func TestManagePullRequestsHandler_ReplyToComment_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:   "reply_to_comment",
		RepoID:   "repo1",
		PRID:     42,
		ThreadID: 10,
		Comment:  "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePullRequestsHandler_UpdateReviewers(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:      "update_reviewers",
		RepoID:      "repo1",
		PRID:        42,
		ReviewerIDs: "user1,user2",
		Status:      "add",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "add")
}

func TestManagePullRequestsHandler_UpdateReviewers_Remove(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{}`))
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:      "update_reviewers",
		RepoID:      "repo1",
		PRID:        42,
		ReviewerIDs: "user1",
		Status:      "remove",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "remove")
}

func TestManagePullRequestsHandler_UpdateReviewers_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action:      "update_reviewers",
		RepoID:      "repo1",
		PRID:        42,
		ReviewerIDs: "user1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePullRequestsHandler_UpdateReviewers_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePullRequestsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePullRequestsInput{
		Action: "update_reviewers",
		RepoID: "repo1",
		PRID:   42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "reviewer_ids are required")
}

func TestParseCSV(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"simple", "a,b,c", 3},
		{"with spaces", " a , b , c ", 3},
		{"single", "a", 1},
		{"empty parts", "a,,b", 2},
		{"empty string", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseCSV(tt.input)
			if len(got) != tt.want {
				t.Errorf("parseCSV(%q) = %d items, want %d", tt.input, len(got), tt.want)
			}
		})
	}
}
