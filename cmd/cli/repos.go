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

func init() {
	RootCmd.AddCommand(reposCmd)
	reposCmd.AddCommand(reposListCmd, reposGetCmd, reposBranchesCmd, reposTreeCmd)

	reposCmd.PersistentFlags().StringP("project", "p", "", "Project name")
}
