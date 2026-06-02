package atlassian

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-kit/internal/mcpkit"
)

// maxAttachmentSize caps attachment downloads; base64 encoding inflates
// the payload by ~33%, so large files would blow up the agent's context.
const maxAttachmentSize = 5 << 20 // 5 MB

type GetJiraTaskInput struct {
	IssueKey string `json:"issue_key" jsonschema:"Jira issue key, e.g. PROJ-123"`
}

type GetJiraAttachmentInput struct {
	AttachmentID string `json:"attachment_id" jsonschema:"Jira attachment ID, found in the issue's fields.attachment[].id"`
}

type SearchJiraInput struct {
	JQL string `json:"jql" jsonschema:"JQL query string, e.g. \"assignee = currentUser() ORDER BY created DESC\""`
}

func handleGetJiraTask(_ context.Context, _ *mcp.CallToolRequest, in GetJiraTaskInput) (*mcp.CallToolResult, any, error) {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s",
		os.Getenv("ATLASSIAN_BASE_URL"),
		in.IssueKey,
	)

	body, err := atlassianRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	return mcpkit.TextResult(string(body)), nil, nil
}

func handleSearchJira(_ context.Context, _ *mcp.CallToolRequest, in SearchJiraInput) (*mcp.CallToolResult, any, error) {
	url := fmt.Sprintf("%s/rest/api/3/search/jql", os.Getenv("ATLASSIAN_BASE_URL"))

	reqBody, err := json.Marshal(map[string]any{
		"jql":    in.JQL,
		"fields": []string{"*all"},
	})
	if err != nil {
		return nil, nil, err
	}

	body, err := atlassianRequest("POST", url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	return mcpkit.TextResult(string(body)), nil, nil
}

func handleGetJiraAttachment(_ context.Context, _ *mcp.CallToolRequest, in GetJiraAttachmentInput) (*mcp.CallToolResult, any, error) {
	url := fmt.Sprintf("%s/rest/api/3/attachment/content/%s",
		os.Getenv("ATLASSIAN_BASE_URL"),
		in.AttachmentID,
	)

	body, contentType, err := atlassianDownload(url, maxAttachmentSize)
	if err != nil {
		return nil, nil, err
	}

	// Strip parameters like "; charset=utf-8" from the media type.
	mediaType := contentType
	if i := strings.Index(mediaType, ";"); i >= 0 {
		mediaType = strings.TrimSpace(mediaType[:i])
	}

	switch {
	case strings.HasPrefix(mediaType, "image/"):
		return mcpkit.ImageResult(body, mediaType), nil, nil
	case strings.HasPrefix(mediaType, "text/"),
		mediaType == "application/json",
		mediaType == "application/xml",
		strings.HasSuffix(mediaType, "+json"),
		strings.HasSuffix(mediaType, "+xml"):
		return mcpkit.TextResult(string(body)), nil, nil
	default:
		return nil, nil, fmt.Errorf("unsupported attachment type %q (%d bytes); only image and text-based attachments can be retrieved", mediaType, len(body))
	}
}
