package mcpkit

import "github.com/modelcontextprotocol/go-sdk/mcp"

// TextResult wraps text content as an MCP tool result.
func TextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

// ImageResult wraps binary image data as an MCP tool result.
func ImageResult(data []byte, mimeType string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.ImageContent{Data: data, MIMEType: mimeType}},
	}
}
