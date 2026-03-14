package mcp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

func TestManageWorkItemsHandler_MyItems(t *testing.T) {
	t.Parallel()
	// my_items uses the 2-step WIQL pattern: POST wiql → POST workitemsbatch
	c := newTestClient(t, muxHandler(map[string]http.HandlerFunc{
		"/test-org/MyProject/_apis/wit/wiql": func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req map[string]interface{}
			_ = json.Unmarshal(body, &req)
			query := req["query"].(string)
			// Verify the WIQL query contains AssignedTo = @Me and project
			if !strings.Contains(query, "@Me") {
				t.Error("WIQL query missing @Me clause")
			}
			if !strings.Contains(query, "MyProject") {
				t.Error("WIQL query missing project name")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"queryType":"flat","queryResultType":"workItem","workItems":[{"id":10},{"id":20}]}`))
		},
		"/test-org/MyProject/_apis/wit/workitemsbatch": jsonHandler(`{"count":2,"value":[{"id":10,"rev":1,"fields":{"System.Title":"My Task","System.State":"Active"}},{"id":20,"rev":1,"fields":{"System.Title":"My Bug","System.State":"New"}}]}`),
	}))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "my_items",
		ProjectKey: "MyProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "My Task"`)
	assertResultSuccess(t, result, `"title": "My Bug"`)
}

func TestManageWorkItemsHandler_MyItems_WithTypeFilter(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, muxHandler(map[string]http.HandlerFunc{
		"/test-org/MyProject/_apis/wit/wiql": func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req map[string]interface{}
			_ = json.Unmarshal(body, &req)
			query := req["query"].(string)
			// Verify Bug type filter is appended
			if !strings.Contains(query, "[System.WorkItemType] = 'Bug'") {
				t.Error("WIQL query missing WorkItemType filter for Bug")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"queryType":"flat","queryResultType":"workItem","workItems":[{"id":20}]}`))
		},
		"/test-org/MyProject/_apis/wit/workitemsbatch": jsonHandler(`{"count":1,"value":[{"id":20,"rev":1,"fields":{"System.Title":"My Bug","System.WorkItemType":"Bug"}}]}`),
	}))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:       "my_items",
		ProjectKey:   "MyProject",
		WorkItemType: "Bug",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "My Bug"`)
}

func TestManageWorkItemsHandler_MyItems_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action: "my_items",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageWorkItemsHandler_IterationItems(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, muxHandler(map[string]http.HandlerFunc{
		"/test-org/Proj/_apis/work/teamsettings/iterations": jsonHandler(`{"workItemRelations":[{"target":{"id":5,"url":"https://test"}},{"target":{"id":6,"url":"https://test"}}]}`),
		"/test-org/Proj/_apis/wit/workitemsbatch":           jsonHandler(`{"count":2,"value":[{"id":5,"rev":1,"fields":{"System.Title":"Sprint Item 1"}},{"id":6,"rev":1,"fields":{"System.Title":"Sprint Item 2"}}]}`),
	}))
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:      "iteration_items",
		ProjectKey:  "Proj",
		IterationID: "iter-123",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Sprint Item 1")
	assertResultSuccess(t, result, "Sprint Item 2")
}

func TestManageWorkItemsHandler_IterationItems_MissingIterationID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "iteration_items",
		ProjectKey: "Proj",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "iteration_id is required")
}

