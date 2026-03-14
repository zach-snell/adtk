package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const advancedSecurityAPIVersion = "7.2-preview.1"

// advancedSecurityQuery returns url.Values with the advanced security preview API version.
func advancedSecurityQuery() url.Values {
	q := url.Values{}
	q.Set("api-version", advancedSecurityAPIVersion)
	return q
}

// GetSecurityAlerts retrieves security alerts for a repository.
// GET /{project}/_apis/alert/repositories/{repoId}/alerts?api-version=7.2-preview.1
func (c *Client) GetSecurityAlerts(project, repoID, states, severities string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/alert/repositories/%s/alerts", repoID)
	query := advancedSecurityQuery()
	if states != "" {
		query.Set("criteria.states", states)
	}
	if severities != "" {
		query.Set("criteria.severities", severities)
	}

	data, err := c.getFrom(HostMain, project, path, query)
	if err != nil {
		return nil, fmt.Errorf("getting security alerts: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling security alerts: %w", err)
	}
	return resp.Value, nil
}

// GetSecurityAlertDetails gets a single security alert by ID.
// GET /{project}/_apis/alert/repositories/{repoId}/alerts/{alertId}?api-version=7.2-preview.1
func (c *Client) GetSecurityAlertDetails(project, repoID string, alertID int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/alert/repositories/%s/alerts/%d", repoID, alertID)

	data, err := c.getFrom(HostMain, project, path, advancedSecurityQuery())
	if err != nil {
		return nil, fmt.Errorf("getting security alert details: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling security alert: %w", err)
	}
	return result, nil
}
