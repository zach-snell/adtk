package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManagePipelinesInput defines the input schema for the manage_pipelines tool.
type ManagePipelinesInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'list', 'get', 'list_runs', 'get_run', 'trigger', 'get_logs', 'get_log', 'get_build_changes', 'list_definitions', 'list_variable_groups', 'get_variable_group', 'list_environments'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	PipelineID int    `json:"pipeline_id,omitempty" jsonschema:"Pipeline ID (required for get, list_runs, trigger, get_logs, get_log)"`
	RunID      int    `json:"run_id,omitempty" jsonschema:"Run ID (required for get_run, get_logs, get_log)"`
	LogID      int    `json:"log_id,omitempty" jsonschema:"Log ID (required for get_log)"`
	BuildID    int    `json:"build_id,omitempty" jsonschema:"Build ID (required for get_build_changes)"`
	Branch     string `json:"branch,omitempty" jsonschema:"Branch name to run pipeline on (for trigger)"`
	Top        int    `json:"top,omitempty" jsonschema:"Max results to return"`
	GroupID    int    `json:"group_id,omitempty" jsonschema:"Variable group ID (required for get_variable_group)"`
}

// ManagePipelinesHandler returns the handler for the manage_pipelines tool.
func ManagePipelinesHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManagePipelinesInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManagePipelinesInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "list":
			pipelines, err := c.ListPipelines(input.ProjectKey, input.Top)
			if err != nil {
				return resultError(fmt.Sprintf("listing pipelines: %v", err))
			}
			return resultJSON(pipelines)
		case "get":
			if input.PipelineID == 0 {
				return resultError("pipeline_id is required for 'get' action")
			}
			pipeline, err := c.GetPipeline(input.ProjectKey, input.PipelineID)
			if err != nil {
				return resultError(fmt.Sprintf("getting pipeline: %v", err))
			}
			return resultJSON(pipeline)
		case "list_runs":
			if input.PipelineID == 0 {
				return resultError("pipeline_id is required for 'list_runs' action")
			}
			runs, err := c.ListPipelineRuns(input.ProjectKey, input.PipelineID, input.Top)
			if err != nil {
				return resultError(fmt.Sprintf("listing pipeline runs: %v", err))
			}
			return resultJSON(runs)
		case "get_run":
			if input.PipelineID == 0 || input.RunID == 0 {
				return resultError("pipeline_id and run_id are required for 'get_run' action")
			}
			run, err := c.GetPipelineRun(input.ProjectKey, input.PipelineID, input.RunID)
			if err != nil {
				return resultError(fmt.Sprintf("getting pipeline run: %v", err))
			}
			return resultJSON(run)
		case "trigger":
			if !enableWrites {
				return resultError("trigger action requires ADTK_ENABLE_WRITES=true")
			}
			if input.PipelineID == 0 {
				return resultError("pipeline_id is required for 'trigger' action")
			}
			run, err := c.TriggerPipeline(input.ProjectKey, input.PipelineID, input.Branch)
			if err != nil {
				return resultError(fmt.Sprintf("triggering pipeline: %v", err))
			}
			return resultJSON(run)
		case "get_logs":
			if input.PipelineID == 0 || input.RunID == 0 {
				return resultError("pipeline_id and run_id are required for 'get_logs' action")
			}
			data, err := c.GetPipelineLogs(input.ProjectKey, input.PipelineID, input.RunID)
			if err != nil {
				return resultError(fmt.Sprintf("getting pipeline logs: %v", err))
			}
			return resultText(string(data))
		case "get_log":
			if input.PipelineID == 0 || input.RunID == 0 || input.LogID == 0 {
				return resultError("pipeline_id, run_id, and log_id are required for 'get_log' action")
			}
			data, err := c.GetPipelineLog(input.ProjectKey, input.PipelineID, input.RunID, input.LogID)
			if err != nil {
				return resultError(fmt.Sprintf("getting pipeline log: %v", err))
			}
			return resultText(string(data))
		case "get_build_changes":
			if input.BuildID == 0 {
				return resultError("build_id is required for 'get_build_changes' action")
			}
			changes, err := c.GetBuildChanges(input.ProjectKey, input.BuildID)
			if err != nil {
				return resultError(fmt.Sprintf("getting build changes: %v", err))
			}
			return resultJSON(changes)
		case "list_definitions":
			defs, err := c.ListBuildDefinitions(input.ProjectKey)
			if err != nil {
				return resultError(fmt.Sprintf("listing build definitions: %v", err))
			}
			return resultJSON(defs)
		case "list_variable_groups":
			groups, err := c.ListVariableGroups(input.ProjectKey)
			if err != nil {
				return resultError(fmt.Sprintf("listing variable groups: %v", err))
			}
			return resultJSON(groups)
		case "get_variable_group":
			if input.GroupID == 0 {
				return resultError("group_id is required for 'get_variable_group' action")
			}
			group, err := c.GetVariableGroup(input.ProjectKey, input.GroupID)
			if err != nil {
				return resultError(fmt.Sprintf("getting variable group: %v", err))
			}
			return resultJSON(group)
		case "list_environments":
			envs, err := c.ListEnvironments(input.ProjectKey)
			if err != nil {
				return resultError(fmt.Sprintf("listing environments: %v", err))
			}
			return resultJSON(envs)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
