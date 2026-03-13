package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageProjectsInput defines the input schema for the manage_projects tool.
type ManageProjectsInput struct {
	Action      string `json:"action" jsonschema:"Action to perform: 'list', 'get', 'list_teams', 'get_team', 'create'"`
	ProjectKey  string `json:"project_key,omitempty" jsonschema:"Project name or ID (required for get, list_teams)"`
	TeamID      string `json:"team_id,omitempty" jsonschema:"Team name or ID (required for get_team)"`
	Name        string `json:"name,omitempty" jsonschema:"Project name (required for create)"`
	Description string `json:"description,omitempty" jsonschema:"Project description (for create)"`
}

// ManageProjectsHandler returns the handler for the manage_projects tool.
func ManageProjectsHandler(c *devops.Client, enableWrites bool) func(context.Context, *mcp.CallToolRequest, ManageProjectsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ManageProjectsInput) (*mcp.CallToolResult, any, error) {
		switch input.Action {
		case "list":
			return handleListProjects(c)
		case "get":
			if input.ProjectKey == "" {
				return resultError("project_key is required for 'get' action")
			}
			return handleGetProject(c, input.ProjectKey)
		case "list_teams":
			if input.ProjectKey == "" {
				return resultError("project_key is required for 'list_teams' action")
			}
			return handleListTeams(c, input.ProjectKey)
		case "get_team":
			if input.ProjectKey == "" {
				return resultError("project_key is required for 'get_team' action")
			}
			if input.TeamID == "" {
				return resultError("team_id is required for 'get_team' action")
			}
			return handleGetTeam(c, input.ProjectKey, input.TeamID)
		case "create":
			if !enableWrites {
				return resultError("create action requires ADTK_ENABLE_WRITES=true")
			}
			if input.Name == "" {
				return resultError("name is required for 'create' action")
			}
			return handleCreateProject(c, input.Name, input.Description)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleListProjects(c *devops.Client) (*mcp.CallToolResult, any, error) {
	result, err := devops.GetJSON[devops.ProjectList](c, "", "/projects", nil)
	if err != nil {
		return resultError(fmt.Sprintf("listing projects: %v", err))
	}
	return resultJSON(result)
}

func handleGetProject(c *devops.Client, projectKey string) (*mcp.CallToolResult, any, error) {
	path := fmt.Sprintf("/projects/%s", projectKey)
	result, err := devops.GetJSON[devops.Project](c, "", path, nil)
	if err != nil {
		return resultError(fmt.Sprintf("getting project: %v", err))
	}
	return resultJSON(result)
}

func handleListTeams(c *devops.Client, projectKey string) (*mcp.CallToolResult, any, error) {
	path := fmt.Sprintf("/projects/%s/teams", projectKey)
	result, err := devops.GetJSON[devops.TeamList](c, "", path, nil)
	if err != nil {
		return resultError(fmt.Sprintf("listing teams: %v", err))
	}
	return resultJSON(result)
}

func handleGetTeam(c *devops.Client, projectKey, teamID string) (*mcp.CallToolResult, any, error) {
	path := fmt.Sprintf("/projects/%s/teams/%s", projectKey, teamID)
	result, err := devops.GetJSON[devops.Team](c, "", path, nil)
	if err != nil {
		return resultError(fmt.Sprintf("getting team: %v", err))
	}
	return resultJSON(result)
}

func handleCreateProject(c *devops.Client, name, description string) (*mcp.CallToolResult, any, error) {
	body := map[string]interface{}{
		"name":        name,
		"description": description,
		"capabilities": map[string]interface{}{
			"versioncontrol": map[string]string{
				"sourceControlType": "Git",
			},
			"processTemplate": map[string]string{
				"templateTypeId": "6b724908-ef14-45cf-84f8-768b5384da45", // Agile
			},
		},
	}
	data, err := c.Post("", "/projects", body)
	if err != nil {
		return resultError(fmt.Sprintf("creating project: %v", err))
	}
	return resultText(fmt.Sprintf("Project creation initiated. Note: project creation is async in Azure DevOps.\n%s", string(data)))
}
