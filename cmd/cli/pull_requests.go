package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var prsCmd = &cobra.Command{
	Use:     "pull-requests",
	Aliases: []string{"pr", "prs"},
	Short:   "Manage Azure DevOps pull requests",
}

var prsListCmd = &cobra.Command{
	Use:   "list <repo>",
	Short: "List pull requests for a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		status, _ := cmd.Flags().GetString("status")
		top, _ := cmd.Flags().GetInt("top")
		prs, err := c.ListPullRequests(project, args[0], status, top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, prs, func() {
			t := NewTable()
			t.Header("ID", "Title", "Status", "Source", "Target", "Author")
			for _, pr := range prs {
				t.Row(
					fmt.Sprintf("%d", pr.PullRequestID),
					Truncate(pr.Title, 40),
					pr.Status,
					pr.SourceRefName,
					pr.TargetRefName,
					pr.CreatedBy.DisplayName,
				)
			}
			t.Flush()
			fmt.Printf("\n%d pull requests\n", len(prs))
		})
	},
}

var prsGetCmd = &cobra.Command{
	Use:   "get <repo> <pr-id>",
	Short: "Get pull request details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		prID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid PR ID: %s\n", args[1])
			os.Exit(1)
		}
		pr, err := c.GetPullRequest(project, args[0], prID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, pr, func() {
			KVf("ID", "%d", pr.PullRequestID)
			KV("Title", pr.Title)
			KV("Status", pr.Status)
			KV("Source", pr.SourceRefName)
			KV("Target", pr.TargetRefName)
			KV("Author", pr.CreatedBy.DisplayName)
			KV("Created", pr.CreationDate.Format("2006-01-02 15:04"))
			KV("Draft", FormatBool(pr.IsDraft))
			KV("Merge Status", pr.MergeStatus)
			if pr.Description != "" {
				fmt.Printf("\n%s\n", pr.Description)
			}
		})
	},
}

var prsReviewersCmd = &cobra.Command{
	Use:   "reviewers <repo> <pr-id>",
	Short: "List reviewers on a pull request",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		prID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid PR ID: %s\n", args[1])
			os.Exit(1)
		}
		reviewers, err := c.ListPRReviewers(project, args[0], prID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, reviewers, func() {
			t := NewTable()
			t.Header("Name", "Vote", "Required")
			for _, r := range reviewers {
				vote := voteText(r.Vote)
				t.Row(r.DisplayName, vote, FormatBool(r.IsRequired))
			}
			t.Flush()
		})
	},
}

func init() {
	RootCmd.AddCommand(prsCmd)
	prsCmd.AddCommand(prsListCmd, prsGetCmd, prsReviewersCmd)

	prsCmd.PersistentFlags().StringP("project", "p", "", "Project name")

	prsListCmd.Flags().String("status", "active", "Filter by status: active, completed, abandoned, all")
	prsListCmd.Flags().Int("top", 25, "Max results to return")
}

func voteText(vote int) string {
	switch vote {
	case 10:
		return "Approved"
	case 5:
		return "Approved with suggestions"
	case 0:
		return "No vote"
	case -5:
		return "Waiting for author"
	case -10:
		return "Rejected"
	default:
		return fmt.Sprintf("%d", vote)
	}
}
