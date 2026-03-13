package mcp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// newTestClient creates an httptest server with the given handler and returns
// a devops.Client pointed at it. The server is cleaned up when the test ends.
func newTestClient(t *testing.T, handler http.Handler) *devops.Client {
	t.Helper()
	if handler == nil {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
	}
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return devops.NewTestClient(srv.URL)
}

// getResultText extracts the text from a CallToolResult's first TextContent.
func getResultText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *TextContent, got %T", result.Content[0])
	}
	return tc.Text
}

// assertResultError checks that the result is an error containing substr.
func assertResultError(t *testing.T, result *mcp.CallToolResult, substr string) {
	t.Helper()
	if !result.IsError {
		t.Fatalf("expected error result, got success: %s", getResultText(t, result))
	}
	text := getResultText(t, result)
	if !strings.Contains(text, substr) {
		t.Errorf("error text %q does not contain %q", text, substr)
	}
}

// assertResultSuccess checks that the result is successful and contains substr.
func assertResultSuccess(t *testing.T, result *mcp.CallToolResult, substr string) {
	t.Helper()
	if result.IsError {
		t.Fatalf("expected success result, got error: %s", getResultText(t, result))
	}
	text := getResultText(t, result)
	if !strings.Contains(text, substr) {
		t.Errorf("result text %q does not contain %q", text, substr)
	}
}

// jsonHandler returns an http.HandlerFunc that writes the given JSON string with 200 OK.
func jsonHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}
}

// errorHandler returns an http.HandlerFunc that responds with the given status code.
func errorHandler(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, _ = w.Write([]byte(`{"message":"test error"}`))
	}
}

// muxHandler creates a handler that routes requests by path prefix.
// Falls back to 404 for unmatched paths.
func muxHandler(routes map[string]http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for prefix, handler := range routes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				handler(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"no route matched"}`))
	})
}

// methodMux routes by HTTP method within a single path.
func methodMux(methods map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := methods[r.Method]; ok {
			handler(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
