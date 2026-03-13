package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Manage Azure DevOps Kanban boards",
}

var boardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List boards for a project/team",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		team, _ := cmd.Flags().GetString("team")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		boards, err := c.ListBoards(project, team)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, boards, func() {
			t := NewTable()
			t.Header("ID", "Name")
			for _, b := range boards {
				t.Row(b.ID, b.Name)
			}
			t.Flush()
			fmt.Printf("\n%d boards\n", len(boards))
		})
	},
}

var boardsColumnsCmd = &cobra.Command{
	Use:   "columns <board>",
	Short: "List columns on a board",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		team, _ := cmd.Flags().GetString("team")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		cols, err := c.GetBoardColumns(project, team, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, cols, func() {
			t := NewTable()
			t.Header("Name", "Item Limit", "Column Type")
			for _, col := range cols {
				limit := "-"
				if col.ItemLimit > 0 {
					limit = fmt.Sprintf("%d", col.ItemLimit)
				}
				t.Row(col.Name, limit, col.ColumnType)
			}
			t.Flush()
		})
	},
}

func init() {
	RootCmd.AddCommand(boardsCmd)
	boardsCmd.AddCommand(boardsListCmd, boardsColumnsCmd)

	boardsCmd.PersistentFlags().StringP("project", "p", "", "Project name (required)")
	boardsCmd.PersistentFlags().StringP("team", "t", "", "Team name (optional)")
}
