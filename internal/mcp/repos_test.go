package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageReposHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":"r1","name":"repo1","remoteUrl":"https://test/repo1"},{"id":"r2","name":"repo2","remoteUrl":"https://test/repo2"}]}`))
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

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
	handler := ManageReposHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageReposInput{
		Action: "invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
