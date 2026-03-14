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
	Action      string `json:"action" jsonschema:"Action to perform: 'get', 'create', 'update', 'delete', 'add_comment', 'list_comments', 'get_links', 'list_types', 'get_history', 'batch_get', 'batch_update', 'add_children', 'link', 'unlink', 'add_artifact_link', 'my_items', 'iteration_items', 'update_comment'"`
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
	ParentID      int    `json:"parent_id,omitempty" jsonschema:"Parent work item ID to link (for create, add_children)"`
	Tags          string `json:"tags,omitempty" jsonschema:"Semicolon-separated tags (for create, update)"`

	// Comment
	Comment   string `json:"comment,omitempty" jsonschema:"Comment text in HTML (for add_comment, update_comment, add_artifact_link)"`
	CommentID int    `json:"comment_id,omitempty" jsonschema:"Comment ID (required for update_comment)"`

	// WIQL query
	Query  string   `json:"query,omitempty" jsonschema:"WIQL query string (for 'list' action via WIQL)"`
	Fields []string `json:"fields,omitempty" jsonschema:"Fields to return (for batch_get). Default: System.Title, System.State, System.AssignedTo"`
	Top    int      `json:"top,omitempty" jsonschema:"Max results to return (for WIQL queries)"`

	// Linking
	TargetID      int    `json:"target_id,omitempty" jsonschema:"Target work item ID (required for link)"`
	LinkType      string `json:"link_type,omitempty" jsonschema:"Link type name (for link, add_artifact_link), e.g. System.LinkTypes.Related, System.LinkTypes.Hierarchy-Forward"`
	RelationIndex int    `json:"relation_index,omitempty" jsonschema:"Relation index to remove (required for unlink)"`

	// Artifact link
	ArtifactURI string `json:"artifact_uri,omitempty" jsonschema:"Artifact URI for artifact links (required for add_artifact_link), e.g. vstfs:///Git/Commit/{projectId}%2F{repoId}%2F{commitId}"`

	// Children
	Titles []string `json:"titles,omitempty" jsonschema:"List of titles for child work items (required for add_children)"`

	// Iteration items
	Team        string `json:"team,omitempty" jsonschema:"Team name (optional, scopes iteration_items to a specific team)"`
	IterationID string `json:"iteration_id,omitempty" jsonschema:"Iteration ID (required for iteration_items)"`

	// My items
	IncludeCompleted bool `json:"include_completed,omitempty" jsonschema:"Include completed/closed work items (for my_items, default false)"`
}

// ManageWorkItemsHandler returns the handler for the manage_work_items tool.
func ManageWorkItemsHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "get":
			return actionGetWorkItem(c, input)
		case "batch_get":
			return actionBatchGetWorkItems(c, input)
		case "create":
			return actionCreateWorkItem(c, enableWrites, input)
		case "update":
			return actionUpdateWorkItem(c, enableWrites, input)
		case "delete":
			return actionDeleteWorkItem(c, enableWrites, input)
		case "add_comment":
			return actionAddWorkItemComment(c, enableWrites, input)
		case "list_comments":
			return actionListWorkItemComments(c, input)
		case "get_links":
			return actionGetWorkItemLinks(c, input)
		case "list_types":
			return actionListWorkItemTypes(c, input)
		case "get_history":
			return actionGetWorkItemHistory(c, input)
		case "batch_update":
			return actionBatchUpdateWorkItems(c, enableWrites, input)
		case "add_children":
			return actionAddChildWorkItems(c, enableWrites, input)
		case "link":
			return actionLinkWorkItems(c, enableWrites, input)
		case "unlink":
			return actionUnlinkWorkItem(c, enableWrites, input)
		case "add_artifact_link":
			return actionAddArtifactLink(c, enableWrites, input)
		case "my_items":
			return actionMyWorkItems(c, input)
		case "iteration_items":
			return actionIterationWorkItems(c, input)
		case "update_comment":
			return actionUpdateWorkItemComment(c, enableWrites, input)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func actionGetWorkItem(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'get' action")
	}
	return handleGetWorkItem(c, input.ProjectKey, input.WorkItemID)
}

func actionBatchGetWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if len(input.WorkItemIDs) == 0 {
		return resultError("work_item_ids is required for 'batch_get' action")
	}
	return handleBatchGetWorkItems(c, input.ProjectKey, input.WorkItemIDs, input.Fields)
}

func actionCreateWorkItem(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
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
}

func actionUpdateWorkItem(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("update action requires ADTK_ENABLE_WRITES=true")
	}
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'update' action")
	}
	return handleUpdateWorkItem(c, input)
}

func actionDeleteWorkItem(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("delete action requires ADTK_ENABLE_WRITES=true")
	}
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'delete' action")
	}
	return handleDeleteWorkItem(c, input.ProjectKey, input.WorkItemID)
}

func actionAddWorkItemComment(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
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
}

func actionListWorkItemComments(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'list_comments' action")
	}
	return handleListWorkItemComments(c, input.ProjectKey, input.WorkItemID)
}

func actionGetWorkItemLinks(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'get_links' action")
	}
	return handleGetWorkItemLinks(c, input.ProjectKey, input.WorkItemID)
}

func actionListWorkItemTypes(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.ProjectKey == "" {
		return resultError("project_key is required for 'list_types' action")
	}
	return handleListWorkItemTypes(c, input.ProjectKey)
}

func actionGetWorkItemHistory(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'get_history' action")
	}
	return handleGetWorkItemHistory(c, input.ProjectKey, input.WorkItemID)
}

func actionBatchUpdateWorkItems(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("batch_update action requires ADTK_ENABLE_WRITES=true")
	}
	if len(input.WorkItemIDs) == 0 {
		return resultError("work_item_ids is required for 'batch_update' action")
	}
	return handleBatchUpdateWorkItems(c, input)
}

func actionAddChildWorkItems(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("add_children action requires ADTK_ENABLE_WRITES=true")
	}
	if input.ParentID == 0 {
		return resultError("parent_id is required for 'add_children' action")
	}
	if input.WorkItemType == "" {
		return resultError("work_item_type is required for 'add_children' action")
	}
	if len(input.Titles) == 0 {
		return resultError("titles is required for 'add_children' action")
	}
	return handleAddChildWorkItems(c, input)
}

func actionLinkWorkItems(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("link action requires ADTK_ENABLE_WRITES=true")
	}
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'link' action")
	}
	if input.TargetID == 0 {
		return resultError("target_id is required for 'link' action")
	}
	if input.LinkType == "" {
		return resultError("link_type is required for 'link' action")
	}
	return handleLinkWorkItems(c, input)
}

func actionUnlinkWorkItem(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("unlink action requires ADTK_ENABLE_WRITES=true")
	}
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'unlink' action")
	}
	return handleUnlinkWorkItem(c, input)
}

func actionAddArtifactLink(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("add_artifact_link action requires ADTK_ENABLE_WRITES=true")
	}
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'add_artifact_link' action")
	}
	if input.ArtifactURI == "" {
		return resultError("artifact_uri is required for 'add_artifact_link' action")
	}
	if input.LinkType == "" {
		return resultError("link_type is required for 'add_artifact_link' action")
	}
	return handleAddArtifactLink(c, input)
}

func actionMyWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.ProjectKey == "" {
		return resultError("project_key is required for 'my_items' action")
	}
	return handleMyWorkItems(c, input)
}

func actionIterationWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if input.ProjectKey == "" {
		return resultError("project_key is required for 'iteration_items' action")
	}
	if input.IterationID == "" {
		return resultError("iteration_id is required for 'iteration_items' action")
	}
	return handleIterationWorkItems(c, input)
}

func actionUpdateWorkItemComment(c *devops.Client, enableWrites bool, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("update_comment action requires ADTK_ENABLE_WRITES=true")
	}
	if input.WorkItemID == 0 {
		return resultError("work_item_id is required for 'update_comment' action")
	}
	if input.CommentID == 0 {
		return resultError("comment_id is required for 'update_comment' action")
	}
	if input.Comment == "" {
		return resultError("comment is required for 'update_comment' action")
	}
	return handleUpdateWorkItemComment(c, input)
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

func handleBatchUpdateWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	// Build a batch update: apply the same field changes to all specified IDs
	fields := buildUpdateFields(input)
	if len(fields) == 0 {
		return resultError("at least one field to update is required for batch_update")
	}

	ops := devops.BuildJSONPatchOps(fields)
	updates := make([]devops.BatchWorkItemUpdate, len(input.WorkItemIDs))
	for i, id := range input.WorkItemIDs {
		updates[i] = devops.BatchWorkItemUpdate{ID: id, Ops: ops}
	}

	items, err := c.UpdateWorkItemsBatch(input.ProjectKey, updates)
	if err != nil {
		return resultError(fmt.Sprintf("batch updating work items: %v", err))
	}

	flat := make([]map[string]interface{}, len(items))
	for i, wi := range items {
		flat[i] = flattenWorkItem(&wi)
	}
	return resultJSON(flat)
}

func handleAddChildWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	items, err := c.AddChildWorkItems(input.ProjectKey, input.ParentID, input.WorkItemType, input.Titles)
	if err != nil {
		return resultError(fmt.Sprintf("adding child work items: %v", err))
	}

	flat := make([]map[string]interface{}, len(items))
	for i, wi := range items {
		flat[i] = flattenWorkItem(&wi)
	}
	return resultJSON(flat)
}

func handleLinkWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	wi, err := c.LinkWorkItems(input.ProjectKey, input.WorkItemID, input.TargetID, input.LinkType)
	if err != nil {
		return resultError(fmt.Sprintf("linking work items %d -> %d: %v", input.WorkItemID, input.TargetID, err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleUnlinkWorkItem(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	wi, err := c.UnlinkWorkItem(input.ProjectKey, input.WorkItemID, input.RelationIndex)
	if err != nil {
		return resultError(fmt.Sprintf("unlinking relation %d from work item %d: %v", input.RelationIndex, input.WorkItemID, err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleAddArtifactLink(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	wi, err := c.AddArtifactLink(input.ProjectKey, input.WorkItemID, input.ArtifactURI, input.LinkType, input.Comment)
	if err != nil {
		return resultError(fmt.Sprintf("adding artifact link to work item %d: %v", input.WorkItemID, err))
	}
	return resultJSON(flattenWorkItem(wi))
}

func handleMyWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	top := input.Top
	if top == 0 {
		top = 25
	}
	items, err := c.GetMyWorkItems(input.ProjectKey, input.WorkItemType, input.IncludeCompleted, top)
	if err != nil {
		return resultError(fmt.Sprintf("getting my work items: %v", err))
	}

	flat := make([]map[string]interface{}, len(items))
	for i, wi := range items {
		flat[i] = flattenWorkItem(&wi)
	}
	return resultJSON(flat)
}

func handleIterationWorkItems(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	items, err := c.GetWorkItemsForIteration(input.ProjectKey, input.Team, input.IterationID)
	if err != nil {
		return resultError(fmt.Sprintf("getting iteration work items: %v", err))
	}

	flat := make([]map[string]interface{}, len(items))
	for i, wi := range items {
		flat[i] = flattenWorkItem(&wi)
	}
	return resultJSON(flat)
}

func handleUpdateWorkItemComment(c *devops.Client, input ManageWorkItemsInput) (*sdkmcp.CallToolResult, any, error) {
	comment, err := c.UpdateWorkItemComment(input.ProjectKey, input.WorkItemID, input.CommentID, input.Comment)
	if err != nil {
		return resultError(fmt.Sprintf("updating comment %d on work item %d: %v", input.CommentID, input.WorkItemID, err))
	}
	return resultJSON(comment)
}

// buildUpdateFields extracts field values from input for update operations.
func buildUpdateFields(input ManageWorkItemsInput) map[string]interface{} {
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
	return fields
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
