package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageWorkItemsHandler_Get(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":42,"rev":1,"fields":{"System.Title":"Test Item","System.State":"Active"},"url":"https://test"}`))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "get",
		WorkItemID: 42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"id": 42`)
	assertResultSuccess(t, result, `"title": "Test Item"`)
}

func TestManageWorkItemsHandler_Get_MissingID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action: "get",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "work_item_id is required")
}

func TestManageWorkItemsHandler_BatchGet(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"id":1,"rev":1,"fields":{"System.Title":"A"}},{"id":2,"rev":1,"fields":{"System.Title":"B"}}]}`))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:      "batch_get",
		WorkItemIDs: []int{1, 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "A"`)
	assertResultSuccess(t, result, `"title": "B"`)
}

func TestManageWorkItemsHandler_BatchGet_MissingIDs(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action: "batch_get",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "work_item_ids is required")
}

func TestManageWorkItemsHandler_Create_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:       "create",
		ProjectKey:   "TestProject",
		WorkItemType: "Task",
		Title:        "New Task",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWorkItemsHandler_Create(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":99,"rev":1,"fields":{"System.Title":"New Task","System.State":"New"},"url":"https://test"}`))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:       "create",
		ProjectKey:   "TestProject",
		WorkItemType: "Task",
		Title:        "New Task",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"id": 99`)
}

func TestManageWorkItemsHandler_Create_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action: "create",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageWorkItemsHandler_Update_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "update",
		WorkItemID: 42,
		Title:      "Updated",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWorkItemsHandler_Delete_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "delete",
		WorkItemID: 42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWorkItemsHandler_ListComments(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"totalCount":1,"count":1,"comments":[{"id":1,"text":"hello","workItemId":42}]}`))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "list_comments",
		WorkItemID: 42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "hello")
}

func TestManageWorkItemsHandler_ListTypes(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":2,"value":[{"name":"Task"},{"name":"Bug"}]}`))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "list_types",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Task")
	assertResultSuccess(t, result, "Bug")
}

func TestManageWorkItemsHandler_GetHistory(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"count":1,"value":[{"id":1,"rev":1}]}`))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "get_history",
		WorkItemID: 42,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"rev"`)
}

func TestManageWorkItemsHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action: "invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}

func TestManageWorkItemsHandler_APIError(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, errorHandler(404))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "get",
		WorkItemID: 999,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "getting work item")
}

func TestFlattenFieldName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"System.Title", "title"},
		{"System.AssignedTo", "assigned_to"},
		{"Microsoft.VSTS.Common.Priority", "priority"},
		{"System.WorkItemType", "work_item_type"},
		{"CustomField", "custom_field"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := flattenFieldName(tt.input)
			if got != tt.want {
				t.Errorf("flattenFieldName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
