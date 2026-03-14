package devops

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// BaseHost constants for Azure DevOps REST APIs.
const (
	HostMain     = "dev.azure.com"
	HostIdentity = "vssps.dev.azure.com"
	HostRelease  = "vsrm.dev.azure.com"
	HostSearch   = "almsearch.dev.azure.com"

	// DefaultAPIVersion is the Azure DevOps REST API version used for all requests.
	DefaultAPIVersion = "7.1"
)

// Client is the Azure DevOps REST API HTTP client.
// It supports PAT authentication and targets multiple base URLs.
type Client struct {
	http         *http.Client
	organization string
	pat          string
	authHeader   string // pre-computed "Basic base64(:pat)"

	rateLimiter *RateLimiter

	// testBaseURL overrides buildURL to route all requests to a test server.
	// Only set by NewTestClient; zero value means production behavior.
	testBaseURL string
}

// RateLimiter implements a token bucket rate limiter.
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a rate limiter with the specified max tokens and refill rate.
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed and consumes a token.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)
	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}
}

// NewClient creates an Azure DevOps API client with PAT authentication.
// The PAT uses empty-username Basic Auth: base64(":" + pat).
func NewClient(organization, pat string) *Client {
	encoded := base64.StdEncoding.EncodeToString([]byte(":" + pat))
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		organization: organization,
		pat:          pat,
		authHeader:   "Basic " + encoded,
		// ADO uses 200 TSTUs per 5-minute sliding window (usage-based).
		// Use conservative token bucket: 30 requests per minute.
		rateLimiter: NewRateLimiter(30, 2*time.Second),
	}
}

// Organization returns the Azure DevOps organization this client is configured for.
func (c *Client) Organization() string {
	return c.organization
}

// buildURL constructs the full API URL for a given host, optional project, path, and query params.
// The api-version query parameter is always appended.
// When testBaseURL is set (test mode), all requests route to the test server.
func (c *Client) buildURL(host, project, path string, query url.Values) string {
	var base string
	switch {
	case c.testBaseURL != "" && project != "":
		// Test mode with project: route to httptest server
		base = fmt.Sprintf("%s/%s/%s/_apis%s", c.testBaseURL, c.organization, project, path)
	case c.testBaseURL != "":
		// Test mode without project
		base = fmt.Sprintf("%s/%s/_apis%s", c.testBaseURL, c.organization, path)
	case project != "":
		base = fmt.Sprintf("https://%s/%s/%s/_apis%s", host, c.organization, project, path)
	default:
		base = fmt.Sprintf("https://%s/%s/_apis%s", host, c.organization, path)
	}

	if query == nil {
		query = url.Values{}
	}
	if query.Get("api-version") == "" {
		// Some APIs require -preview suffix; callers set it explicitly.
		query.Set("api-version", DefaultAPIVersion)
	}

	return base + "?" + query.Encode()
}

// do executes an HTTP request with auth headers and rate limiting.
func (c *Client) do(method, requestURL string, bodyData []byte, contentType string) (*http.Response, error) {
	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded: please wait and retry")
	}

	var bodyReader io.Reader
	if bodyData != nil {
		bodyReader = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}

// Get performs a GET request to the main ADO API (dev.azure.com).
func (c *Client) Get(project, path string, query url.Values) ([]byte, error) {
	return c.getFrom(HostMain, project, path, query)
}

// GetPreview performs a GET request with a preview API version.
// Some ADO endpoints (e.g. connectionData, graph, comments) require the -preview suffix.
func (c *Client) GetPreview(project, path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("api-version", DefaultAPIVersion+"-preview")
	return c.getFrom(HostMain, project, path, query)
}

// GetIdentity performs a GET request to the identity API (vssps.dev.azure.com).
func (c *Client) GetIdentity(path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("api-version", DefaultAPIVersion+"-preview")
	return c.getFrom(HostIdentity, "", path, query)
}

// PostPreview performs a POST request with a preview API version.
func (c *Client) PostPreview(project, path string, body interface{}) ([]byte, error) {
	return c.postToWithVersion(HostMain, project, path, body, "application/json", DefaultAPIVersion+"-preview")
}

// GetSearch performs a POST request to the search API (almsearch.dev.azure.com).
func (c *Client) PostSearch(project, path string, body interface{}) ([]byte, error) {
	return c.postTo(HostSearch, project, path, body, "application/json")
}

func (c *Client) getFrom(host, project, path string, query url.Values) ([]byte, error) {
	requestURL := c.buildURL(host, project, path, query)
	resp, err := c.do(http.MethodGet, requestURL, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, data)
	}

	return data, nil
}

// Post performs a POST request with a JSON body to the main ADO API.
func (c *Client) Post(project, path string, body interface{}) ([]byte, error) {
	return c.postTo(HostMain, project, path, body, "application/json")
}

