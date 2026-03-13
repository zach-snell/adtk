package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageIterationsInput defines the input schema for the manage_iterations tool.
type ManageIterationsInput struct {
	Action      string `json:"action" jsonschema:"required,description=Action to perform: 'list'\\, 'get'\\, 'get_current'"`
	ProjectKey  string `json:"project_key,omitempty" jsonschema:"description=Project name (required)"`
	Team        string `json:"team,omitempty" jsonschema:"description=Team name (optional\\, scopes to a specific team)"`
	IterationID string `json:"iteration_id,omitempty" jsonschema:"description=Iteration ID (required for get)"`
}

// ManageIterationsHandler returns the handler for the manage_iterations tool.
func ManageIterationsHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "list":
			iters, err := c.ListIterations(input.ProjectKey, input.Team)
			if err != nil {
				return resultError(fmt.Sprintf("listing iterations: %v", err))
			}
			return resultJSON(iters)
		case "get":
			if input.IterationID == "" {
				return resultError("iteration_id is required for 'get' action")
			}
			iter, err := c.GetIteration(input.ProjectKey, input.Team, input.IterationID)
			if err != nil {
				return resultError(fmt.Sprintf("getting iteration: %v", err))
			}
			return resultJSON(iter)
		case "get_current":
			iter, err := c.GetCurrentIteration(input.ProjectKey, input.Team)
			if err != nil {
				return resultError(fmt.Sprintf("getting current iteration: %v", err))
			}
			return resultJSON(iter)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
