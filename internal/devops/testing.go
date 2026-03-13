package devops

import (
	"encoding/base64"
	"net/http"
	"time"
)

// testClient wraps Client and overrides buildURL to route all requests to a test server.
// This is not exported because it's in internal/ and only used by tests.

// NewTestClient creates a Client that sends requests to the given base URL.
// Intended for use in tests with httptest.NewServer. Since this is in internal/,
// it is not accessible to external consumers.
//
// The test client stores the base URL and overrides buildURL to route requests
// to the test server instead of real Azure DevOps endpoints.
func NewTestClient(baseURL string) *Client {
	encoded := base64.StdEncoding.EncodeToString([]byte(":test-pat"))
	c := &Client{
		http:         &http.Client{Timeout: 5 * time.Second},
		organization: "test-org",
		pat:          "test-pat",
		authHeader:   "Basic " + encoded,
		rateLimiter:  NewRateLimiter(10000, time.Millisecond),
		testBaseURL:  baseURL,
	}
	return c
}
