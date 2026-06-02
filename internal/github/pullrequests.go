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

type ListPullRequestsInput struct {
	Repo    string `json:"repo" jsonschema:"Repository name (slug) within the GitHub owner/organization"`
	State   string `json:"state,omitempty" jsonschema:"Filter by state: 'open', 'closed', 'all'. Default is 'open'."`
	PerPage int    `json:"per_page,omitempty" jsonschema:"Results per page (max 100). Default is 30."`
	Page    int    `json:"page,omitempty" jsonschema:"Page number of results. Default is 1."`
}

func handleListPullRequests(_ context.Context, _ *mcp.CallToolRequest, in ListPullRequestsInput) (*mcp.CallToolResult, any, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=%s&per_page=%d&page=%d",
		apiBase(),
		os.Getenv("GITHUB_OWNER"),
		in.Repo,
		defaultIfEmpty(in.State, "open"),
		maxInt(in.PerPage, 30),
		maxInt(in.Page, 1),
	)
	body, err := get(url, "")
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
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

func defaultIfEmpty(val, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}

func maxInt(val, defaultVal int) int {
	if val <= 0 {
		return defaultVal
	}
	if val > 100 {
		return 100
	}
	return val
}
