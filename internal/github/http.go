package github

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const defaultAPIBase = "https://api.github.com"

func apiBase() string {
	if v := os.Getenv("GITHUB_API_URL"); v != "" {
		return v
	}
	return defaultAPIBase
}

// get performs a GET against the GitHub API with Bearer token auth.
func get(url string, accept string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// prURL builds the URL for a pull request resource.
// Owner comes from GITHUB_OWNER env, repo is passed per-call.
func prURL(repo string, id int, subpath string) string {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d",
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