func (c *Client) postTo(host, project, path string, body interface{}, contentType string) ([]byte, error) {
	return c.postToWithVersion(host, project, path, body, contentType, "")
}

func (c *Client) postToWithVersion(host, project, path string, body interface{}, contentType, apiVersion string) ([]byte, error) {
	var bodyData []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyData = b
	}

	var query url.Values
	if apiVersion != "" {
		query = url.Values{}
		query.Set("api-version", apiVersion)
	}

	requestURL := c.buildURL(host, project, path, query)
	resp, err := c.do(http.MethodPost, requestURL, bodyData, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respData))
	}

	return respData, nil
}

// Patch performs a PATCH request with a JSON body.
func (c *Client) Patch(project, path string, body interface{}) ([]byte, error) {
	return c.patchWith(project, path, body, "application/json")
}

// PatchPreview performs a PATCH request with a preview API version.
func (c *Client) PatchPreview(project, path string, body interface{}) ([]byte, error) {
	return c.patchWithVersion(project, path, body, "application/json", DefaultAPIVersion+"-preview")
}

// PatchJSONPatch performs a PATCH request with JSON Patch content type.
// Used for work item updates which require application/json-patch+json.
func (c *Client) PatchJSONPatch(project, path string, body interface{}) ([]byte, error) {
	return c.patchWith(project, path, body, "application/json-patch+json")
}

func (c *Client) patchWith(project, path string, body interface{}, contentType string) ([]byte, error) {
	return c.patchWithVersion(project, path, body, contentType, "")
}

func (c *Client) patchWithVersion(project, path string, body interface{}, contentType, apiVersion string) ([]byte, error) {
	var bodyData []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyData = b
	}

	var query url.Values
	if apiVersion != "" {
		query = url.Values{}
		query.Set("api-version", apiVersion)
	}

	requestURL := c.buildURL(HostMain, project, path, query)
	resp, err := c.do(http.MethodPatch, requestURL, bodyData, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respData))
	}

	return respData, nil
}

// PutWithETag performs a PUT request with an If-Match header for optimistic concurrency.
// Used by wiki page updates which require the ETag (git SHA).
func (c *Client) PutWithETag(project, path string, query url.Values, body interface{}, etag string) ([]byte, error) {
	var bodyData []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyData = b
	}

	requestURL := c.buildURL(HostMain, project, path, query)

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded: please wait and retry")
	}

	var bodyReader io.Reader
	if bodyData != nil {
		bodyReader = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequest(http.MethodPut, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("If-Match", etag)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respData))
	}

	return respData, nil
}

// GetWithETag performs a GET and returns both the body and the ETag header.
func (c *Client) GetWithETag(project, path string, query url.Values) (data []byte, etag string, err error) {
	requestURL := c.buildURL(HostMain, project, path, query)
	resp, err := c.do(http.MethodGet, requestURL, nil, "")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, "", fmt.Errorf("reading response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		return nil, "", parseAPIError(resp.StatusCode, data)
	}

	etag = resp.Header.Get("ETag")
	return data, etag, nil
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(project, path string, body interface{}) ([]byte, error) {
	var bodyData []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyData = b
	}

	requestURL := c.buildURL(HostMain, project, path, nil)
	resp, err := c.do(http.MethodPut, requestURL, bodyData, "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respData))
	}

	return respData, nil
}

// Delete performs a DELETE request.
func (c *Client) Delete(project, path string) error {
	requestURL := c.buildURL(HostMain, project, path, nil)
	resp, err := c.do(http.MethodDelete, requestURL, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return parseAPIError(resp.StatusCode, data)
	}

	return nil
}

// GetJSON performs a GET and unmarshals the JSON response into a typed result.
func GetJSON[T any](c *Client, project, path string, query url.Values) (*T, error) {
	data, err := c.Get(project, path, query)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &result, nil
}

// GetAbsolute performs a GET request to an absolute URL (not relative to org).
// Useful for pagination continuation URLs or download links.
func (c *Client) GetAbsolute(absoluteURL string) (*http.Response, error) {
	return c.do(http.MethodGet, absoluteURL, nil, "")
}

// unmarshalJSON is a helper that unmarshals JSON data into a value.
func unmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}
	return nil
}

// readBody reads and returns the full response body.
func readBody(resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	return data, nil
}

func parseAPIError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusForbidden:
		return fmt.Errorf("403 Forbidden: Permission denied. Ensure your PAT has the required scopes. Details: %s", string(body))
	case http.StatusUnauthorized:
		return fmt.Errorf("401 Unauthorized: Authentication failed. Check your PAT. Details: %s", string(body))
	case http.StatusNotFound:
		return fmt.Errorf("404 Not Found: Resource not found. Check project/org name and path. Details: %s", string(body))
	default:
		return fmt.Errorf("API error %d: %s", statusCode, string(body))
	}
}
