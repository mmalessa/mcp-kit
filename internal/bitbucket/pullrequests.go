package bitbucket

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-kit/internal/mcpkit"
)

type PullRequestInput struct {
	Repo string `json:"repo" jsonschema:"Repository slug - the segment after the workspace in the PR URL (e.g. 'my-repo' in bitbucket.org/my-workspace/my-repo/pull-requests/123)"`
	ID   int    `json:"id" jsonschema:"Pull request ID (the number in the PR URL)"`
}

type ListPullRequestsInput struct {
	Repo     string `json:"repo" jsonschema:"Repository slug within the Bitbucket workspace"`
	State    string `json:"state,omitempty" jsonschema:"Filter by state: 'OPEN', 'MERGED', 'DECLINED', 'SUPERSEDED'. Default is 'OPEN'."`
	PageLen  int    `json:"pagelen,omitempty" jsonschema:"Results per page (max 50). Default is 10."`
	Page     int    `json:"page,omitempty" jsonschema:"Page number of results. Default is 1."`
}

func handleListPullRequests(_ context.Context, _ *mcp.CallToolRequest, in ListPullRequestsInput) (*mcp.CallToolResult, any, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?state=%s&pagelen=%d&page=%d",
		apiBase(),
		os.Getenv("BITBUCKET_WORKSPACE"),
		in.Repo,
		defaultIfEmpty(in.State, "OPEN"),
		clampInt(in.PageLen, 1, 50, 10),
		maxInt(in.Page, 1),
	)
	body, err := get(url)
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
}

func fetchPR(repo string, id int, subpath string) (*mcp.CallToolResult, any, error) {
	url := prURL(repo, id)
	if subpath != "" {
		url += "/" + subpath
	}
	body, err := get(url)
	if err != nil {
		return nil, nil, err
	}
	return mcpkit.TextResult(string(body)), nil, nil
}

func handleGetPullRequest(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	return fetchPR(in.Repo, in.ID, "")
}

func handleGetPullRequestDiff(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	return fetchPR(in.Repo, in.ID, "diff")
}

func handleGetPullRequestDiffstat(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	return fetchPR(in.Repo, in.ID, "diffstat")
}

func handleGetPullRequestComments(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	return fetchPR(in.Repo, in.ID, "comments")
}

func handleGetPullRequestCommits(_ context.Context, _ *mcp.CallToolRequest, in PullRequestInput) (*mcp.CallToolResult, any, error) {
	return fetchPR(in.Repo, in.ID, "commits")
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
	return val
}

func clampInt(val, min, max, defaultVal int) int {
	if val <= 0 {
		return defaultVal
	}
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
