package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Azure DevOps (code, work items, wiki)",
}

var searchCodeCmd = &cobra.Command{
	Use:   "code <query>",
	Short: "Search code across repositories",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		top, _ := cmd.Flags().GetInt("top")
		result, err := c.SearchCode(project, args[0], top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, result, func() {
			t := NewTable()
			t.Header("File", "Path", "Repository", "Project")
			for _, r := range result.Results {
				t.Row(r.FileName, r.Path, r.Repository.Name, r.Project.Name)
			}
			t.Flush()
			fmt.Printf("\n%d results\n", result.Count)
		})
	},
}

var searchWICmd = &cobra.Command{
	Use:   "work-items <query>",
	Short: "Search work items by text",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		top, _ := cmd.Flags().GetInt("top")
		result, err := c.SearchWorkItems(project, args[0], top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, result, func() {
			t := NewTable()
			t.Header("ID", "Type", "Title", "State", "Project")
			for _, r := range result.Results {
				t.Row(
					r.Fields["system.id"],
					r.Fields["system.workitemtype"],
					Truncate(r.Fields["system.title"], 40),
					r.Fields["system.state"],
					r.Project.Name,
				)
			}
			t.Flush()
			fmt.Printf("\n%d results\n", result.Count)
		})
	},
}

var searchWiqlCmd = &cobra.Command{
	Use:   "wiql <query>",
	Short: "Execute a WIQL query",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		top, _ := cmd.Flags().GetInt("top")
		items, err := c.WIQLAndFetch(project, args[0], nil, top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, items, func() {
			t := NewTable()
			t.Header("ID", "Type", "Title", "State")
			for _, wi := range items {
				wiType := fieldStr(wi.Fields, "System.WorkItemType")
				title := fieldStr(wi.Fields, "System.Title")
				state := fieldStr(wi.Fields, "System.State")
				t.Row(fmt.Sprintf("%d", wi.ID), wiType, Truncate(title, 50), state)
			}
			t.Flush()
			fmt.Printf("\n%d work items\n", len(items))
		})
	},
}

var searchQueryCmd = &cobra.Command{
	Use:   "query <id>",
	Short: "Run a saved query by ID or path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		top, _ := cmd.Flags().GetInt("top")
		items, err := c.RunQueryByID(project, args[0], top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, items, func() {
			t := NewTable()
			t.Header("ID", "Type", "Title", "State")
			for _, wi := range items {
				wiType := fieldStr(wi.Fields, "System.WorkItemType")
				title := fieldStr(wi.Fields, "System.Title")
				state := fieldStr(wi.Fields, "System.State")
				t.Row(fmt.Sprintf("%d", wi.ID), wiType, Truncate(title, 50), state)
			}
			t.Flush()
			fmt.Printf("\n%d work items\n", len(items))
		})
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(searchCodeCmd, searchWICmd, searchWiqlCmd, searchQueryCmd)

	searchCmd.PersistentFlags().StringP("project", "p", "", "Project name (optional, scopes search)")
	searchCmd.PersistentFlags().Int("top", 25, "Max results to return")
}
