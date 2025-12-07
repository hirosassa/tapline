package git

import (
	"testing"
)

func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "HTTPS GitHub URL with .git",
			url:      "https://github.com/hirosassa/tapline.git",
			expected: "hirosassa/tapline",
		},
		{
			name:     "SSH GitHub URL with .git",
			url:      "git@github.com:hirosassa/tapline.git",
			expected: "hirosassa/tapline",
		},
		{
			name:     "HTTPS GitHub URL without .git",
			url:      "https://github.com/user/my-repo",
			expected: "user/my-repo",
		},
		{
			name:     "SSH GitLab URL with nested path",
			url:      "git@gitlab.com:company/team/project.git",
			expected: "company/team/project",
		},
		{
			name:     "HTTPS GitLab URL",
			url:      "https://gitlab.com/company/project.git",
			expected: "company/project",
		},
		{
			name:     "HTTPS GitLab URL with nested path",
			url:      "https://gitlab.com/company/team/project.git",
			expected: "company/team/project",
		},
		{
			name:     "SSH URL with protocol and port",
			url:      "ssh://git@github.com:22/user/repo.git",
			expected: "user/repo",
		},
		{
			name:     "SSH URL with protocol, port, and nested path",
			url:      "ssh://git@gitlab.com:22/company/team/project.git",
			expected: "company/team/project",
		},
		{
			name:     "URL with trailing whitespace",
			url:      "  https://github.com/user/repo.git  ",
			expected: "user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRepoName(tt.url)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetRepoInfo(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a Git repository")
	}

	info, err := GetRepoInfo()
	if err != nil {
		t.Fatalf("Failed to get repo info: %v", err)
	}

	if info.OriginURL == "" {
		t.Error("Expected OriginURL to be set")
	}

	if info.RepoName == "" {
		t.Error("Expected RepoName to be set")
	}

	t.Logf("Repository info: URL=%s, Name=%s, Branch=%s, Commit=%s",
		info.OriginURL, info.RepoName, info.Branch, info.Commit)
}

func TestIsGitRepo(t *testing.T) {
	result := IsGitRepo()
	if !result {
		t.Error("Expected to be in a Git repository")
	}
}

func TestGetGitBranch(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a Git repository")
	}

	branch, err := getGitBranch()
	if err != nil {
		t.Fatalf("Failed to get branch: %v", err)
	}

	if branch == "" {
		t.Error("Expected branch name to be non-empty")
	}

	t.Logf("Current branch: %s", branch)
}

func TestGetGitCommit(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a Git repository")
	}

	commit, err := getGitCommit()
	if err != nil {
		t.Fatalf("Failed to get commit: %v", err)
	}

	if commit == "" {
		t.Error("Expected commit hash to be non-empty")
	}

	if len(commit) != 7 {
		t.Errorf("Expected short commit hash (7 chars), got %d chars: %s", len(commit), commit)
	}

	t.Logf("Current commit: %s", commit)
}

func TestGetGitOriginURL(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a Git repository")
	}

	url, err := getGitOriginURL()
	if err != nil {
		t.Fatalf("Failed to get origin URL: %v", err)
	}

	if url == "" {
		t.Error("Expected origin URL to be non-empty")
	}

	t.Logf("Origin URL: %s", url)
}
