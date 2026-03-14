package devops_test

import (
	"regexp"
	"strconv"
	"testing"
)

// workItemIDRegex matches the same pattern as the production code in git.go.
// We test the regex logic directly with known branch patterns.
var workItemIDRegex = regexp.MustCompile(`(?:^|/)(\d+)`)

// extractWorkItemID simulates DetectWorkItemFromBranch but without exec.Command.
// This lets us test the regex parsing logic in isolation.
func extractWorkItemID(branch string) int {
	if branch == "" || branch == "HEAD" {
		return 0
	}
	matches := workItemIDRegex.FindStringSubmatch(branch)
	if len(matches) < 2 {
		return 0
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}
	return id
}

func TestDetectWorkItemFromBranch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		branch string
		want   int
	}{
		{"feature/12345-add-login", 12345},
		{"users/zach/67890-fix-bug", 67890},
		{"bugfix/111", 111},
		{"main", 0},
		{"feature/no-number", 0},
		{"12345-at-start", 12345},
		{"hotfix/99999-urgent", 99999},
		{"refs/heads/feature/42-task", 42},
		{"", 0},
		{"HEAD", 0},
		{"dev", 0},
		{"release/v1.2.3", 0},
		{"users/jane/0-edge-case", 0},
	}
	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			t.Parallel()
			got := extractWorkItemID(tt.branch)
			if got != tt.want {
				t.Errorf("extractWorkItemID(%q) = %d, want %d", tt.branch, got, tt.want)
			}
		})
	}
}
