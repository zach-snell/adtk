package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ListRepositories lists all Git repositories in a project.
func (c *Client) ListRepositories(project string) ([]GitRepository, error) {
	result, err := GetJSON[GitRepositoryList](c, project, "/git/repositories", nil)
	if err != nil {
		return nil, fmt.Errorf("listing repositories: %w", err)
	}
	return result.Value, nil
}

// GetRepository gets a single repository by name or ID.
func (c *Client) GetRepository(project, repoID string) (*GitRepository, error) {
	path := fmt.Sprintf("/git/repositories/%s", repoID)
	return GetJSON[GitRepository](c, project, path, nil)
}

// ListBranches lists all branches (refs) for a repository.
func (c *Client) ListBranches(project, repoID string) ([]GitRef, error) {
	path := fmt.Sprintf("/git/repositories/%s/refs", repoID)
	query := url.Values{}
	query.Set("filter", "heads/")
	result, err := GetJSON[GitRefList](c, project, path, query)
	if err != nil {
		return nil, fmt.Errorf("listing branches: %w", err)
	}
	return result.Value, nil
}

// GetFileContent gets the content of a file at a given path and optional version.
func (c *Client) GetFileContent(project, repoID, filePath, version string) (*GitItem, error) {
	path := fmt.Sprintf("/git/repositories/%s/items", repoID)
	query := url.Values{}
	query.Set("path", filePath)
	query.Set("includeContent", "true")
	if version != "" {
		query.Set("version", version)
	}
	return GetJSON[GitItem](c, project, path, query)
}

// GetTree gets the tree (directory listing) for a repository at a given path.
func (c *Client) GetTree(project, repoID, treePath string) ([]GitItem, error) {
	path := fmt.Sprintf("/git/repositories/%s/items", repoID)
	query := url.Values{}
	query.Set("scopePath", treePath)
	query.Set("recursionLevel", "OneLevel")

	data, err := c.Get(project, path, query)
	if err != nil {
		return nil, fmt.Errorf("getting tree: %w", err)
	}

	var result struct {
		Count int       `json:"count"`
		Value []GitItem `json:"value"`
	}
	if err := unmarshalJSON(data, &result); err != nil {
		return nil, err
	}
	return result.Value, nil
}

// CreateBranch creates a new branch from a source branch.
func (c *Client) CreateBranch(project, repoID, branchName, sourceBranch string) error {
	// First, resolve the source branch objectId
	sourceOID, err := c.resolveSourceBranch(project, repoID, sourceBranch)
	if err != nil {
		return err
	}

	// Create the new ref
	path := fmt.Sprintf("/git/repositories/%s/refs", repoID)
	refName := branchName
	if !strings.HasPrefix(refName, "refs/heads/") {
		refName = "refs/heads/" + refName
	}
	body := []map[string]string{
		{
			"name":        refName,
			"oldObjectId": "0000000000000000000000000000000000000000",
			"newObjectId": sourceOID,
		},
	}
	_, err = c.Post(project, path, body)
	if err != nil {
		return fmt.Errorf("creating branch: %w", err)
	}
	return nil
}

// resolveSourceBranch gets the objectId for a source branch name.
func (c *Client) resolveSourceBranch(project, repoID, sourceBranch string) (string, error) {
	branches, err := c.ListBranches(project, repoID)
	if err != nil {
		return "", fmt.Errorf("resolving source branch: %w", err)
	}

	refName := sourceBranch
	if !strings.HasPrefix(refName, "refs/heads/") {
		refName = "refs/heads/" + refName
	}

	for _, b := range branches {
		if b.Name == refName {
			return b.ObjectID, nil
		}
	}
	return "", fmt.Errorf("source branch %q not found", sourceBranch)
}

// ListBranchPolicies lists branch policy configurations for a project.
// Optionally filters by repository ID.
// GET /{project}/_apis/policy/configurations?api-version=7.1
func (c *Client) ListBranchPolicies(project, repoID string) ([]map[string]interface{}, error) {
	data, err := c.Get(project, "/policy/configurations", nil)
	if err != nil {
		return nil, fmt.Errorf("listing branch policies: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling policies: %w", err)
	}

	// Filter by repository if specified
	if repoID == "" {
		return resp.Value, nil
	}
	return filterPoliciesByRepo(resp.Value, repoID), nil
}

// filterPoliciesByRepo filters policy configurations to those matching a repo ID.
func filterPoliciesByRepo(policies []map[string]interface{}, repoID string) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, p := range policies {
		settings, ok := p["settings"].(map[string]interface{})
		if !ok {
			continue
		}
		scopes, ok := settings["scope"].([]interface{})
		if !ok {
			continue
		}
		for _, s := range scopes {
			scope, ok := s.(map[string]interface{})
			if !ok {
				continue
			}
			if rid, ok := scope["repositoryId"].(string); ok && rid == repoID {
				filtered = append(filtered, p)
				break
			}
		}
	}
	return filtered
}

// ListTags lists git tags in a repository.
// GET /{project}/_apis/git/repositories/{repoId}/refs?filter=tags/&api-version=7.1
func (c *Client) ListTags(project, repoID string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/git/repositories/%s/refs", repoID)
	query := url.Values{}
	query.Set("filter", "tags/")

	data, err := c.Get(project, path, query)
	if err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}

	var resp struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling tags: %w", err)
	}
	return resp.Value, nil
}

// CreateTag creates a git tag pointing to a specific commit.
// POST /{project}/_apis/git/repositories/{repoId}/refs?api-version=7.1
func (c *Client) CreateTag(project, repoID, tagName, commitSHA string) error {
	path := fmt.Sprintf("/git/repositories/%s/refs", repoID)
	refName := tagName
	if !strings.HasPrefix(refName, "refs/tags/") {
		refName = "refs/tags/" + refName
	}
	body := []map[string]string{
		{
			"name":        refName,
			"oldObjectId": "0000000000000000000000000000000000000000",
			"newObjectId": commitSHA,
		},
	}
	_, err := c.Post(project, path, body)
	if err != nil {
		return fmt.Errorf("creating tag: %w", err)
	}
	return nil
}

// SearchCommits searches for commits with filters using the commitsbatch API.
func (c *Client) SearchCommits(project, repoID string, params map[string]string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/git/repositories/%s/commitsbatch", repoID)
	body := map[string]interface{}{}

	if author, ok := params["author"]; ok && author != "" {
		body["author"] = author
	}
	if fromDate, ok := params["fromDate"]; ok && fromDate != "" {
		body["fromDate"] = fromDate
	}
	if toDate, ok := params["toDate"]; ok && toDate != "" {
		body["toDate"] = toDate
	}

	data, err := c.Post(project, path, body)
	if err != nil {
		return nil, fmt.Errorf("searching commits: %w", err)
	}

	var result struct {
		Count int                      `json:"count"`
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling commits: %w", err)
	}
	return result.Value, nil
}
