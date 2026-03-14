package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManageTestPlansHandler_ListPlans(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"id":1,"name":"Release 2.0 Plan","state":"Active","iteration":"Project\\Sprint 5"},{"id":2,"name":"Regression Suite","state":"Active"}]}`))
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "list_plans",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Release 2.0 Plan")
	assertResultSuccess(t, result, "Regression Suite")
}

func TestManageTestPlansHandler_CreatePlan(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":3,"name":"New Test Plan","state":"Active","iteration":"Project\\Sprint 6"}`))
	handler := ManageTestPlansHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "create_plan",
		ProjectKey: "TestProject",
		Name:       "New Test Plan",
		Iteration:  "Project\\Sprint 6",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "New Test Plan")
}

func TestManageTestPlansHandler_CreatePlan_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "create_plan",
		ProjectKey: "TestProject",
		Name:       "Plan",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageTestPlansHandler_CreatePlan_MissingName(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "create_plan",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "name is required")
}

func TestManageTestPlansHandler_ListSuites(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"id":100,"name":"Default Suite","suiteType":"staticTestSuite","parentSuite":{"id":99}},{"id":101,"name":"Login Tests","suiteType":"staticTestSuite"}]}`))
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "list_suites",
		ProjectKey: "TestProject",
		PlanID:     1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Default Suite")
	assertResultSuccess(t, result, "Login Tests")
}

func TestManageTestPlansHandler_ListSuites_MissingPlanID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "list_suites",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "plan_id is required")
}

func TestManageTestPlansHandler_CreateSuite(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"id":102,"name":"Payment Tests","suiteType":"staticTestSuite","parentSuite":{"id":100}}`))
	handler := ManageTestPlansHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:        "create_suite",
		ProjectKey:    "TestProject",
		PlanID:        1,
		ParentSuiteID: 100,
		Name:          "Payment Tests",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Payment Tests")
}

func TestManageTestPlansHandler_CreateSuite_WritesDisabled(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "create_suite",
		ProjectKey: "TestProject",
		PlanID:     1,
		Name:       "Suite",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "ADTK_ENABLE_WRITES=true")
}

func TestManageTestPlansHandler_CreateSuite_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, true)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "create_suite",
		ProjectKey: "TestProject",
		PlanID:     1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "name are required")
}

func TestManageTestPlansHandler_ListCases(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"workItem":{"id":500,"name":"Login succeeds with valid creds"}},{"workItem":{"id":501,"name":"Login fails with bad password"}}]}`))
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "list_cases",
		ProjectKey: "TestProject",
		PlanID:     1,
		SuiteID:    100,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Login succeeds with valid creds")
	assertResultSuccess(t, result, "Login fails with bad password")
}

func TestManageTestPlansHandler_ListCases_MissingFields(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "list_cases",
		ProjectKey: "TestProject",
		PlanID:     1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "suite_id are required")
}

func TestManageTestPlansHandler_GetTestResults(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, jsonHandler(`{"value":[{"id":1,"name":"Login Test Run","totalTests":10,"passedTests":8,"failedTests":2,"state":"Completed","buildConfiguration":{"buildId":500}}]}`))
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "get_test_results",
		ProjectKey: "TestProject",
		BuildID:    500,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultSuccess(t, result, "Login Test Run")
	assertResultSuccess(t, result, "Completed")
}

func TestManageTestPlansHandler_GetTestResults_MissingBuildID(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "get_test_results",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "build_id is required")
}

func TestManageTestPlansHandler_MissingProject(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action: "list_plans",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "project_key is required")
}

func TestManageTestPlansHandler_UnknownAction(t *testing.T) {
	t.Parallel()
	c := newTestClient(t, nil)
	handler := ManageTestPlansHandler(c, false)

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, ManageTestPlansInput{
		Action:     "invalid",
		ProjectKey: "TestProject",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertResultError(t, result, "unknown action")
}
