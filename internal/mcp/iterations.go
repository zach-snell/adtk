package mcp

import (
	"context"
	"fmt"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageIterationsInput defines the input schema for the manage_iterations tool.
type ManageIterationsInput struct {
	Action      string `json:"action" jsonschema:"Action to perform: 'list', 'get', 'get_current', 'create', 'get_team_settings'"`
	ProjectKey  string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	Team        string `json:"team,omitempty" jsonschema:"Team name (optional, scopes to a specific team)"`
	IterationID string `json:"iteration_id,omitempty" jsonschema:"Iteration ID (required for get)"`
	Name        string `json:"name,omitempty" jsonschema:"Iteration name (required for create)"`
	StartDate   string `json:"start_date,omitempty" jsonschema:"Start date in YYYY-MM-DD format (optional, for create)"`
	FinishDate  string `json:"finish_date,omitempty" jsonschema:"Finish date in YYYY-MM-DD format (optional, for create)"`
}

// ManageIterationsHandler returns the handler for the manage_iterations tool.
func ManageIterationsHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "list":
			return handleIterationList(c, input)
		case "get":
			return handleIterationGet(c, input)
		case "get_current":
			return handleIterationGetCurrent(c, input)
		case "create":
			return handleIterationCreate(c, input, enableWrites)
		case "get_team_settings":
			return handleIterationGetTeamSettings(c, input)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleIterationList(c *devops.Client, input ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
	iters, err := c.ListIterations(input.ProjectKey, input.Team)
	if err != nil {
		return resultError(fmt.Sprintf("listing iterations: %v", err))
	}
	return resultJSON(iters)
}

func handleIterationGet(c *devops.Client, input ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.IterationID == "" {
		return resultError("iteration_id is required for 'get' action")
	}
	iter, err := c.GetIteration(input.ProjectKey, input.Team, input.IterationID)
	if err != nil {
		return resultError(fmt.Sprintf("getting iteration: %v", err))
	}
	return resultJSON(iter)
}

func handleIterationGetCurrent(c *devops.Client, input ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
	iter, err := c.GetCurrentIteration(input.ProjectKey, input.Team)
	if err != nil {
		return resultError(fmt.Sprintf("getting current iteration: %v", err))
	}
	return resultJSON(iter)
}

func handleIterationCreate(c *devops.Client, input ManageIterationsInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("create action requires ADTK_ENABLE_WRITES=true")
	}
	if input.Name == "" {
		return resultError("name is required for 'create' action")
	}

	var startDate, finishDate *time.Time
	if input.StartDate != "" {
		t, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			return resultError(fmt.Sprintf("invalid start_date format (use YYYY-MM-DD): %v", err))
		}
		startDate = &t
	}
	if input.FinishDate != "" {
		t, err := time.Parse("2006-01-02", input.FinishDate)
		if err != nil {
			return resultError(fmt.Sprintf("invalid finish_date format (use YYYY-MM-DD): %v", err))
		}
		finishDate = &t
	}

	if err := c.CreateIteration(input.ProjectKey, input.Name, startDate, finishDate); err != nil {
		return resultError(fmt.Sprintf("creating iteration: %v", err))
	}
	return resultText(fmt.Sprintf("Iteration %q created successfully", input.Name))
}

func handleIterationGetTeamSettings(c *devops.Client, input ManageIterationsInput) (*sdkmcp.CallToolResult, any, error) {
	settings, err := c.GetTeamSettings(input.ProjectKey, input.Team)
	if err != nil {
		return resultError(fmt.Sprintf("getting team settings: %v", err))
	}
	return resultJSON(settings)
}
