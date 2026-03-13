package devops

import (
	"encoding/json"
	"fmt"
)

// SearchWorkItems searches work items using the search API.
func (c *Client) SearchWorkItems(project, searchText string, top int) (*WorkItemSearchResult, error) {
	body := map[string]interface{}{
		"searchText": searchText,
		"$top":       top,
	}
	if top == 0 {
		body["$top"] = 25
	}

	data, err := c.PostSearch(project, "/search/workitemsearchresults", body)
	if err != nil {
		return nil, fmt.Errorf("searching work items: %w", err)
	}

	var result WorkItemSearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling search results: %w", err)
	}
	return &result, nil
}

// SearchCode searches code across repositories.
func (c *Client) SearchCode(project, searchText string, top int) (*CodeSearchResult, error) {
	body := map[string]interface{}{
		"searchText": searchText,
		"$top":       top,
	}
	if top == 0 {
		body["$top"] = 25
	}

	data, err := c.PostSearch(project, "/search/codesearchresults", body)
	if err != nil {
		return nil, fmt.Errorf("searching code: %w", err)
	}

	var result CodeSearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling code search results: %w", err)
	}
	return &result, nil
}

// SearchWiki searches wiki pages.
func (c *Client) SearchWiki(project, searchText string, top int) ([]byte, error) {
	body := map[string]interface{}{
		"searchText": searchText,
		"$top":       top,
	}
	if top == 0 {
		body["$top"] = 25
	}

	return c.PostSearch(project, "/search/wikisearchresults", body)
}
