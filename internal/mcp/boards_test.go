package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageBoardsHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":"b1","name":"Stories"},{"id":"b2","name":"Bugs"}]}`))
	handler := ManageBoardsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageBoardsInput{
		Action:     "list",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Stories")
	assertResultSuccess(t, result, "Bugs")
}

func TestManageBoardsHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":"b1","name":"Stories","columns":[{"name":"New"},{"name":"Done"}]}`))
	handler := ManageBoardsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageBoardsInput{
		Action:     "get",
		ProjectKey: "TestProject",
		BoardID:    "Stories",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Stories")
}

func TestManageBoardsHandler_Get_MissingBoardID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageBoardsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageBoardsInput{
		Action:     "get",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "board_id is required")
}

func TestManageBoardsHandler_GetColumns(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":3,"value":[{"id":"c1","name":"New","columnType":"incoming"},{"id":"c2","name":"Active","columnType":"inProgress"},{"id":"c3","name":"Done","columnType":"outgoing"}]}`))
	handler := ManageBoardsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageBoardsInput{
		Action:     "get_columns",
		ProjectKey: "TestProject",
		BoardID:    "Stories",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "New")
	assertResultSuccess(t, result, "Active")
	assertResultSuccess(t, result, "Done")
}

func TestManageBoardsHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageBoardsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageBoardsInput{
		Action: "list",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageBoardsHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageBoardsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageBoardsInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
