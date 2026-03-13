package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageUsersHandler_GetCurrent(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"authenticatedUser":{"id":"u1","displayName":"Test User","uniqueName":"test@example.com"},"authorizedUser":{"id":"u1"},"instanceId":"inst1"}`))
	handler := ManageUsersHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageUsersInput{
		Action: "get_current",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Test User")
}

func TestManageUsersHandler_Search(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":"u1","displayName":"Found User"}]}`))
	handler := ManageUsersHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageUsersInput{
		Action: "search",
		Query:  "Found",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Found User")
}

func TestManageUsersHandler_Search_MissingQuery(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageUsersHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageUsersInput{
		Action: "search",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "query is required")
}

func TestManageUsersHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageUsersHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageUsersInput{
		Action: "invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
