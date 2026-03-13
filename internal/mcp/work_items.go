package mcp

import (
	"context"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageWorkItemsInput defines the input schema for the manage_work_items tool.
type ManageWorkItemsInput struct {
	Action      string `json:"action" jsonschema:"Action to perform: 'get', 'create', 'update', 'delete', 'add_comment', 'list_comments', 'get_links', 'list_types', 'get_history', 'batch_get'"`
	ProjectKey  string `json:"project_key,omitempty" jsonschema:"Project name (required for most actions)"`
	WorkItemID  int    `json:"work_item_id,omitempty" jsonschema:"Work item ID (required for get, update, delete, add_comment, list_comments, get_links, get_history)"`
	WorkItemIDs []int  `json:"work_item_ids,omitempty" jsonschema:"Work item IDs (for batch_get, max 200)"`

	// Create/update fields
	WorkItemType  string `json:"work_item_type,omitempty" jsonschema:"Work item type: Task, Bug, User Story, Epic, Feature, Issue (required for create)"`
	Title         string `json:"title,omitempty" jsonschema:"Work item title (required for create)"`
	Description   string `json:"description,omitempty" jsonschema:"Work item description in HTML (for create, update)"`
	State         string `json:"state,omitempty" jsonschema:"Work item state: New, Active, Closed, etc. (for update)"`
	AssignedTo    string `json:"assigned_to,omitempty" jsonschema:"Assignee display name or email (for create, update)"`
	AreaPath      string `json:"area_path,omitempty" jsonschema:"Area path e.g. Project\\Team (for create, update)"`
	IterationPath string `json:"iteration_path,omitempty" jsonschema:"Iteration path e.g. Project\\Sprint 1 (for create, update)"`
	Priority      int    `json:"priority,omitempty" jsonschema:"Priority: 1=Critical, 2=High, 3=Medium, 4=Low (for create, update)"`
	ParentID      int    `json:"parent_id,omitempty" jsonschema:"Parent work item ID to link (for create)"`
	Tags          string `json:"tags,omitempty" jsonschema:"Semicolon-separated tags (for create, update)"`

	// Comment
	Comment string `json:"comment,omitempty" jsonschema:"Comment text in HTML (for add_comment)"`

	// WIQL query
	Query  string   `json:"query,omitempty" jsonschema:"WIQL query string (for 'list' action via WIQL)"`
	Fields []string `json:"fields,omitempty" jsonschema:"Fields to return (for batch_get). Default: System.Title, System.State, System.AssignedTo"`
	Top    int      `json:"top,omitempty" jsonschema:"Max results to return (for WIQL queries)"`
}

// ManageWorkItemsHandler returns the handler for the manage_work_items tool.
func ManageWorkItemsHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "get":
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'get' action")
			}
			return handleGetWorkItem(c, input.ProjectKey, input.WorkItemID)

		case "batch_get":
			if len(input.WorkItemIDs) == 0 {
				return resultError("work_item_ids is required for 'batch_get' action")
			}
			return handleBatchGetWorkItems(c, input.ProjectKey, input.WorkItemIDs, input.Fields)

		case "create":
			if !enableWrites {
				return resultError("create action requires ADTK_ENABLE_WRITES=true")
			}
			if input.ProjectKey == "" {
				return resultError("project_key is required for 'create' action")
			}
			if input.WorkItemType == "" {
				return resultError("work_item_type is required for 'create' action")
			}
			if input.Title == "" {
				return resultError("title is required for 'create' action")
			}
			return handleCreateWorkItem(c, input)

		case "update":
			if !enableWrites {
				return resultError("update action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'update' action")
			}
			return handleUpdateWorkItem(c, input)

		case "delete":
			if !enableWrites {
				return resultError("delete action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'delete' action")
			}
			return handleDeleteWorkItem(c, input.ProjectKey, input.WorkItemID)

		case "add_comment":
			if !enableWrites {
				return resultError("add_comment action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'add_comment' action")
			}
			if input.Comment == "" {
				return resultError("comment is required for 'add_comment' action")
			}
			return handleAddWorkItemComment(c, input.ProjectKey, input.WorkItemID, input.Comment)

		case "list_comments":
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'list_comments' action")
			}
			return handleListWorkItemComments(c, input.ProjectKey, input.WorkItemID)

		case "get_links":
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'get_links' action")
			}
			return handleGetWorkItemLinks(c, input.ProjectKey, input.WorkItemID)

		case "list_types":
			if input.ProjectKey == "" {
				return resultError("project_key is required for 'list_types' action")
			}
			return handleListWorkItemTypes(c, input.ProjectKey)

		case "get_history":
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'get_history' action")
			}
			return handleGetWorkItemHistory(c, input.ProjectKey, input.WorkItemID)

		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleGetWorkItem(c *devops.Client, project string, id int) (*sdkmcp.CallToolResult, any, error) {
	wi, err := c.GetWorkItem(project, id, "All")
	if err != nil {
		return resultError(fmt.Sprintf("getting work item %d: %v", id, err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleBatchGetWorkItems(c *devops.Client, project string, ids []int, fields []string) (*sdkmcp.CallToolResult, any, error) {
	if len(fields) == 0 {
		fields = []string{
			"System.Id",
			"System.Title",
			"System.State",
			"System.AssignedTo",
			"System.WorkItemType",
			"System.CreatedDate",
			"System.ChangedDate",
		}
	}
	items, err := c.GetWorkItemsBatch(project, ids, fields)
	if err != nil {
		return resultError(fmt.Sprintf("batch getting work items: %v", err))
	}

	flat := make([]map[string]interface{}, len(items))
	for i, wi := range items {
		flat[i] = flattenWorkItem(&wi)
	}
	return resultJSON(flat)
}

func handleCreateWorkItem(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	fields := map[string]interface{}{
		"Title": input.Title,
	}
	if input.Description != "" {
		fields["Description"] = input.Description
	}
	if input.State != "" {
		fields["State"] = input.State
	}
	if input.AssignedTo != "" {
		fields["AssignedTo"] = input.AssignedTo
	}
	if input.AreaPath != "" {
		fields["AreaPath"] = input.AreaPath
	}
	if input.IterationPath != "" {
		fields["IterationPath"] = input.IterationPath
	}
	if input.Priority > 0 {
		fields["Microsoft.VSTS.Common.Priority"] = input.Priority
	}
	if input.Tags != "" {
		fields["Tags"] = input.Tags
	}

	ops := devops.BuildJSONPatchOps(fields)

	// Add parent link if specified
	if input.ParentID > 0 {
		parentURL := fmt.Sprintf("https://%s/%s/_apis/wit/workitems/%d",
			devops.HostMain, c.Organization(), input.ParentID)
		ops = append(ops, devops.JSONPatchOp{
			Op:   "add",
			Path: "/relations/-",
			Value: map[string]interface{}{
				"rel": "System.LinkTypes.Hierarchy-Reverse",
				"url": parentURL,
			},
		})
	}

	wi, err := c.CreateWorkItem(input.ProjectKey, input.WorkItemType, ops)
	if err != nil {
		return resultError(fmt.Sprintf("creating work item: %v", err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleUpdateWorkItem(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	fields := make(map[string]interface{})
	if input.Title != "" {
		fields["Title"] = input.Title
	}
	if input.Description != "" {
		fields["Description"] = input.Description
	}
	if input.State != "" {
		fields["State"] = input.State
	}
	if input.AssignedTo != "" {
		fields["AssignedTo"] = input.AssignedTo
	}
	if input.AreaPath != "" {
		fields["AreaPath"] = input.AreaPath
	}
	if input.IterationPath != "" {
		fields["IterationPath"] = input.IterationPath
	}
	if input.Priority > 0 {
		fields["Microsoft.VSTS.Common.Priority"] = input.Priority
	}
	if input.Tags != "" {
		fields["Tags"] = input.Tags
	}

	if len(fields) == 0 {
		return resultError("at least one field to update is required")
	}

	ops := devops.BuildJSONPatchOps(fields)
	wi, err := c.UpdateWorkItem(input.ProjectKey, input.WorkItemID, ops)
	if err != nil {
		return resultError(fmt.Sprintf("updating work item %d: %v", input.WorkItemID, err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleDeleteWorkItem(c *devops.Client, project string, id int) (*sdkmcp.CallToolResult, any, error) {
	if err := c.DeleteWorkItem(project, id); err != nil {
		return resultError(fmt.Sprintf("deleting work item %d: %v", id, err))
	}
	return resultText(fmt.Sprintf("Work item %d deleted successfully", id))
}

func handleAddWorkItemComment(c *devops.Client, project string, id int, text string) (*sdkmcp.CallToolResult, any, error) {
	comment, err := c.AddWorkItemComment(project, id, text)
	if err != nil {
		return resultError(fmt.Sprintf("adding comment to work item %d: %v", id, err))
	}
	return resultJSON(comment)
}

func handleListWorkItemComments(c *devops.Client, project string, id int) (*sdkmcp.CallToolResult, any, error) {
	comments, err := c.GetWorkItemComments(project, id)
	if err != nil {
		return resultError(fmt.Sprintf("listing comments for work item %d: %v", id, err))
	}
	return resultJSON(comments)
}

func handleGetWorkItemLinks(c *devops.Client, project string, id int) (*sdkmcp.CallToolResult, any, error) {
	wi, err := c.GetWorkItem(project, id, "Relations")
	if err != nil {
		return resultError(fmt.Sprintf("getting work item %d relations: %v", id, err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleListWorkItemTypes(c *devops.Client, project string) (*sdkmcp.CallToolResult, any, error) {
	types, err := c.GetWorkItemTypes(project)
	if err != nil {
		return resultError(fmt.Sprintf("listing work item types: %v", err))
	}
	return resultJSON(types)
}

func handleGetWorkItemHistory(c *devops.Client, project string, id int) (*sdkmcp.CallToolResult, any, error) {
	updates, err := c.GetWorkItemUpdates(project, id)
	if err != nil {
		return resultError(fmt.Sprintf("getting history for work item %d: %v", id, err))
	}
	return resultJSON(updates)
}

// flattenWorkItem converts a WorkItem to a flat map, stripping _links and
// converting System.* field names to readable snake_case.
func flattenWorkItem(wi *devops.WorkItem) map[string]interface{} {
	flat := map[string]interface{}{
		"id":  wi.ID,
		"rev": wi.Rev,
		"url": wi.URL,
	}

	for key, val := range wi.Fields {
		flatKey := flattenFieldName(key)
		flat[flatKey] = val
	}

	return flat
}

// flattenFieldName converts "System.Title" → "title", "Microsoft.VSTS.Common.Priority" → "priority".
func flattenFieldName(name string) string {
	// Strip common prefixes
	name = strings.TrimPrefix(name, "System.")
	name = strings.TrimPrefix(name, "Microsoft.VSTS.Common.")
	name = strings.TrimPrefix(name, "Microsoft.VSTS.Scheduling.")
	name = strings.TrimPrefix(name, "Microsoft.VSTS.TCM.")

	// Convert CamelCase to snake_case
	var result []byte
	for i, ch := range name {
		if ch >= 'A' && ch <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(ch-'A'+'a'))
		} else {
			result = append(result, byte(ch))
		}
	}
	return string(result)
}