func TestManageWorkItemsHandler_AddChildren(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify it's a PATCH to create work item
		body, _ := io.ReadAll(r.Body)
		var ops []map[string]interface{}
		_ = json.Unmarshal(body, &ops)

		// Find the relation op that has parent link
		foundParentLink := false
		for _, op := range ops {
			if op["path"] == "/relations/-" {
				val := op["value"].(map[string]interface{})
				if val["rel"] == "System.LinkTypes.Hierarchy-Reverse" {
					foundParentLink = true
				}
			}
		}
		if !foundParentLink {
			t.Error("expected parent link in created child items")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":101,"rev":1,"fields":{"System.Title":"Child Task"},"url":"https://test"}`))
	}))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:       "add_children",
		ProjectKey:   "Proj",
		ParentID:     50,
		WorkItemType: "Task",
		Titles:       []string{"Child Task"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "Child Task"`)
}

func TestManageWorkItemsHandler_AddChildren_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:       "add_children",
		ProjectKey:   "Proj",
		ParentID:     50,
		WorkItemType: "Task",
		Titles:       []string{"Child"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWorkItemsHandler_AddChildren_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	tests := []struct {
		name  string
		input ManageWorkItemsInput
		want  string
	}{
		{
			name:  "missing parent_id",
			input: ManageWorkItemsInput{Action: "add_children", WorkItemType: "Task", Titles: []string{"A"}},
			want:  "parent_id is required",
		},
		{
			name:  "missing work_item_type",
			input: ManageWorkItemsInput{Action: "add_children", ParentID: 50, Titles: []string{"A"}},
			want:  "work_item_type is required",
		},
		{
			name:  "missing titles",
			input: ManageWorkItemsInput{Action: "add_children", ParentID: 50, WorkItemType: "Task"},
			want:  "titles is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assertResultError(t, result, tt.want)
		})
	}
}

func TestManageWorkItemsHandler_Link(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var ops []map[string]interface{}
		_ = json.Unmarshal(body, &ops)

		// Verify the relation type and target URL are correct
		if len(ops) == 0 {
			t.Fatal("expected at least one patch operation")
		}
		val := ops[0]["value"].(map[string]interface{})
		rel := val["rel"].(string)
		if rel != "System.LinkTypes.Related" {
			t.Errorf("expected rel = System.LinkTypes.Related, got %q", rel)
		}
		urlStr := val["url"].(string)
		if !strings.Contains(urlStr, "/workitems/200") {
			t.Errorf("expected target URL to contain /workitems/200, got %q", urlStr)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":100,"rev":2,"fields":{"System.Title":"Linked Item"},"url":"https://test"}`))
	}))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "link",
		ProjectKey: "Proj",
		WorkItemID: 100,
		TargetID:   200,
		LinkType:   "System.LinkTypes.Related",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "Linked Item"`)
}

func TestManageWorkItemsHandler_Link_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	tests := []struct {
		name  string
		input ManageWorkItemsInput
		want  string
	}{
		{
			name:  "missing work_item_id",
			input: ManageWorkItemsInput{Action: "link", TargetID: 200, LinkType: "System.LinkTypes.Related"},
			want:  "work_item_id is required",
		},
		{
			name:  "missing target_id",
			input: ManageWorkItemsInput{Action: "link", WorkItemID: 100, LinkType: "System.LinkTypes.Related"},
			want:  "target_id is required",
		},
		{
			name:  "missing link_type",
			input: ManageWorkItemsInput{Action: "link", WorkItemID: 100, TargetID: 200},
			want:  "link_type is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assertResultError(t, result, tt.want)
		})
	}
}

func TestManageWorkItemsHandler_Unlink(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var ops []map[string]interface{}
		_ = json.Unmarshal(body, &ops)

		// Verify the operation is a remove at the correct index
		if len(ops) == 0 {
			t.Fatal("expected at least one patch operation")
		}
		if ops[0]["op"] != "remove" {
			t.Errorf("expected op = remove, got %q", ops[0]["op"])
		}
		if ops[0]["path"] != "/relations/2" {
			t.Errorf("expected path = /relations/2, got %q", ops[0]["path"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":100,"rev":3,"fields":{"System.Title":"Unlinked Item"},"url":"https://test"}`))
	}))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:        "unlink",
		ProjectKey:    "Proj",
		WorkItemID:    100,
		RelationIndex: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "Unlinked Item"`)
}

func TestManageWorkItemsHandler_Unlink_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:        "unlink",
		WorkItemID:    100,
		RelationIndex: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWorkItemsHandler_AddArtifactLink(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var ops []map[string]interface{}
		_ = json.Unmarshal(body, &ops)

		if len(ops) == 0 {
			t.Fatal("expected at least one patch operation")
		}
		val := ops[0]["value"].(map[string]interface{})
		if val["rel"] != "ArtifactLink" {
			t.Errorf("expected rel = ArtifactLink, got %q", val["rel"])
		}
		if val["url"] != "vstfs:///Git/Commit/proj%2Frepo%2Fabc123" {
			t.Errorf("unexpected artifact URI: %q", val["url"])
		}
		attrs := val["attributes"].(map[string]interface{})
		if attrs["name"] != "Fixed in Commit" {
			t.Errorf("expected name = Fixed in Commit, got %q", attrs["name"])
		}
		if attrs["comment"] != "linked commit" {
			t.Errorf("expected comment = linked commit, got %q", attrs["comment"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":100,"rev":4,"fields":{"System.Title":"With Artifact"},"url":"https://test"}`))
	}))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:      "add_artifact_link",
		ProjectKey:  "Proj",
		WorkItemID:  100,
		ArtifactURI: "vstfs:///Git/Commit/proj%2Frepo%2Fabc123",
		LinkType:    "Fixed in Commit",
		Comment:     "linked commit",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, `"title": "With Artifact"`)
}

