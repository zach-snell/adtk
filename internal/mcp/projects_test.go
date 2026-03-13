package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageProjectsHandler_List(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":"1","name":"ProjectA","state":"wellFormed"},{"id":"2","name":"ProjectB","state":"wellFormed"}]}`))
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action: "list",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "ProjectA")
	assertResultSuccess(t, result, "ProjectB")
}

func TestManageProjectsHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":"1","name":"ProjectA","state":"wellFormed","visibility":"private"}`))
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action:     "get",
		ProjectKey: "ProjectA",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "ProjectA")
}

func TestManageProjectsHandler_Get_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action: "get",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageProjectsHandler_ListTeams(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":"t1","name":"TeamAlpha"}]}`))
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action:     "list_teams",
		ProjectKey: "ProjectA",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "TeamAlpha")
}

func TestManageProjectsHandler_GetTeam(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":"t1","name":"TeamAlpha","description":"Primary team"}`))
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action:     "get_team",
		ProjectKey: "ProjectA",
		TeamID:     "TeamAlpha",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "TeamAlpha")
}

func TestManageProjectsHandler_GetTeam_MissingTeamID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action:     "get_team",
		ProjectKey: "ProjectA",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "team_id is required")
}

func TestManageProjectsHandler_Create_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action: "create",
		Name:   "NewProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageProjectsHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageProjectsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageProjectsInput{
		Action: "invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
