package devops

import (
	"encoding/base64"
	"net/http"
	"time"
)

// NewTestClient creates a Client that sends requests to the given base URL.
// Intended for use in tests with httptest.NewServer. Since this is in internal/,
// it is not accessible to external consumers.
//
// The test client overrides buildURL to route all requests to the test server,
// so callers don't need a real Azure DevOps instance.
func NewTestClient(baseURL string) *Client {
	encoded := base64.StdEncoding.EncodeToString([]byte(":test-pat"))
	return &Client{
		http:         &http.Client{Timeout: 5 * time.Second},
		organization: "test-org",
		pat:          "test-pat",
		authHeader:   "Basic " + encoded,
		rateLimiter:  NewRateLimiter(10000, time.Millisecond),
	}
}
