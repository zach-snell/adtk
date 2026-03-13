package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageUsersInput defines the input schema for the manage_users tool.
type ManageUsersInput struct {
	Action string `json:"action" jsonschema:"Action to perform: 'get_current', 'search'"`
	Query  string `json:"query,omitempty" jsonschema:"Search query - display name or email (for search)"`
}

// ManageUsersHandler returns the handler for the manage_users tool.
func ManageUsersHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageUsersInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageUsersInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "get_current":
			return handleGetCurrentUser(c)
		case "search":
			if input.Query == "" {
				return resultError("query is required for 'search' action")
			}
			return handleSearchUsers(c, input.Query)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleGetCurrentUser(c *devops.Client) (*sdkmcp.CallToolResult, any, error) {
	data, err := c.GetPreview("", "/connectionData", nil)
	if err != nil {
		return resultError(fmt.Sprintf("getting current user: %v", err))
	}
	var result devops.ConnectionData
	if err := json.Unmarshal(data, &result); err != nil {
		return resultError(fmt.Sprintf("parsing current user: %v", err))
	}
	return resultJSON(result.AuthenticatedUser)
}

func handleSearchUsers(c *devops.Client, query string) (*sdkmcp.CallToolResult, any, error) {
	params := url.Values{}
	params.Set("searchFilter", "General")
	params.Set("filterValue", query)
	data, err := c.GetIdentity("/graph/users", params)
	if err != nil {
		return resultError(fmt.Sprintf("searching users: %v", err))
	}
	return resultText(string(data))
}
