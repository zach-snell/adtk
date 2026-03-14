package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageMetricsHandler_GetMetrics(t *testing.T) {
	t.Parallel()
	// Mock work item updates with state transitions:
	// Rev 1: Created (New) at T0
	// Rev 2: New -> Active at T0+1h
	// Rev 3: Active -> Closed at T0+25h
	now := time.Now().UTC()
	t0 := now.Add(-48 * time.Hour).Format(time.RFC3339)
	t1 := now.Add(-47 * time.Hour).Format(time.RFC3339)
	t2 := now.Add(-23 * time.Hour).Format(time.RFC3339)

	c := newTestClient(t, jsonHandler(`{"count":3,"value":[`+
		`{"id":1,"rev":1,"revisedDate":"`+t0+`","fields":{"System.State":{"oldValue":"","newValue":"New"}}},`+
		`{"id":1,"rev":2,"revisedDate":"`+t1+`","fields":{"System.State":{"oldValue":"New","newValue":"Active"}}},`+
		`{"id":1,"rev":3,"revisedDate":"`+t2+`","fields":{"System.State":{"oldValue":"Active","newValue":"Closed"}}}`+
		`]}`))
	handler := ManageMetricsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageMetricsInput{
		Action:     "get_metrics",
		ProjectKey: "TestProject",
		WorkItemID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	// Verify key metrics are present in the response
	assertResultSuccess(t, result, `"current_status": "Closed"`)
	assertResultSuccess(t, result, `"cycle_time"`)
	assertResultSuccess(t, result, `"lead_time"`)
	assertResultSuccess(t, result, `"time_in_status"`)
	assertResultSuccess(t, result, `"status_transitions"`)
}

func TestManageMetricsHandler_GetMetrics_NoStateChanges(t *testing.T) {
	t.Parallel()
	// No state field changes — only other field updates
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":1,"rev":1,"revisedDate":"2024-01-01T00:00:00Z","fields":{"System.Title":{"oldValue":"","newValue":"My Item"}}}]}`))
	handler := ManageMetricsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageMetricsInput{
		Action:     "get_metrics",
		ProjectKey: "TestProject",
		WorkItemID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	// With no state changes, cycle_time and lead_time should be "0s"
	assertResultSuccess(t, result, `"cycle_time": "0s"`)
	assertResultSuccess(t, result, `"lead_time": "0s"`)
}

func TestManageMetricsHandler_GetMetrics_MissingWorkItemID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageMetricsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageMetricsInput{
		Action:     "get_metrics",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "work_item_id is required")
}

func TestManageMetricsHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageMetricsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageMetricsInput{
		Action:     "get_metrics",
		WorkItemID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageMetricsHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageMetricsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageMetricsInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}

func TestManageMetricsHandler_APIError(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, errorHandler(500))
	handler := ManageMetricsHandler(c)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageMetricsInput{
		Action:     "get_metrics",
		ProjectKey: "TestProject",
		WorkItemID: 999,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "computing metrics")
}

func TestFormatDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"zero", 0, "0s"},
		{"hours", 5 * time.Hour, "5h0m0s"},
		{"one day", 24 * time.Hour, "1d"},
		{"days and hours", 50 * time.Hour, "2d 2h"},
		{"minutes", 30 * time.Minute, "30m0s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatDuration(tt.d)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

func TestFormatTimeInStatus(t *testing.T) {
	t.Parallel()
	input := map[string]time.Duration{
		"Active": 48 * time.Hour,
		"New":    2 * time.Hour,
	}
	result := formatTimeInStatus(input)
	if result["Active"] != "2d" {
		t.Errorf("expected Active = 2d, got %q", result["Active"])
	}
	if result["New"] != "2h0m0s" {
		t.Errorf("expected New = 2h0m0s, got %q", result["New"])
	}
}
