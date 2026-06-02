package atlassian

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func atlassianRequest(method, url string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(
		os.Getenv("ATLASSIAN_EMAIL"),
		os.Getenv("ATLASSIAN_API_TOKEN"),
	)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("atlassian API %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// atlassianDownload fetches a binary resource (e.g. an attachment) and
// returns its body together with the response Content-Type. The read is
// capped at maxSize bytes; larger responses return an error.
func atlassianDownload(url string, maxSize int64) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(
		os.Getenv("ATLASSIAN_EMAIL"),
		os.Getenv("ATLASSIAN_API_TOKEN"),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, "", fmt.Errorf("atlassian API %d: %s", resp.StatusCode, string(respBody))
	}

	if resp.ContentLength > maxSize {
		return nil, "", fmt.Errorf("attachment too large: %d bytes (limit %d)", resp.ContentLength, maxSize)
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxSize+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(respBody)) > maxSize {
		return nil, "", fmt.Errorf("attachment too large: exceeds %d bytes limit", maxSize)
	}

	return respBody, resp.Header.Get("Content-Type"), nil
}
