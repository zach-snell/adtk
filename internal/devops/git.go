package devops

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// workItemIDRegex matches a sequence of digits, used to extract work item IDs
// from branch names like feature/12345-description, users/name/12345-fix,
// bugfix/12345, or 12345-description.
var workItemIDRegex = regexp.MustCompile(`(?:^|/)(\d+)`)

// DetectWorkItemFromBranch extracts a work item ID from the current git branch name.
// Patterns: feature/12345-description, users/name/12345-fix, bugfix/12345, 12345-description
// Returns 0 if no work item ID detected.
func DetectWorkItemFromBranch() int {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	branch := strings.TrimSpace(string(output))
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

// DetectWorkItemIDOrArg returns the work item ID from the argument if provided,
// otherwise attempts to detect it from the current git branch.
func DetectWorkItemIDOrArg(arg string) (int, error) {
	if arg != "" {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return 0, fmt.Errorf("invalid work item ID: %s", arg)
		}
		return id, nil
	}

	id := DetectWorkItemFromBranch()
	if id == 0 {
		return 0, fmt.Errorf("no work item ID provided and could not detect from git branch")
	}
	return id, nil
}
