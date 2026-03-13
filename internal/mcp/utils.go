package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// resultText returns a successful MCP result with a text content block.
func resultText(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, nil, nil
}

// resultJSON returns a successful MCP result with a JSON-formatted text block.
func resultJSON(data any) (*mcp.CallToolResult, any, error) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return resultError(fmt.Sprintf("marshaling response: %v", err))
	}
	return resultText(string(out))
}

// resultError returns an MCP result with an error message.
func resultError(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Error: " + msg},
		},
		IsError: true,
	}, nil, nil
}
