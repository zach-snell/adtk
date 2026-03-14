package mcp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageReposHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":"r1","name":"repo1","remoteUrl":"https://test/repo1"},{"id":"r2","name":"repo2","remoteUrl":"https://test/repo2"}]}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "list",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "repo1")
	assertResultSuccess(t, result, "repo2")
}

func TestManageReposHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":"r1","name":"repo1","defaultBranch":"refs/heads/main","remoteUrl":"https://test/repo1","size":12345}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "get",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "repo1")
	assertResultSuccess(t, result, "12345")
}

func TestManageReposHandler_Get_MissingRepoID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action: "get",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "repo_id is required")
}

func TestManageReposHandler_ListBranches(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"name":"refs/heads/main","objectId":"abc123"},{"name":"refs/heads/dev","objectId":"def456"}]}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "list_branches",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "refs/heads/main")
	assertResultSuccess(t, result, "refs/heads/dev")
}

func TestManageReposHandler_GetFile(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"objectId":"abc","path":"/README.md","content":"# Hello","gitObjectType":"blob"}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "get_file",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
		FilePath:   "/README.md",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Hello")
}

func TestManageReposHandler_GetFile_MissingPath(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "get_file",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "file_path is required")
}

func TestManageReposHandler_GetTree(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"path":"/src","gitObjectType":"tree"},{"path":"/README.md","gitObjectType":"blob"}]}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "get_tree",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "/src")
	assertResultSuccess(t, result, "/README.md")
}

func TestManageReposHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action: "invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}

func TestManageReposHandler_CreateBranch(t *testing.T) {
	t.Parallel()
	// create_branch does: list branches to find source OID → POST refs
	c := newTestClient(t, muxHandler(map[string]http.HandlerFunc{
		"/test-org/Proj/_apis/git/repositories/repo1/refs": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"count":1,"value":[{"name":"refs/heads/main","objectId":"abc123def456"}]}`))
				return
			}
			if r.Method == "POST" {
				body, _ := io.ReadAll(r.Body)
				var refs []map[string]string
				_ = json.Unmarshal(body, &refs)
				if len(refs) == 0 {
					t.Fatal("expected ref create body")
				}
				if refs[0]["newObjectId"] != "abc123def456" {
					t.Errorf("expected newObjectId from source branch, got %q", refs[0]["newObjectId"])
				}
				if !strings.Contains(refs[0]["name"], "refs/heads/feature-branch") {
					t.Errorf("expected branch name with refs/heads/ prefix, got %q", refs[0]["name"])
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"value":[{"name":"refs/heads/feature-branch"}]}`))
				return
			}
		},
	}))
	handler := ManageReposHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:       "create_branch",
		ProjectKey:   "Proj",
		RepoID:       "repo1",
		BranchName:   "feature-branch",
		SourceBranch: "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "feature-branch")
	assertResultSuccess(t, result, "main")
}

func TestManageReposHandler_CreateBranch_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:       "create_branch",
		ProjectKey:   "Proj",
		RepoID:       "repo1",
		BranchName:   "feature",
		SourceBranch: "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageReposHandler_CreateBranch_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "create_branch",
		ProjectKey: "Proj",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "branch_name, and source_branch are required")
}

func TestManageReposHandler_SearchCommits(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"commitId":"abc123","comment":"fix bug","author":{"name":"Jane"}},{"commitId":"def456","comment":"add feature","author":{"name":"John"}}]}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "search_commits",
		ProjectKey: "Proj",
		RepoID:     "repo1",
		Author:     "Jane",
		FromDate:   "2024-01-01",
		ToDate:     "2024-12-31",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "abc123")
	assertResultSuccess(t, result, "fix bug")
}

func TestManageReposHandler_SearchCommits_MissingRepoID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action: "search_commits",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "repo_id is required")
}

func TestManageReposHandler_ListPolicies(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"id":1,"type":{"displayName":"Minimum number of reviewers"},"isEnabled":true,"settings":{"minimumApproverCount":2,"scope":[{"repositoryId":"repo1"}]}}]}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "list_policies",
		ProjectKey: "Proj",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Minimum number of reviewers")
}

func TestManageReposHandler_ListTags(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"name":"refs/tags/v1.0.0","objectId":"aaa111"},{"name":"refs/tags/v1.1.0","objectId":"bbb222"}]}`))
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "list_tags",
		ProjectKey: "Proj",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "refs/tags/v1.0.0")
	assertResultSuccess(t, result, "refs/tags/v1.1.0")
}

func TestManageReposHandler_ListTags_MissingRepoID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action: "list_tags",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "repo_id is required")
}

func TestManageReposHandler_CreateTag(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"name":"refs/tags/v2.0.0","newObjectId":"sha123"}]}`))
	handler := ManageReposHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:     "create_tag",
		ProjectKey: "Proj",
		RepoID:     "repo1",
		TagName:    "v2.0.0",
		CommitSHA:  "sha123",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "v2.0.0")
	assertResultSuccess(t, result, "sha123")
}

func TestManageReposHandler_CreateTag_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action:    "create_tag",
		RepoID:    "repo1",
		TagName:   "v1.0",
		CommitSHA: "abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageReposHandler_CreateTag_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageReposHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action: "create_tag",
		RepoID: "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "tag_name, and commit_sha are required")
}
