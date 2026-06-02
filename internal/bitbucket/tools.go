package bitbucket

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pull_requests",
		Description: "List pull requests in a Bitbucket repository. Optionally filter by state and paginate.",
	}, handleListPullRequests)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request",
		Description: "Get Bitbucket pull request metadata: title, description, author, state, source/destination branches, reviewers.",
	}, handleGetPullRequest)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_diff",
		Description: "Get the unified diff of a Bitbucket pull request (full text of code changes, suitable for code review).",
	}, handleGetPullRequestDiff)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_diffstat",
		Description: "Get per-file change summary of a Bitbucket pull request (status, lines added/removed).",
	}, handleGetPullRequestDiffstat)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_comments",
		Description: "Get comments on a Bitbucket pull request (inline code comments include file path and line).",
	}, handleGetPullRequestComments)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_commits",
		Description: "Get the list of commits in a Bitbucket pull request (hash, author, message, date).",
	}, handleGetPullRequestCommits)
}