func TestManageWorkItemsHandler_AddArtifactLink_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	tests := []struct {
		name  string
		input ManageWorkItemsInput
		want  string
	}{
		{
			name:  "missing work_item_id",
			input: ManageWorkItemsInput{Action: "add_artifact_link", ArtifactURI: "vstfs:///test", LinkType: "Build"},
			want:  "work_item_id is required",
		},
		{
			name:  "missing artifact_uri",
			input: ManageWorkItemsInput{Action: "add_artifact_link", WorkItemID: 100, LinkType: "Build"},
			want:  "artifact_uri is required",
		},
		{
			name:  "missing link_type",
			input: ManageWorkItemsInput{Action: "add_artifact_link", WorkItemID: 100, ArtifactURI: "vstfs:///test"},
			want:  "link_type is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assertResultError(t, result, tt.want)
		})
	}
}

func TestManageWorkItemsHandler_BatchUpdate(t *testing.T) {
	t.Parallel()
	// Each item gets its own PATCH request (sequential calls)
	c := newTestClient(t, jsonHandler(`{"id":10,"rev":2,"fields":{"System.Title":"Updated","System.State":"Active"},"url":"https://test"}`))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:      "batch_update",
		ProjectKey:  "Proj",
		WorkItemIDs: []int{10, 20},
		State:       "Active",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Updated")
}

func TestManageWorkItemsHandler_BatchUpdate_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:      "batch_update",
		WorkItemIDs: []int{1, 2},
		State:       "Closed",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageWorkItemsHandler_BatchUpdate_MissingIDs(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action: "batch_update",
		State:  "Active",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "work_item_ids is required")
}

func TestManageWorkItemsHandler_BatchUpdate_NoFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:      "batch_update",
		WorkItemIDs: []int{1, 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "at least one field to update is required")
}

func TestManageWorkItemsHandler_UpdateComment(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":5,"workItemId":42,"text":"updated comment text"}`))
	handler := ManageWorkItemsHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "update_comment",
		ProjectKey: "Proj",
		WorkItemID: 42,
		CommentID:  5,
		Comment:    "updated comment text",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "updated comment text")
}

func TestManageWorkItemsHandler_UpdateComment_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, true)

	tests := []struct {
		name  string
		input ManageWorkItemsInput
		want  string
	}{
		{
			name:  "missing work_item_id",
			input: ManageWorkItemsInput{Action: "update_comment", CommentID: 5, Comment: "text"},
			want:  "work_item_id is required",
		},
		{
			name:  "missing comment_id",
			input: ManageWorkItemsInput{Action: "update_comment", WorkItemID: 42, Comment: "text"},
			want:  "comment_id is required",
		},
		{
			name:  "missing comment",
			input: ManageWorkItemsInput{Action: "update_comment", WorkItemID: 42, CommentID: 5},
			want:  "comment is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assertResultError(t, result, tt.want)
		})
	}
}

func TestManageWorkItemsHandler_UpdateComment_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageWorkItemsHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageWorkItemsInput{
		Action:     "update_comment",
		WorkItemID: 42,
		CommentID:  5,
		Comment:    "text",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestBuildUpdateFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input ManageWorkItemsInput
		want  int // number of fields expected
	}{
		{"empty", ManageWorkItemsInput{}, 0},
		{"title only", ManageWorkItemsInput{Title: "test"}, 1},
		{"all fields", ManageWorkItemsInput{
			Title: "t", Description: "d", State: "s", AssignedTo: "a",
			AreaPath: "ap", IterationPath: "ip", Priority: 2, Tags: "tag1",
		}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildUpdateFields(tt.input)
			if len(got) != tt.want {
				t.Errorf("buildUpdateFields returned %d fields, want %d", len(got), tt.want)
			}
		})
	}
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
