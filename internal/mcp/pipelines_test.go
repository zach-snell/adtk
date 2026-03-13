package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManagePipelinesHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":1,"name":"CI","folder":"\\"},{"id":2,"name":"Deploy","folder":"\\prod"}]}`))
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "list",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "CI")
	assertResultSuccess(t, result, "Deploy")
}

func TestManagePipelinesHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":1,"name":"CI","folder":"\\","url":"https://test"}`))
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "get",
		ProjectKey: "TestProject",
		PipelineID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "CI")
}

func TestManagePipelinesHandler_Get_MissingID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "get",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "pipeline_id is required")
}

func TestManagePipelinesHandler_ListRuns(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":100,"name":"CI #100","state":"completed","result":"succeeded"}]}`))
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "list_runs",
		ProjectKey: "TestProject",
		PipelineID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "CI #100")
	assertResultSuccess(t, result, "succeeded")
}

func TestManagePipelinesHandler_GetRun(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":100,"name":"CI #100","state":"completed","result":"succeeded"}`))
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "get_run",
		ProjectKey: "TestProject",
		PipelineID: 1,
		RunID:      100,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "succeeded")
}

func TestManagePipelinesHandler_Trigger_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "trigger",
		ProjectKey: "TestProject",
		PipelineID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManagePipelinesHandler_Trigger(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":101,"name":"CI #101","state":"inProgress"}`))
	handler := ManagePipelinesHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "trigger",
		ProjectKey: "TestProject",
		PipelineID: 1,
		Branch:     "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "CI #101")
}

func TestManagePipelinesHandler_GetLogs(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"logs":[{"id":1},{"id":2}]}`))
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "get_logs",
		ProjectKey: "TestProject",
		PipelineID: 1,
		RunID:      100,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "logs")
}

func TestManagePipelinesHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action: "list",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManagePipelinesHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManagePipelinesHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManagePipelinesInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
