package github

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mrbryside/config"
	"log"
)

// PRFile represents a single file change in a PR
type PRFile struct {
	Filename  string `json:"filename"`  // Path of the file
	Additions int    `json:"additions"` // Number of lines added
	Deletions int    `json:"deletions"` // Number of lines deleted
	Changes   int    `json:"changes"`   // Total changes (additions + deletions)
	Patch     string `json:"patch"`     // The diff (first few lines of changes)
}

// FetchPRChanges fetches the changed files in a PR
func FetchPRChanges(owner, repo string, prNumber int) ([]PRFile, error) {
	client := resty.New()
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files", owner, repo, prNumber)

	token := config.Cfg.GitHubToken
	var files []PRFile
	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&files). // Bind the result to the PRFile slice
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status())
	}

	return files, nil
}

// GitHubCommentRequest represents the request body for creating a comment in a PR
type GitHubCommentRequest struct {
	Body     string `json:"body"`
	Path     string `json:"path"` // filePath
	Line     int    `json:"line"`
	CommitID string `json:"commit_id"`
	Side     string `json:"side"`
}

// CommentOnPullRequest creates a comment in a pull request.
func CommentOnPullRequest(owner, repo string, prNumber int, request GitHubCommentRequest) error {
	//token := "YOUR_GITHUB_PERSONAL_ACCESS_TOKEN"
	//owner := "your-github-username"
	//repo := "your-repo-name"
	//prNumber := 1
	//filePath := "example.go"
	//lineNumber := 10
	//commitID := "abc1234def5"
	//comment := "This is a test comment from Resty!"
	baseUrl := "https://api.github.com"
	token := config.Cfg.GitHubToken
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/comments", baseUrl, owner, repo, prNumber)
	client := resty.New().
		SetBaseURL(url).
		SetHeader("Authorization", "token "+token).
		SetHeader("Accept", "application/vnd.github.v3+json")
	resp, err := client.R().
		SetBody(request).
		Post(url)

	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("GitHub API error: %s", resp.String())
	}

	log.Println("Comment created successfully:", resp.String())
	return nil
}
