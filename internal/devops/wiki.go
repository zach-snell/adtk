package devops

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListWikis lists all wikis in a project.
func (c *Client) ListWikis(project string) ([]Wiki, error) {
	result, err := GetJSON[WikiList](c, project, "/wiki/wikis", nil)
	if err != nil {
		return nil, fmt.Errorf("listing wikis: %w", err)
	}
	return result.Value, nil
}

// GetWikiPage gets a wiki page by path.
func (c *Client) GetWikiPage(project, wikiID, pagePath string, includeContent bool) (*WikiPage, error) {
	path := fmt.Sprintf("/wiki/wikis/%s/pages", wikiID)
	query := url.Values{}
	query.Set("path", pagePath)
	if includeContent {
		query.Set("includeContent", "true")
	}
	return GetJSON[WikiPage](c, project, path, query)
}

// CreateWikiPage creates or updates a wiki page.
// Wiki pages are markdown-native — no format conversion needed.
func (c *Client) CreateWikiPage(project, wikiID, pagePath, content string) (*WikiPage, error) {
	path := fmt.Sprintf("/wiki/wikis/%s/pages", wikiID)
	query := url.Values{}
	query.Set("path", pagePath)

	requestURL := c.buildURL(HostMain, project, path, query)
	body := map[string]string{"content": content}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling wiki page: %w", err)
	}

	resp, err := c.do("PUT", requestURL, bodyBytes, "application/json")
	if err != nil {
		return nil, fmt.Errorf("creating wiki page: %w", err)
	}
	defer resp.Body.Close()

	data, err := readBody(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, data)
	}

	var result WikiPage
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling wiki page: %w", err)
	}
	return &result, nil
}

// UpdateWikiPage updates an existing wiki page content.
// Automatically fetches the current ETag for optimistic concurrency.
func (c *Client) UpdateWikiPage(project, wikiID, pagePath, content string, _ int) (*WikiPage, error) {
	apiPath := fmt.Sprintf("/wiki/wikis/%s/pages", wikiID)
	query := url.Values{}
	query.Set("path", pagePath)

	// First GET the page to retrieve its ETag (git SHA).
	_, etag, err := c.GetWithETag(project, apiPath, query)
	if err != nil {
		return nil, fmt.Errorf("getting wiki page ETag: %w", err)
	}
	if etag == "" {
		return nil, fmt.Errorf("wiki page ETag not found; page may not exist at path %q", pagePath)
	}

	body := map[string]string{"content": content}
	data, err := c.PutWithETag(project, apiPath, query, body, etag)
	if err != nil {
		return nil, fmt.Errorf("updating wiki page: %w", err)
	}

	var result WikiPage
	if err := unmarshalJSON(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteWikiPage deletes a wiki page by path.
func (c *Client) DeleteWikiPage(project, wikiID, pagePath string) error {
	path := fmt.Sprintf("/wiki/wikis/%s/pages", wikiID)
	query := url.Values{}
	query.Set("path", pagePath)

	requestURL := c.buildURL(HostMain, project, path, query)
	resp, err := c.do("DELETE", requestURL, nil, "")
	if err != nil {
		return fmt.Errorf("deleting wiki page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		data, _ := readBody(resp)
		return parseAPIError(resp.StatusCode, data)
	}

	return nil
}
