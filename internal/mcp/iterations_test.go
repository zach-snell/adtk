package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageIterationsHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":"i1","name":"Sprint 1","path":"\\Project\\Sprint 1"},{"id":"i2","name":"Sprint 2","path":"\\Project\\Sprint 2"}]}`))
	handler := ManageIterationsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageIterationsInput{
		Action:     "list",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Sprint 1")
	assertResultSuccess(t, result, "Sprint 2")
}

func TestManageIterationsHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":"i1","name":"Sprint 1","path":"\\Project\\Sprint 1","attributes":{"timeFrame":"current"}}`))
	handler := ManageIterationsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageIterationsInput{
		Action:      "get",
		ProjectKey:  "TestProject",
		IterationID: "i1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Sprint 1")
}

func TestManageIterationsHandler_Get_MissingID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageIterationsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageIterationsInput{
		Action:     "get",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "iteration_id is required")
}

func TestManageIterationsHandler_GetCurrent(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":"i1","name":"Sprint 1","attributes":{"timeFrame":"past"}},{"id":"i2","name":"Sprint 2","attributes":{"timeFrame":"current"}}]}`))
	handler := ManageIterationsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageIterationsInput{
		Action:     "get_current",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Sprint 2")
}

func TestManageIterationsHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageIterationsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageIterationsInput{
		Action: "list",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageIterationsHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageIterationsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageIterationsInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
