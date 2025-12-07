// Package git provides Git repository information retrieval functionality.
package git

import (
	"os/exec"
	"strings"
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
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getGitBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getGitCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// extractRepoName extracts repository name from Git URL
// Examples:
//   https://github.com/user/repo.git -> user/repo
//   git@github.com:user/repo.git -> user/repo
func extractRepoName(url string) string {
	url = strings.TrimSpace(url)

	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		parts := strings.Split(url, "/")
		// Ensure the URL has at least 5 parts: protocol, domain, owner, repo
		if len(parts) >= 5 {
			return parts[3] + "/" + parts[4]
		}
	}

	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}

	return url
}

// IsGitRepo checks if the current directory is inside a Git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}
