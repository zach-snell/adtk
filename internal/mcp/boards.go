package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageBoardsInput defines the input schema for the manage_boards tool.
type ManageBoardsInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'list', 'get', 'get_columns'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	Team       string `json:"team,omitempty" jsonschema:"Team name (optional, scopes to a specific team)"`
	BoardID    string `json:"board_id,omitempty" jsonschema:"Board name or ID (required for get, get_columns)"`
}

// ManageBoardsHandler returns the handler for the manage_boards tool.
func ManageBoardsHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageBoardsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageBoardsInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "list":
			boards, err := c.ListBoards(input.ProjectKey, input.Team)
			if err != nil {
				return resultError(fmt.Sprintf("listing boards: %v", err))
			}
			return resultJSON(boards)
		case "get":
			if input.BoardID == "" {
				return resultError("board_id is required for 'get' action")
			}
			data, err := c.GetBoard(input.ProjectKey, input.Team, input.BoardID)
			if err != nil {
				return resultError(fmt.Sprintf("getting board: %v", err))
			}
			return resultText(string(data))
		case "get_columns":
			if input.BoardID == "" {
				return resultError("board_id is required for 'get_columns' action")
			}
			cols, err := c.GetBoardColumns(input.ProjectKey, input.Team, input.BoardID)
			if err != nil {
				return resultError(fmt.Sprintf("getting board columns: %v", err))
			}
			return resultJSON(cols)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
