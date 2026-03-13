package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var attachmentsCmd = &cobra.Command{
	Use:     "attachments",
	Aliases: []string{"attach"},
	Short:   "Manage Azure DevOps work item attachments",
}

var attachmentsListCmd = &cobra.Command{
	Use:   "list <work-item-id>",
	Short: "List attachments on a work item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid work item ID: %s\n", args[0])
			os.Exit(1)
		}
		rels, err := c.ListWorkItemAttachments(project, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, rels, func() {
			if len(rels) == 0 {
				fmt.Println("No attachments found")
				return
			}
			t := NewTable()
			t.Header("Type", "URL")
			for _, r := range rels {
				t.Row(r.Rel, r.URL)
			}
			t.Flush()
			fmt.Printf("\n%d attachments\n", len(rels))
		})
	},
}

func init() {
	RootCmd.AddCommand(attachmentsCmd)
	attachmentsCmd.AddCommand(attachmentsListCmd)

	attachmentsCmd.PersistentFlags().StringP("project", "p", "", "Project name")
}
