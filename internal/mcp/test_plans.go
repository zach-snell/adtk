package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageTestPlansInput defines the input schema for the manage_test_plans tool.
type ManageTestPlansInput struct {
	Action        string `json:"action" jsonschema:"Action to perform: 'list_plans', 'create_plan', 'list_suites', 'create_suite', 'list_cases', 'get_test_results'"`
	ProjectKey    string `json:"project_key,omitempty" jsonschema:"Project name (required)"`
	PlanID        int    `json:"plan_id,omitempty" jsonschema:"Test plan ID"`
	SuiteID       int    `json:"suite_id,omitempty" jsonschema:"Test suite ID"`
	ParentSuiteID int    `json:"parent_suite_id,omitempty" jsonschema:"Parent suite ID (for create_suite)"`
	BuildID       int    `json:"build_id,omitempty" jsonschema:"Build ID (for get_test_results)"`
	Name          string `json:"name,omitempty" jsonschema:"Name (for create_plan, create_suite)"`
	Iteration     string `json:"iteration,omitempty" jsonschema:"Iteration path (for create_plan)"`
}

// ManageTestPlansHandler returns the handler for the manage_test_plans tool.
func ManageTestPlansHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageTestPlansInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageTestPlansInput) (*sdkmcp.CallToolResult, any, error) {
		if input.ProjectKey == "" {
			return resultError("project_key is required")
		}

		switch input.Action {
		case "list_plans":
			plans, err := c.ListTestPlans(input.ProjectKey)
			if err != nil {
				return resultError(fmt.Sprintf("listing test plans: %v", err))
			}
			return resultJSON(plans)
		case "create_plan":
			if !enableWrites {
				return resultError("create_plan action requires ADTK_ENABLE_WRITES=true")
			}
			if input.Name == "" {
				return resultError("name is required for 'create_plan' action")
			}
			plan, err := c.CreateTestPlan(input.ProjectKey, input.Name, input.Iteration)
			if err != nil {
				return resultError(fmt.Sprintf("creating test plan: %v", err))
			}
			return resultJSON(plan)
		case "list_suites":
			if input.PlanID == 0 {
				return resultError("plan_id is required for 'list_suites' action")
			}
			suites, err := c.ListTestSuites(input.ProjectKey, input.PlanID)
			if err != nil {
				return resultError(fmt.Sprintf("listing test suites: %v", err))
			}
			return resultJSON(suites)
		case "create_suite":
			if !enableWrites {
				return resultError("create_suite action requires ADTK_ENABLE_WRITES=true")
			}
			if input.PlanID == 0 || input.Name == "" {
				return resultError("plan_id and name are required for 'create_suite' action")
			}
			suite, err := c.CreateTestSuite(input.ProjectKey, input.PlanID, input.ParentSuiteID, input.Name)
			if err != nil {
				return resultError(fmt.Sprintf("creating test suite: %v", err))
			}
			return resultJSON(suite)
		case "list_cases":
			if input.PlanID == 0 || input.SuiteID == 0 {
				return resultError("plan_id and suite_id are required for 'list_cases' action")
			}
			cases, err := c.ListTestCases(input.ProjectKey, input.PlanID, input.SuiteID)
			if err != nil {
				return resultError(fmt.Sprintf("listing test cases: %v", err))
			}
			return resultJSON(cases)
		case "get_test_results":
			if input.BuildID == 0 {
				return resultError("build_id is required for 'get_test_results' action")
			}
			results, err := c.GetTestResultsForBuild(input.ProjectKey, input.BuildID)
			if err != nil {
				return resultError(fmt.Sprintf("getting test results: %v", err))
			}
			return resultJSON(results)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
