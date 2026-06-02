package github

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pull_requests",
		Description: "List pull requests in a GitHub repository. Optionally filter by state and paginate.",
	}, handleListPullRequests)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request",
		Description: "Get GitHub pull request metadata: title, description, author, state, source/destination branches, reviewers.",
	}, handleGetPullRequest)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_diff",
		Description: "Get the unified diff of a GitHub pull request (full text of code changes, suitable for code review).",
	}, handleGetPullRequestDiff)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_diffstat",
		Description: "Get per-file change summary of a GitHub pull request (status, lines added/removed).",
	}, handleGetPullRequestDiffstat)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_comments",
		Description: "Get comments on a GitHub pull request (review comments include file path and line; issue comments are general PR discussion).",
	}, handleGetPullRequestComments)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pull_request_commits",
		Description: "Get the list of commits in a GitHub pull request (hash, author, message, date).",
	}, handleGetPullRequestCommits)
}
