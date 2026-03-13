package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListPullRequests lists pull requests for a repository.
func (c *Client) ListPullRequests(project, repoID, status string, top int) ([]PullRequest, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests", repoID)
	query := url.Values{}
	if status != "" {
		query.Set("searchCriteria.status", status)
	}
	if top > 0 {
		query.Set("$top", fmt.Sprintf("%d", top))
	}
	result, err := GetJSON[PullRequestList](c, project, path, query)
	if err != nil {
		return nil, fmt.Errorf("listing pull requests: %w", err)
	}
	return result.Value, nil
}

// GetPullRequest gets a single pull request by ID.
func (c *Client) GetPullRequest(project, repoID string, prID int) (*PullRequest, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests/%d", repoID, prID)
	return GetJSON[PullRequest](c, project, path, nil)
}

// CreatePullRequest creates a new pull request.
func (c *Client) CreatePullRequest(project, repoID string, pr *PullRequest) (*PullRequest, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests", repoID)
	data, err := c.Post(project, path, pr)
	if err != nil {
		return nil, fmt.Errorf("creating pull request: %w", err)
	}

	var result PullRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling pull request: %w", err)
	}
	return &result, nil
}

// UpdatePullRequest updates a pull request (title, description, status, etc).
func (c *Client) UpdatePullRequest(project, repoID string, prID int, update map[string]interface{}) (*PullRequest, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests/%d", repoID, prID)
	data, err := c.Patch(project, path, update)
	if err != nil {
		return nil, fmt.Errorf("updating pull request: %w", err)
	}

	var result PullRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling pull request: %w", err)
	}
	return &result, nil
}

// ListPRThreads lists comment threads on a pull request.
func (c *Client) ListPRThreads(project, repoID string, prID int) ([]PRThread, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests/%d/threads", repoID, prID)
	result, err := GetJSON[PRThreadList](c, project, path, nil)
	if err != nil {
		return nil, fmt.Errorf("listing PR threads: %w", err)
	}
	return result.Value, nil
}

// AddPRComment adds a comment thread to a pull request.
func (c *Client) AddPRComment(project, repoID string, prID int, content string) (*PRThread, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests/%d/threads", repoID, prID)
	body := map[string]interface{}{
		"comments": []map[string]string{
			{"content": content, "commentType": "text"},
		},
		"status": "active",
	}
	data, err := c.Post(project, path, body)
	if err != nil {
		return nil, fmt.Errorf("adding PR comment: %w", err)
	}

	var result PRThread
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling PR thread: %w", err)
	}
	return &result, nil
}

// VotePR submits a vote on a pull request.
// Vote values: 10=approved, 5=approved with suggestions, 0=no vote, -5=waiting, -10=rejected
func (c *Client) VotePR(project, repoID string, prID int, reviewerID string, vote int) error {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests/%d/reviewers/%s", repoID, prID, reviewerID)
	_, err := c.Put(project, path, map[string]int{"vote": vote})
	return err
}

// ListPRReviewers lists reviewers on a pull request.
func (c *Client) ListPRReviewers(project, repoID string, prID int) ([]Reviewer, error) {
	path := fmt.Sprintf("/git/repositories/%s/pullrequests/%d/reviewers", repoID, prID)
	result, err := GetJSON[struct {
		Count int        `json:"count"`
		Value []Reviewer `json:"value"`
	}](c, project, path, nil)
	if err != nil {
		return nil, fmt.Errorf("listing PR reviewers: %w", err)
	}
	return result.Value, nil
}
