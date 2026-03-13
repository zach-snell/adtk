package mcp

import (
	"context"
	"net/http"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageSearchHandler_WIQL(t *testing.T) {
	t.Parallel()
	// WIQL needs 2-step: first POST returns IDs, then POST batch returns items.
	// With a single handler, we return the WIQL result with workItems then the batch response.
	c := newTestClient(t, muxHandler(map[string]http.HandlerFunc{
		"/test-org/_apis/wit/wiql":           jsonHandler(`{"queryType":"flat","queryResultType":"workItem","workItems":[{"id":1},{"id":2}]}`),
		"/test-org/_apis/wit/workitemsbatch": jsonHandler(`{"count":2,"value":[{"id":1,"rev":1,"fields":{"System.Title":"First"}},{"id":2,"rev":1,"fields":{"System.Title":"Second"}}]}`),
	}))
	handler := ManageSearchHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageSearchInput{
		Action: "wiql",
		Query:  "SELECT [System.Id] FROM WorkItems",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "First")
}

func TestManageSearchHandler_Code(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"results":[{"fileName":"main.go","path":"/src/main.go","repository":{"name":"repo1"},"project":{"name":"proj1"}}]}`))
	handler := ManageSearchHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageSearchInput{
		Action: "code",
		Query:  "func main",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "main.go")
}

func TestManageSearchHandler_WorkItems(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"results":[{"project":{"name":"proj"},"fields":{"system.title":"Found Item"}}]}`))
	handler := ManageSearchHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageSearchInput{
		Action: "work_items",
		Query:  "bug",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Found Item")
}

func TestManageSearchHandler_Wiki(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"results":[{"title":"Wiki Page","path":"/Home"}]}`))
	handler := ManageSearchHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageSearchInput{
		Action: "wiki",
		Query:  "home",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Wiki Page")
}

func TestManageSearchHandler_MissingQuery(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageSearchHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageSearchInput{
		Action: "code",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "query is required")
}

func TestManageSearchHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageSearchHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageSearchInput{
		Action: "invalid",
		Query:  "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
