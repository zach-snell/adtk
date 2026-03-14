package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const testPlanAPIVersion = "7.1-preview.1"

// testPlanQuery returns url.Values with the test plan preview API version.
func testPlanQuery() url.Values {
	q := url.Values{}
	q.Set("api-version", testPlanAPIVersion)
	return q
}

// ListTestPlans lists test plans in a project.
// GET /{project}/_apis/testplan/plans?api-version=7.1-preview.1
func (c *Client) ListTestPlans(project string) ([]map[string]interface{}, error) {
	data, err := c.getFrom(HostMain, project, "/testplan/plans", testPlanQuery())
	if err != nil {
		return nil, fmt.Errorf("listing test plans: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling test plans: %w", err)
	}
	return resp.Value, nil
}

// CreateTestPlan creates a test plan in a project.
// POST /{project}/_apis/testplan/plans?api-version=7.1-preview.1
func (c *Client) CreateTestPlan(project, name, iteration string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"name": name,
	}
	if iteration != "" {
		body["iteration"] = iteration
	}

	data, err := c.postToWithVersion(HostMain, project, "/testplan/plans", body, "application/json", testPlanAPIVersion)
	if err != nil {
		return nil, fmt.Errorf("creating test plan: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling test plan: %w", err)
	}
	return result, nil
}

// ListTestSuites lists test suites in a test plan.
// GET /{project}/_apis/testplan/Plans/{planId}/suites?api-version=7.1-preview.1
func (c *Client) ListTestSuites(project string, planID int) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/testplan/Plans/%d/suites", planID)
	data, err := c.getFrom(HostMain, project, path, testPlanQuery())
	if err != nil {
		return nil, fmt.Errorf("listing test suites: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling test suites: %w", err)
	}
	return resp.Value, nil
}

// CreateTestSuite creates a test suite in a test plan.
// POST /{project}/_apis/testplan/Plans/{planId}/suites?api-version=7.1-preview.1
func (c *Client) CreateTestSuite(project string, planID, parentSuiteID int, name string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/testplan/Plans/%d/suites", planID)
	body := map[string]interface{}{
		"suiteType": "staticTestSuite",
		"name":      name,
	}
	if parentSuiteID > 0 {
		body["parentSuite"] = map[string]interface{}{
			"id": parentSuiteID,
		}
	}

	data, err := c.postToWithVersion(HostMain, project, path, body, "application/json", testPlanAPIVersion)
	if err != nil {
		return nil, fmt.Errorf("creating test suite: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling test suite: %w", err)
	}
	return result, nil
}

// ListTestCases lists test cases in a test suite.
// GET /{project}/_apis/testplan/Plans/{planId}/Suites/{suiteId}/TestCase?api-version=7.1-preview.1
func (c *Client) ListTestCases(project string, planID, suiteID int) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/testplan/Plans/%d/Suites/%d/TestCase", planID, suiteID)
	data, err := c.getFrom(HostMain, project, path, testPlanQuery())
	if err != nil {
		return nil, fmt.Errorf("listing test cases: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling test cases: %w", err)
	}
	return resp.Value, nil
}

// GetTestResultsForBuild gets test results (runs) for a build.
// GET /{project}/_apis/test/Runs?buildUri=vstfs:///Build/Build/{buildId}&api-version=7.1
func (c *Client) GetTestResultsForBuild(project string, buildID int) ([]map[string]interface{}, error) {
	query := url.Values{}
	query.Set("buildUri", fmt.Sprintf("vstfs:///Build/Build/%d", buildID))

	data, err := c.Get(project, "/test/Runs", query)
	if err != nil {
		return nil, fmt.Errorf("getting test results for build: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling test results: %w", err)
	}
	return resp.Value, nil
}
