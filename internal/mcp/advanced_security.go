package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageAdvancedSecurityInput defines the input schema for the manage_advanced_security tool.
type ManageAdvancedSecurityInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'list_alerts', 'get_alert'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	RepoID     string `json:"repo_id,omitempty" jsonschema:"Repository name or ID (required)"`
	AlertID    int    `json:"alert_id,omitempty" jsonschema:"Alert ID (required for get_alert)"`
	States     string `json:"states,omitempty" jsonschema:"Filter by alert states (for list_alerts)"`
	Severities string `json:"severities,omitempty" jsonschema:"Filter by severities (for list_alerts)"`
}

// ManageAdvancedSecurityHandler returns the handler for the manage_advanced_security tool.
func ManageAdvancedSecurityHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageAdvancedSecurityInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageAdvancedSecurityInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}
		if input.RepoID == "" {
			return resultError("repo_id is required")
		}

		switch input.Action {
		case "list_alerts":
			alerts, err := c.GetSecurityAlerts(input.ProjectKey, input.RepoID, input.States, input.Severities)
			if err != nil {
				return resultError(fmt.Sprintf("listing security alerts: %v", err))
			}
			return resultJSON(alerts)
		case "get_alert":
			if input.AlertID == 0 {
				return resultError("alert_id is required for 'get_alert' action")
			}
			alert, err := c.GetSecurityAlertDetails(input.ProjectKey, input.RepoID, input.AlertID)
			if err != nil {
				return resultError(fmt.Sprintf("getting security alert: %v", err))
			}
			return resultJSON(alert)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
