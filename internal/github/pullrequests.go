package github

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-kit/internal/mcpkit"
)

type PullRequestInput struct {
	Repo string `json:"repo" jsonschema:"Repository name (slug) within the GitHub owner/organization (e.g. 'my-repo' in github.com/my-org/my-repo/pull/123)"`
	ID   int    `json:"id" jsonschema:"Pull request number (the number in the PR URL)"`
}

func handleGetPullRequest(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	body, err := get(prURL(in.Repo, in.ID, ""), "")
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
}

func handleGetPullRequestDiff(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	body, err := get(prURL(in.Repo, in.ID, ""), "application/vnd.github.v3.diff")
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
}

func handleGetPullRequestDiffstat(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	body, err := get(prURL(in.Repo, in.ID, "files"), "")
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
}

func handleGetPullRequestComments(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	reviewBody, err := get(prURL(in.Repo, in.ID, "comments"), "")
	if err != nil {
		return nil, nil, err
	}
	issueBody, err := get(issueURL(in.Repo, in.ID, "comments"), "")
	if err != nil {
		return nil, nil, err
	}
	combined := fmt.Sprintf(`{"review_comments":%s,"issue_comments":%s}`, reviewBody, issueBody)
	return mcpkit.TextResult(combined), nil, nil
}

func handleGetPullRequestCommits(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	body, err := get(prURL(in.Repo, in.ID, "commits"), "")
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
}

// issueURL builds the URL for an issue resource (PRs are also issues on GitHub).
func issueURL(repo string, id int, subpath string) string {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d",
		apiBase(),
		os.Getenv("GITHUB_OWNER"),
		repo,
		id,
	)
	if subpath != "" {
		url += "/" + subpath
	}
	return url
}
