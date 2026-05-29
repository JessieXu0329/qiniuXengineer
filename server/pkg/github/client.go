package github

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client handles communications with the GitHub API.
type Client struct {
	restyClient *resty.Client
}

// NewClient creates a configured Resty client for GitHub.
func NewClient() *Client {
	rc := resty.New().
		SetBaseURL("https://api.github.com").
		SetTimeout(15 * time.Second).
		SetHeader("User-Agent", "AI-PR-Review-Assistant").
		SetHeader("Accept", "application/vnd.github.v3.diff") // Get the raw code diff

	// Check if a GITHUB_TOKEN environment variable is present to bypass rate limits
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		rc.SetHeader("Authorization", "token "+token)
	}

	return &Client{
		restyClient: rc,
	}
}

// FetchPRDiff requests the raw unified diff text for a pull request.
func (c *Client) FetchPRDiff(owner, repo string, prNumber int) (string, error) {
	url := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, prNumber)
	
	resp, err := c.restyClient.R().Get(url)
	if err != nil {
		return "", fmt.Errorf("github api request failed: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("github api returned error status %d: %s", resp.StatusCode(), resp.String())
	}

	diffText := resp.String()
	if diffText == "" {
		return "", fmt.Errorf("empty diff returned for PR %d of %s/%s", prNumber, owner, repo)
	}

	return diffText, nil
}
