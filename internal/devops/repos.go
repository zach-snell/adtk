package devops

import (
	"fmt"
	"net/url"
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
