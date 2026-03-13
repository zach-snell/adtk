package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageWikiInput defines the input schema for the manage_wiki tool.
type ManageWikiInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'list', 'get_page', 'create_page', 'update_page', 'delete_page'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	WikiID     string `json:"wiki_id,omitempty" jsonschema:"Wiki name or ID (required for page operations)"`
	PagePath   string `json:"page_path,omitempty" jsonschema:"Wiki page path e.g. /Home or /Design/Architecture (required for page operations)"`
	Content    string `json:"content,omitempty" jsonschema:"Page content in Markdown (for create_page, update_page). Wiki is markdown-native."`
	Version    int    `json:"version,omitempty" jsonschema:"Page version for optimistic concurrency (for update_page)"`
}

// ManageWikiHandler returns the handler for the manage_wiki tool.
func ManageWikiHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageWikiInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageWikiInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "list":
			wikis, err := c.ListWikis(input.ProjectKey)
			if err != nil {
				return resultError(fmt.Sprintf("listing wikis: %v", err))
			}
			return resultJSON(wikis)
		case "get_page":
			if input.WikiID == "" || input.PagePath == "" {
				return resultError("wiki_id and page_path are required for 'get_page' action")
			}
			page, err := c.GetWikiPage(input.ProjectKey, input.WikiID, input.PagePath, true)
			if err != nil {
				return resultError(fmt.Sprintf("getting wiki page: %v", err))
			}
			return resultJSON(page)
		case "create_page":
			if !enableWrites {
				return resultError("create_page action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WikiID == "" || input.PagePath == "" || input.Content == "" {
				return resultError("wiki_id, page_path, and content are required for 'create_page' action")
			}
			page, err := c.CreateWikiPage(input.ProjectKey, input.WikiID, input.PagePath, input.Content)
			if err != nil {
				return resultError(fmt.Sprintf("creating wiki page: %v", err))
			}
			return resultJSON(page)
		case "update_page":
			if !enableWrites {
				return resultError("update_page action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WikiID == "" || input.PagePath == "" || input.Content == "" {
				return resultError("wiki_id, page_path, and content are required for 'update_page' action")
			}
			page, err := c.UpdateWikiPage(input.ProjectKey, input.WikiID, input.PagePath, input.Content, input.Version)
			if err != nil {
				return resultError(fmt.Sprintf("updating wiki page: %v", err))
			}
			return resultJSON(page)
		case "delete_page":
			if !enableWrites {
				return resultError("delete_page action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WikiID == "" || input.PagePath == "" {
				return resultError("wiki_id and page_path are required for 'delete_page' action")
			}
			if err := c.DeleteWikiPage(input.ProjectKey, input.WikiID, input.PagePath); err != nil {
				return resultError(fmt.Sprintf("deleting wiki page: %v", err))
			}
			return resultText(fmt.Sprintf("Wiki page '%s' deleted successfully", input.PagePath))
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
