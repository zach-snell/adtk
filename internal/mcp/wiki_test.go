package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageWikiHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":"w1","name":"ProjectWiki","type":"projectWiki"}]}`))
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "list",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "ProjectWiki")
}

func TestManageWikiHandler_GetPage(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":1,"path":"/Home","content":"# Welcome","gitItemPath":"/Home.md"}`))
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "get_page",
		ProjectKey: "TestProject",
		WikiID:     "ProjectWiki",
		PagePath:   "/Home",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Welcome")
}

func TestManageWikiHandler_GetPage_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "get_page",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "wiki_id and page_path are required")
}

func TestManageWikiHandler_CreatePage_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "create_page",
		ProjectKey: "TestProject",
		WikiID:     "w1",
		PagePath:   "/NewPage",
		Content:    "# New",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWikiHandler_CreatePage(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":2,"path":"/NewPage","content":"# New"}`))
	handler := ManageWikiHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "create_page",
		ProjectKey: "TestProject",
		WikiID:     "w1",
		PagePath:   "/NewPage",
		Content:    "# New",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "NewPage")
}

func TestManageWikiHandler_UpdatePage_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "update_page",
		ProjectKey: "TestProject",
		WikiID:     "w1",
		PagePath:   "/Home",
		Content:    "# Updated",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWikiHandler_DeletePage_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "delete_page",
		ProjectKey: "TestProject",
		WikiID:     "w1",
		PagePath:   "/OldPage",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWikiHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action: "list",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageWikiHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}

func TestManageWikiHandler_ListPages(t *testing.T) {
	t.Parallel()
	// The API returns a nested tree. The handler flattens it.
	c := newTestClient(t, jsonHandler(`{"id":1,"path":"/","subPages":[{"id":2,"path":"/Home","subPages":[{"id":3,"path":"/Home/Setup"}]},{"id":4,"path":"/Design"}]}`))
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "list_pages",
		ProjectKey: "TestProject",
		WikiID:     "ProjectWiki",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Verify the nested tree is flattened to a list
	assertResultSuccess(t, result, "/Home")
	assertResultSuccess(t, result, "/Home/Setup")
	assertResultSuccess(t, result, "/Design")
}

func TestManageWikiHandler_ListPages_MissingWikiID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWikiHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "list_pages",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "wiki_id is required")
}

func TestManageWikiHandler_DeletePage(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{}`))
	handler := ManageWikiHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWikiInput{
		Action:     "delete_page",
		ProjectKey: "TestProject",
		WikiID:     "w1",
		PagePath:   "/OldPage",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "deleted successfully")
}
