package atlassian

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_jira_task",
		Description: "Get a Jira issue by its key (e.g. PROJ-123).",
	}, handleGetJiraTask)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_jira",
		Description: "Search Jira issues using a JQL query.",
	}, handleSearchJira)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_jira_attachment",
		Description: "Download a Jira attachment by its ID (listed in fields.attachment of get_jira_task; images embedded in descriptions and comments are also issue attachments). Returns images as viewable image content and text files as text. Max 5 MB.",
	}, handleGetJiraAttachment)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_confluence_page",
		Description: "Get a Confluence page by its ID (includes body.storage).",
	}, handleGetConfluencePage)
}
