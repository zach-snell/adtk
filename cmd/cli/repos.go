package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Manage Azure DevOps Git repositories",
}

var reposListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories in a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		repos, err := c.ListRepositories(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, repos, func() {
			t := NewTable()
			t.Header("Name", "Default Branch", "Size", "Remote URL")
			for _, r := range repos {
				t.Row(r.Name, r.DefaultBranch, fmt.Sprintf("%d", r.Size), r.RemoteURL)
			}
			t.Flush()
			fmt.Printf("\n%d repositories\n", len(repos))
		})
	},
}

var reposGetCmd = &cobra.Command{
	Use:   "get <repo>",
	Short: "Get repository details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		repo, err := c.GetRepository(project, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, repo, func() {
			KV("ID", repo.ID)
			KV("Name", repo.Name)
			KV("Default Branch", repo.DefaultBranch)
			KVf("Size", "%d bytes", repo.Size)
			KV("Remote URL", repo.RemoteURL)
			KV("Web URL", repo.WebURL)
			KV("Project", repo.Project.Name)
		})
	},
}

var reposBranchesCmd = &cobra.Command{
	Use:   "branches <repo>",
	Short: "List branches for a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		branches, err := c.ListBranches(project, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, branches, func() {
			t := NewTable()
			t.Header("Name", "Object ID")
			for _, b := range branches {
				t.Row(b.Name, b.ObjectID[:12])
			}
			t.Flush()
			fmt.Printf("\n%d branches\n", len(branches))
		})
	},
}

var reposTreeCmd = &cobra.Command{
	Use:   "tree <repo> [path]",
	Short: "List files in a repository directory",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		treePath := "/"
		if len(args) > 1 {
			treePath = args[1]
		}
		items, err := c.GetTree(project, args[0], treePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, items, func() {
			t := NewTable()
			t.Header("Type", "Path")
			for _, item := range items {
				t.Row(item.GitObjectType, item.Path)
			}
			t.Flush()
		})
	},
}

var reposCommitsCmd = &cobra.Command{
	Use:   "commits <repo>",
	Short: "Search commits in a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		author, _ := cmd.Flags().GetString("author")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		params := map[string]string{}
		if author != "" {
			params["author"] = author
		}
		if from != "" {
			params["fromDate"] = from
		}
		if to != "" {
			params["toDate"] = to
		}
		commits, err := c.SearchCommits(project, args[0], params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(commits)
	},
}

var reposPoliciesCmd = &cobra.Command{
	Use:   "policies <repo>",
	Short: "List branch policies for a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		policies, err := c.ListBranchPolicies(project, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(policies)
	},
}

var reposTagsCmd = &cobra.Command{
	Use:   "tags <repo>",
	Short: "List tags in a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		tags, err := c.ListTags(project, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, tags, func() {
			t := NewTable()
			t.Header("Name", "Object ID")
			for _, tag := range tags {
				name, _ := tag["name"].(string)
				oid, _ := tag["objectId"].(string)
				short := oid
				if len(short) > 12 {
					short = short[:12]
				}
				t.Row(name, short)
			}
			t.Flush()
			fmt.Printf("\n%d tags\n", len(tags))
		})
	},
}

func init() {
	RootCmd.AddCommand(reposCmd)
	reposCmd.AddCommand(reposListCmd, reposGetCmd, reposBranchesCmd, reposTreeCmd, reposCommitsCmd, reposPoliciesCmd, reposTagsCmd)

	reposCmd.PersistentFlags().StringP("project", "p", "", "Project name")

	reposCommitsCmd.Flags().String("author", "", "Filter by commit author")
	reposCommitsCmd.Flags().String("from", "", "Filter commits from date")
	reposCommitsCmd.Flags().String("to", "", "Filter commits to date")
}
