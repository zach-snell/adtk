package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetWorkItem retrieves a single work item by ID with optional field expansion.
func (c *Client) GetWorkItem(project string, id int, expand string) (*WorkItem, error) {
	query := url.Values{}
	if expand != "" {
		query.Set("$expand", expand)
	}
	path := fmt.Sprintf("/wit/workitems/%d", id)
	return GetJSON[WorkItem](c, project, path, query)
}

// GetWorkItemsBatch retrieves multiple work items by IDs (max 200).
// This is the second step of the WIQL 2-step pattern.
func (c *Client) GetWorkItemsBatch(project string, ids []int, fields []string) ([]WorkItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if len(ids) > 200 {
		ids = ids[:200]
	}

	body := map[string]interface{}{
		"ids": ids,
	}
	if len(fields) > 0 {
		body["fields"] = fields
	}

	path := "/wit/workitemsbatch"
	data, err := c.Post(project, path, body)
	if err != nil {
		return nil, fmt.Errorf("batch get work items: %w", err)
	}

	var result WorkItemList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling batch result: %w", err)
	}

	return result.Value, nil
}

// CreateWorkItem creates a new work item of the given type.
// Uses JSON Patch format: [{"op":"add","path":"/fields/System.Title","value":"..."}]
func (c *Client) CreateWorkItem(project, workItemType string, ops []JSONPatchOp) (*WorkItem, error) {
	path := fmt.Sprintf("/wit/workitems/$%s", workItemType)
	data, err := c.PatchJSONPatch(project, path, ops)
	if err != nil {
		return nil, fmt.Errorf("creating work item: %w", err)
	}

	var result WorkItem
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling work item: %w", err)
	}

	return &result, nil
}

// UpdateWorkItem updates an existing work item with JSON Patch operations.
func (c *Client) UpdateWorkItem(project string, id int, ops []JSONPatchOp) (*WorkItem, error) {
	path := fmt.Sprintf("/wit/workitems/%d", id)
	data, err := c.PatchJSONPatch(project, path, ops)
	if err != nil {
		return nil, fmt.Errorf("updating work item: %w", err)
	}

	var result WorkItem
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling work item: %w", err)
	}

	return &result, nil
}

// DeleteWorkItem deletes a work item by ID.
func (c *Client) DeleteWorkItem(project string, id int) error {
	path := fmt.Sprintf("/wit/workitems/%d", id)
	return c.Delete(project, path)
}

// RunWIQL executes a WIQL query and returns work item IDs.
func (c *Client) RunWIQL(project, query string, top int) (*WIQLResult, error) {
	body := map[string]interface{}{
		"query": query,
	}

	path := "/wit/wiql"
	q := url.Values{}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	requestURL := c.buildURL(HostMain, project, path, q)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling WIQL body: %w", err)
	}

	resp, doErr := c.do("POST", requestURL, bodyBytes, "application/json")
	if doErr != nil {
		return nil, fmt.Errorf("running WIQL: %w", doErr)
	}
	defer resp.Body.Close()

	var buf []byte
	buf, err = readBody(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, buf)
	}

	var result WIQLResult
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling WIQL result: %w", err)
	}

	return &result, nil
}

// WIQLAndFetch executes a WIQL query and fetches the full work items in one call.
// This implements the 2-step WIQL pattern: query IDs → batch fetch.
func (c *Client) WIQLAndFetch(project, query string, fields []string, top int) ([]WorkItem, error) {
	result, err := c.RunWIQL(project, query, top)
	if err != nil {
		return nil, err
	}

	if len(result.WorkItems) == 0 {
		return nil, nil
	}

	ids := make([]int, len(result.WorkItems))
	for i, wi := range result.WorkItems {
		ids[i] = wi.ID
	}

	return c.GetWorkItemsBatch(project, ids, fields)
}

// GetWorkItemComments retrieves comments for a work item.
func (c *Client) GetWorkItemComments(project string, id int) (*WorkItemCommentList, error) {
	path := fmt.Sprintf("/wit/workitems/%d/comments", id)
	return GetJSON[WorkItemCommentList](c, project, path, nil)
}

// AddWorkItemComment adds a comment to a work item.
func (c *Client) AddWorkItemComment(project string, id int, text string) (*WorkItemComment, error) {
	path := fmt.Sprintf("/wit/workitems/%d/comments", id)
	body := map[string]string{"text": text}
	data, err := c.Post(project, path, body)
	if err != nil {
		return nil, fmt.Errorf("adding comment: %w", err)
	}

	var result WorkItemComment
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling comment: %w", err)
	}

	return &result, nil
}

// GetWorkItemTypes lists the work item types for a project.
func (c *Client) GetWorkItemTypes(project string) ([]WorkItemType, error) {
	result, err := GetJSON[WorkItemTypeList](c, project, "/wit/workitemtypes", nil)
	if err != nil {
		return nil, fmt.Errorf("getting work item types: %w", err)
	}
	return result.Value, nil
}

// GetWorkItemUpdates retrieves the update history for a work item.
func (c *Client) GetWorkItemUpdates(project string, id int) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/wit/workitems/%d/updates", id)
	data, err := c.Get(project, path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting work item history: %w", err)
	}

	var result struct {
		Count int                      `json:"count"`
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling history: %w", err)
	}

	return result.Value, nil
}

// BuildJSONPatchOps constructs JSON Patch operations from a map of field names to values.
// Field names are automatically prefixed with "/fields/System." if they don't start with "/".
func BuildJSONPatchOps(fields map[string]interface{}) []JSONPatchOp {
	ops := make([]JSONPatchOp, 0, len(fields))
	for field, value := range fields {
		path := field
		if !strings.HasPrefix(path, "/") {
			// Auto-prefix with /fields/ and add System. if not already qualified
			if !strings.Contains(path, ".") {
				path = "/fields/System." + path
			} else {
				path = "/fields/" + path
			}
		}
		ops = append(ops, JSONPatchOp{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return ops
}
