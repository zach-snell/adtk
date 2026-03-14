package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageSearchInput defines the input schema for the manage_search tool.
type ManageSearchInput struct {
	Action     string   `json:"action" jsonschema:"Action to perform: 'wiql', 'code', 'work_items', 'wiki', 'get_query', 'run_query'"`
	ProjectKey string   `json:"project_key,omitempty" jsonschema:"Project name (optional, scopes search)"`
	Query      string   `json:"query,omitempty" jsonschema:"WIQL query string (for wiql action) or search text (for code, work_items, wiki)"`
	QueryID    string   `json:"query_id,omitempty" jsonschema:"Saved query ID or path (for get_query, run_query)"`
	Top        int      `json:"top,omitempty" jsonschema:"Max results to return (default 25)"`
	Fields     []string `json:"fields,omitempty" jsonschema:"Fields to return for WIQL results"`
}

// ManageSearchHandler returns the handler for the manage_search tool.
func ManageSearchHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageSearchInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageSearchInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "get_query":
			return handleGetQuery(c, input)
		case "run_query":
			return handleRunQuery(c, input)
		case "wiql", "code", "work_items", "wiki":
			if input.Query == "" {
				return resultError("query is required")
			}
			return handleSearchAction(c, input)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleSearchAction(c *devops.Client, input ManageSearchInput) (*sdkmcp.CallToolResult, any, error) {
	switch input.Action {
	case "wiql":
		items, err := c.WIQLAndFetch(input.ProjectKey, input.Query, input.Fields, input.Top)
		if err != nil {
			return resultError(fmt.Sprintf("WIQL search: %v", err))
		}
		flat := make([]map[string]interface{}, len(items))
		for i, wi := range items {
			flat[i] = flattenWorkItem(&wi)
		}
		return resultJSON(flat)
	case "code":
		result, err := c.SearchCode(input.ProjectKey, input.Query, input.Top)
		if err != nil {
			return resultError(fmt.Sprintf("code search: %v", err))
		}
		return resultJSON(result)
	case "work_items":
		result, err := c.SearchWorkItems(input.ProjectKey, input.Query, input.Top)
		if err != nil {
			return resultError(fmt.Sprintf("work item search: %v", err))
		}
		return resultJSON(result)
	case "wiki":
		data, err := c.SearchWiki(input.ProjectKey, input.Query, input.Top)
		if err != nil {
			return resultError(fmt.Sprintf("wiki search: %v", err))
		}
		return resultText(string(data))
	default:
		return resultError(fmt.Sprintf("unknown action: %s", input.Action))
	}
}

func handleGetQuery(c *devops.Client, input ManageSearchInput) (*sdkmcp.CallToolResult, any, error) {
	if input.QueryID == "" {
		return resultError("query_id is required for 'get_query' action")
	}
	result, err := c.GetQuery(input.ProjectKey, input.QueryID)
	if err != nil {
		return resultError(fmt.Sprintf("getting query: %v", err))
	}
	return resultJSON(result)
}

func handleRunQuery(c *devops.Client, input ManageSearchInput) (*sdkmcp.CallToolResult, any, error) {
	if input.QueryID == "" {
		return resultError("query_id is required for 'run_query' action")
	}
	items, err := c.RunQueryByID(input.ProjectKey, input.QueryID, input.Top)
	if err != nil {
		return resultError(fmt.Sprintf("running query: %v", err))
	}
	flat := make([]map[string]interface{}, len(items))
	for i, wi := range items {
		flat[i] = flattenWorkItem(&wi)
	}
	return resultJSON(flat)
}
