package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageReposInput defines the input schema for the manage_repos tool.
type ManageReposInput struct {
	Action       string `json:"action" jsonschema:"Action to perform: 'list', 'get', 'list_branches', 'get_file', 'get_tree', 'create_branch', 'search_commits', 'list_policies', 'list_tags', 'create_tag'"`
	ProjectKey   string `json:"project_key,omitempty" jsonschema:"Project name (required for most actions)"`
	RepoID       string `json:"repo_id,omitempty" jsonschema:"Repository name or ID (required for get, list_branches, get_file, get_tree)"`
	FilePath     string `json:"file_path,omitempty" jsonschema:"File path within the repo (for get_file, get_tree)"`
	Version      string `json:"version,omitempty" jsonschema:"Branch name or commit SHA (for get_file)"`
	BranchName   string `json:"branch_name,omitempty" jsonschema:"New branch name (required for create_branch)"`
	SourceBranch string `json:"source_branch,omitempty" jsonschema:"Source branch to create from (required for create_branch)"`
	Author       string `json:"author,omitempty" jsonschema:"Filter commits by author (for search_commits)"`
	FromDate     string `json:"from_date,omitempty" jsonschema:"Filter commits from this date (for search_commits)"`
	ToDate       string `json:"to_date,omitempty" jsonschema:"Filter commits to this date (for search_commits)"`
	TagName      string `json:"tag_name,omitempty" jsonschema:"Tag name (required for create_tag)"`
	CommitSHA    string `json:"commit_sha,omitempty" jsonschema:"Commit SHA to tag (required for create_tag)"`
}

// ManageReposHandler returns the handler for the manage_repos tool.
func ManageReposHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "list":
			return handleRepoList(c, input)
		case "get":
			return handleRepoGet(c, input)
		case "list_branches":
			return handleRepoListBranches(c, input)
		case "get_file":
			return handleRepoGetFile(c, input)
		case "get_tree":
			return handleRepoGetTree(c, input)
		case "create_branch":
			return handleRepoCreateBranch(c, input, enableWrites)
		case "search_commits":
			return handleRepoSearchCommits(c, input)
		case "list_policies":
			return handleRepoListPolicies(c, input)
		case "list_tags":
			return handleRepoListTags(c, input)
		case "create_tag":
			return handleRepoCreateTag(c, input, enableWrites)
		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}

func handleRepoList(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	repos, err := c.ListRepositories(input.ProjectKey)
	if err != nil {
		return resultError(fmt.Sprintf("listing repos: %v", err))
	}
	return resultJSON(repos)
}

func handleRepoGet(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'get' action")
	}
	repo, err := c.GetRepository(input.ProjectKey, input.RepoID)
	if err != nil {
		return resultError(fmt.Sprintf("getting repo: %v", err))
	}
	return resultJSON(repo)
}

func handleRepoListBranches(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'list_branches' action")
	}
	branches, err := c.ListBranches(input.ProjectKey, input.RepoID)
	if err != nil {
		return resultError(fmt.Sprintf("listing branches: %v", err))
	}
	return resultJSON(branches)
}

func handleRepoGetFile(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'get_file' action")
	}
	if input.FilePath == "" {
		return resultError("file_path is required for 'get_file' action")
	}
	item, err := c.GetFileContent(input.ProjectKey, input.RepoID, input.FilePath, input.Version)
	if err != nil {
		return resultError(fmt.Sprintf("getting file: %v", err))
	}
	return resultJSON(item)
}

func handleRepoGetTree(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'get_tree' action")
	}
	treePath := input.FilePath
	if treePath == "" {
		treePath = "/"
	}
	items, err := c.GetTree(input.ProjectKey, input.RepoID, treePath)
	if err != nil {
		return resultError(fmt.Sprintf("getting tree: %v", err))
	}
	return resultJSON(items)
}

func handleRepoCreateBranch(c *devops.Client, input ManageReposInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("create_branch action requires ADTK_ENABLE_WRITES=true")
	}
	if input.RepoID == "" || input.BranchName == "" || input.SourceBranch == "" {
		return resultError("repo_id, branch_name, and source_branch are required for 'create_branch' action")
	}
	if err := c.CreateBranch(input.ProjectKey, input.RepoID, input.BranchName, input.SourceBranch); err != nil {
		return resultError(fmt.Sprintf("creating branch: %v", err))
	}
	return resultText(fmt.Sprintf("Branch %q created from %q", input.BranchName, input.SourceBranch))
}

func handleRepoSearchCommits(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'search_commits' action")
	}
	params := map[string]string{}
	if input.Author != "" {
		params["author"] = input.Author
	}
	if input.FromDate != "" {
		params["fromDate"] = input.FromDate
	}
	if input.ToDate != "" {
		params["toDate"] = input.ToDate
	}
	commits, err := c.SearchCommits(input.ProjectKey, input.RepoID, params)
	if err != nil {
		return resultError(fmt.Sprintf("searching commits: %v", err))
	}
	return resultJSON(commits)
}

func handleRepoListPolicies(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	policies, err := c.ListBranchPolicies(input.ProjectKey, input.RepoID)
	if err != nil {
		return resultError(fmt.Sprintf("listing policies: %v", err))
	}
	return resultJSON(policies)
}

func handleRepoListTags(c *devops.Client, input ManageReposInput) (*sdkmcp.CallToolResult, any, error) {
	if input.RepoID == "" {
		return resultError("repo_id is required for 'list_tags' action")
	}
	tags, err := c.ListTags(input.ProjectKey, input.RepoID)
	if err != nil {
		return resultError(fmt.Sprintf("listing tags: %v", err))
	}
	return resultJSON(tags)
}

func handleRepoCreateTag(c *devops.Client, input ManageReposInput, enableWrites bool) (*sdkmcp.CallToolResult, any, error) {
	if !enableWrites {
		return resultError("create_tag action requires ADTK_ENABLE_WRITES=true")
	}
	if input.RepoID == "" || input.TagName == "" || input.CommitSHA == "" {
		return resultError("repo_id, tag_name, and commit_sha are required for 'create_tag' action")
	}
	if err := c.CreateTag(input.ProjectKey, input.RepoID, input.TagName, input.CommitSHA); err != nil {
		return resultError(fmt.Sprintf("creating tag: %v", err))
	}
	return resultText(fmt.Sprintf("Tag %q created at %s", input.TagName, input.CommitSHA))
}
