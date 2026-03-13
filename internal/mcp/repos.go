package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageReposInput defines the input schema for the manage_repos tool.
type ManageReposInput struct {
	Action     string `json:"action" jsonschema:"required,description=Action to perform: 'list'\\, 'get'\\, 'list_branches'\\, 'get_file'\\, 'get_tree'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"description=Project name (required for most actions)"`
	RepoID     string `json:"repo_id,omitempty" jsonschema:"description=Repository name or ID (required for get\\, list_branches\\, get_file\\, get_tree)"`
	FilePath   string `json:"file_path,omitempty" jsonschema:"description=File path within the repo (for get_file\\, get_tree)"`
	Version    string `json:"version,omitempty" jsonschema:"description=Branch name or commit SHA (for get_file)"`
}

// ManageReposHandler returns the handler for the manage_repos tool.
func ManageReposHandler(c *devops.Client) func(context.Context, *sdkmcp.CallToolRequest, ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "list":
			repos, err := c.ListRepositories(input.ProjectKey)
			if err != nil {
				return resultError(fmt.Sprintf("listing repos: %v", err))
			}
			return resultJSON(repos)
		case "get":
			if input.RepoID == "" {
				return resultError("repo_id is required for 'get' action")
			}
			repo, err := c.GetRepository(input.ProjectKey, input.RepoID)
			if err != nil {
				return resultError(fmt.Sprintf("getting repo: %v", err))
			}
			return resultJSON(repo)
		case "list_branches":
			if input.RepoID == "" {
				return resultError("repo_id is required for 'list_branches' action")
			}
			branches, err := c.ListBranches(input.ProjectKey, input.RepoID)
			if err != nil {
				return resultError(fmt.Sprintf("listing branches: %v", err))
			}
			return resultJSON(branches)
		case "get_file":
			if input.RepoID == "" {
				return resultError("repo_id is required for 'get_file' action")
			}
			if input.FilePath == "" {
				return resultError("file_path is required for 'get_file' action")
			}
			item, err := c.GetFileContent(input.ProjectKey, input.RepoID, input.FilePath, input.Version)
			if err != nil {
				return resultError(fmt.Sprintf("getting file: %v", err))
			}
			return resultJSON(item)
		case "get_tree":
			if input.RepoID == "" {
				return resultError("repo_id is required for 'get_tree' action")
			}
			treePath := input.FilePath
			if treePath == "" {
				treePath = "/"
			}
			items, err := c.GetTree(input.ProjectKey, input.RepoID, treePath)
			if err != nil {
				return resultError(fmt.Sprintf("getting tree: %v", err))
			}
			return resultJSON(items)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
