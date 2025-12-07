// Package git provides Git repository information retrieval functionality.
package git

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// RepoInfo contains Git repository information
type RepoInfo struct {
	OriginURL string `json:"origin_url"`
	RepoName  string `json:"repo_name"`
	Branch    string `json:"branch,omitempty"`
	Commit    string `json:"commit,omitempty"`
}

// GetRepoInfo returns Git repository information from the current directory
func GetRepoInfo() (*RepoInfo, error) {
	info := &RepoInfo{}

	originURL, err := getGitOriginURL()
	if err != nil {
		return nil, err
	}
	info.OriginURL = originURL
	info.RepoName = extractRepoName(originURL)

	if branch, err := getGitBranch(); err == nil {
		info.Branch = branch
	}

	if commit, err := getGitCommit(); err == nil {
		info.Commit = commit
	}

	return info, nil
}

func getGitOriginURL() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getGitBranch() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getGitCommit() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// extractRepoName extracts repository name from Git URL
// Examples:
//
//	https://github.com/user/repo.git -> user/repo
//	git@github.com:user/repo.git -> user/repo
func extractRepoName(url string) string {
	url = strings.TrimSpace(url)

	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		parts := strings.Split(url, "/")
		// Ensure the URL has at least 5 parts: protocol, empty, domain, owner/path, repo
		if len(parts) >= 5 {
			// Join everything after the domain (parts[3] onwards) to support nested paths
			return strings.Join(parts[3:], "/")
		}
	}

	// Handle SSH URLs with explicit protocol (e.g., ssh://git@github.com:22/user/repo.git)
	if strings.HasPrefix(url, "ssh://") {
		// Remove ssh:// prefix and optional user@host:port
		url = strings.TrimPrefix(url, "ssh://")
		// Find the first "/" which separates host from path
		if idx := strings.Index(url, "/"); idx != -1 {
			// Join everything after the first "/" to support nested paths
			parts := strings.Split(url[idx+1:], "/")
			return strings.Join(parts, "/")
		}
	}

	// Handle SCP-like SSH URLs (e.g., git@github.com:user/repo.git)
	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			// Take the last part after the last colon (handles both with and without port)
			return parts[len(parts)-1]
		}
	}

	return url
}

// IsGitRepo checks if the current directory is inside a Git repository
func IsGitRepo() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}
