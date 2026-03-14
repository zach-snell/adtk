package mcp

import (
	"context"
	"fmt"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageMetricsInput defines the input schema for the manage_metrics tool.
type ManageMetricsInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'get_metrics'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	WorkItemID int    `json:"work_item_id,omitempty" jsonschema:"Work item ID (required for get_metrics)"`
}

// ManageMetricsHandler returns the handler for the manage_metrics tool.
func ManageMetricsHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageMetricsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageMetricsInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "get_metrics":
			return handleGetMetrics(c, input)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleGetMetrics(c *devops.Client, input ManageMetricsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'get_metrics' action")
	}
	metrics, err := c.ComputeWorkItemMetrics(input.ProjectKey, input.WorkItemID)
	if err != nil {
		return resultError(fmt.Sprintf("computing metrics for work item %d: %v", input.WorkItemID, err))
	}

	// Format durations for readability
	result := map[string]interface{}{
		"work_item_id":       input.WorkItemID,
		"current_status":     metrics.CurrentStatus,
		"cycle_time":         formatDuration(metrics.CycleTime),
		"lead_time":          formatDuration(metrics.LeadTime),
		"time_in_status":     formatTimeInStatus(metrics.TimeInStatus),
		"status_transitions": metrics.StatusTransitions,
	}
	return resultJSON(result)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}
	hours := int(d.Hours())
	if hours >= 24 {
		days := hours / 24
		remaining := hours % 24
		if remaining > 0 {
			return fmt.Sprintf("%dd %dh", days, remaining)
		}
		return fmt.Sprintf("%dd", days)
	}
	return d.String()
}

func formatTimeInStatus(m map[string]time.Duration) map[string]string {
	result := make(map[string]string, len(m))
	for status, d := range m {
		result[status] = formatDuration(d)
	}
	return result
}
