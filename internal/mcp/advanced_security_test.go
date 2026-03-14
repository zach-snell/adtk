package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageAdvancedSecurityHandler_ListAlerts(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"alertId":1,"alertType":"dependency","title":"CVE-2024-1234 in lodash","severity":"high","state":"active"},{"alertId":2,"alertType":"secret","title":"Exposed API key","severity":"critical","state":"active"}]}`))
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "list_alerts",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "CVE-2024-1234")
	assertResultSuccess(t, result, "Exposed API key")
}

func TestManageAdvancedSecurityHandler_ListAlerts_WithFilters(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"alertId":2,"alertType":"secret","title":"Exposed API key","severity":"critical","state":"active"}]}`))
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "list_alerts",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
		States:     "active",
		Severities: "critical",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Exposed API key")
	assertResultSuccess(t, result, "critical")
}

func TestManageAdvancedSecurityHandler_GetAlert(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"alertId":1,"alertType":"dependency","title":"CVE-2024-1234 in lodash","severity":"high","state":"active","firstSeenDate":"2024-06-01","fixedDate":null,"dismissal":null,"logicalLocations":[{"fullyQualifiedName":"package.json"}]}`))
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "get_alert",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
		AlertID:    1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "CVE-2024-1234")
	assertResultSuccess(t, result, "package.json")
}

func TestManageAdvancedSecurityHandler_GetAlert_MissingAlertID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "get_alert",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "alert_id is required")
}

func TestManageAdvancedSecurityHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action: "list_alerts",
		RepoID: "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageAdvancedSecurityHandler_MissingRepoID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "list_alerts",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "repo_id is required")
}

func TestManageAdvancedSecurityHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}

func TestManageAdvancedSecurityHandler_APIError(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, errorHandler(500))
	handler := ManageAdvancedSecurityHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageAdvancedSecurityInput{
		Action:     "list_alerts",
		ProjectKey: "TestProject",
		RepoID:     "repo1",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "listing security alerts")
}
