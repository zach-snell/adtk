package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListPipelines lists all pipeline definitions in a project.
func (c *Client) ListPipelines(project string, top int) ([]Pipeline, error) {
	query := url.Values{}
	if top > 0 {
		query.Set("$top", fmt.Sprintf("%d", top))
	}
	result, err := GetJSON[PipelineList](c, project, "/pipelines", query)
	if err != nil {
		return nil, fmt.Errorf("listing pipelines: %w", err)
	}
	return result.Value, nil
}

// GetPipeline gets a single pipeline definition by ID.
func (c *Client) GetPipeline(project string, pipelineID int) (*Pipeline, error) {
	path := fmt.Sprintf("/pipelines/%d", pipelineID)
	return GetJSON[Pipeline](c, project, path, nil)
}

// ListPipelineRuns lists runs for a pipeline.
func (c *Client) ListPipelineRuns(project string, pipelineID, top int) ([]PipelineRun, error) {
	path := fmt.Sprintf("/pipelines/%d/runs", pipelineID)
	query := url.Values{}
	if top > 0 {
		query.Set("$top", fmt.Sprintf("%d", top))
	}
	result, err := GetJSON[PipelineRunList](c, project, path, query)
	if err != nil {
		return nil, fmt.Errorf("listing pipeline runs: %w", err)
	}
	return result.Value, nil
}

// GetPipelineRun gets a single pipeline run.
func (c *Client) GetPipelineRun(project string, pipelineID, runID int) (*PipelineRun, error) {
	path := fmt.Sprintf("/pipelines/%d/runs/%d", pipelineID, runID)
	return GetJSON[PipelineRun](c, project, path, nil)
}

// TriggerPipeline triggers a new run of a pipeline.
func (c *Client) TriggerPipeline(project string, pipelineID int, branch string) (*PipelineRun, error) {
	path := fmt.Sprintf("/pipelines/%d/runs", pipelineID)

	body := map[string]interface{}{}
	if branch != "" {
		body["resources"] = map[string]interface{}{
			"repositories": map[string]interface{}{
				"self": map[string]string{
					"refName": "refs/heads/" + branch,
				},
			},
		}
	}

	data, err := c.Post(project, path, body)
	if err != nil {
		return nil, fmt.Errorf("triggering pipeline: %w", err)
	}

	var result PipelineRun
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling pipeline run: %w", err)
	}
	return &result, nil
}

// GetPipelineLogs gets the logs for a pipeline run.
func (c *Client) GetPipelineLogs(project string, pipelineID, runID int) ([]byte, error) {
	path := fmt.Sprintf("/pipelines/%d/runs/%d/logs", pipelineID, runID)
	return c.Get(project, path, nil)
}

// GetPipelineLog gets a specific log by ID for a pipeline run.
func (c *Client) GetPipelineLog(project string, pipelineID, runID, logID int) ([]byte, error) {
	path := fmt.Sprintf("/pipelines/%d/runs/%d/logs/%d", pipelineID, runID, logID)
	return c.Get(project, path, nil)
}
