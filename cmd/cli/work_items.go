package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zach-snell/adtk/internal/devops"
)

var workItemsCmd = &cobra.Command{
	Use:     "work-items",
	Aliases: []string{"wi"},
	Short:   "Manage Azure DevOps work items",
}

var wiGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a work item by ID (auto-detects from git branch if omitted)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		arg := ""
		if len(args) > 0 {
			arg = args[0]
		}
		id, err := devops.DetectWorkItemIDOrArg(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		project, _ := cmd.Flags().GetString("project")
		wi, err := c.GetWorkItem(project, id, "All")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, wi, func() {
			KVf("ID", "%d", wi.ID)
			KVf("Rev", "%d", wi.Rev)
			for key, val := range wi.Fields {
				KV(key, fmt.Sprintf("%v", val))
			}
		})
	},
}

var wiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List work items via WIQL query",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		query, _ := cmd.Flags().GetString("query")
		top, _ := cmd.Flags().GetInt("top")

		if query == "" {
			query = "SELECT [System.Id], [System.Title], [System.State], [System.AssignedTo] FROM WorkItems ORDER BY [System.ChangedDate] DESC"
		}

		items, err := c.WIQLAndFetch(project, query, nil, top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		PrintOrJSON(cmd, items, func() {
			t := NewTable()
			t.Header("ID", "Type", "Title", "State", "Assigned To")
			for _, wi := range items {
				wiType := fieldStr(wi.Fields, "System.WorkItemType")
				title := fieldStr(wi.Fields, "System.Title")
				state := fieldStr(wi.Fields, "System.State")
				assignedTo := fieldStr(wi.Fields, "System.AssignedTo")
				t.Row(fmt.Sprintf("%d", wi.ID), wiType, Truncate(title, 50), state, assignedTo)
			}
			t.Flush()
			fmt.Printf("\n%d work items\n", len(items))
		})
	},
}

var wiTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List available work item types for a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		types, err := c.GetWorkItemTypes(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, types, func() {
			t := NewTable()
			t.Header("Name", "Description")
			for _, wt := range types {
				t.Row(wt.Name, Truncate(wt.Description, 60))
			}
			t.Flush()
		})
	},
}

var wiMyCmd = &cobra.Command{
	Use:   "my",
	Short: "List my assigned work items",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		items, err := c.GetMyWorkItems(project, "", false, 50)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, items, func() {
			t := NewTable()
			t.Header("ID", "Type", "Title", "State", "Assigned To")
			for _, wi := range items {
				wiType := fieldStr(wi.Fields, "System.WorkItemType")
				title := fieldStr(wi.Fields, "System.Title")
				state := fieldStr(wi.Fields, "System.State")
				assignedTo := fieldStr(wi.Fields, "System.AssignedTo")
				t.Row(fmt.Sprintf("%d", wi.ID), wiType, Truncate(title, 50), state, assignedTo)
			}
			t.Flush()
			fmt.Printf("\n%d work items\n", len(items))
		})
	},
}

var wiCommentsCmd = &cobra.Command{
	Use:   "comments <id>",
	Short: "List comments on a work item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid work item ID: %s\n", args[0])
			os.Exit(1)
		}
		project, _ := cmd.Flags().GetString("project")
		comments, err := c.GetWorkItemComments(project, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, comments, func() {
			t := NewTable()
			t.Header("ID", "Author", "Date", "Text")
			for _, c := range comments.Comments {
				t.Row(
					fmt.Sprintf("%d", c.ID),
					c.CreatedBy.DisplayName,
					c.CreatedDate.Format("2006-01-02 15:04"),
					Truncate(c.Text, 60),
				)
			}
			t.Flush()
			fmt.Printf("\n%d comments\n", comments.TotalCount)
		})
	},
}

var wiMetricsCmd = &cobra.Command{
	Use:   "metrics <id>",
	Short: "Show lifecycle metrics for a work item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid work item ID: %s\n", args[0])
			os.Exit(1)
		}
		project, _ := cmd.Flags().GetString("project")
		metrics, err := c.ComputeWorkItemMetrics(project, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, metrics, func() {
			KV("Current Status", metrics.CurrentStatus)
			KV("Cycle Time", metrics.CycleTime.String())
			KV("Lead Time", metrics.LeadTime.String())
			fmt.Println("\n  Time in Status:")
			for status, d := range metrics.TimeInStatus {
				KV("  "+status, d.String())
			}
			if len(metrics.StatusTransitions) > 0 {
				fmt.Println("\n  Status Transitions:")
				t := NewTable()
				t.Header("From", "To", "At")
				for _, tr := range metrics.StatusTransitions {
					t.Row(tr.From, tr.To, tr.At.Format("2006-01-02 15:04"))
				}
				t.Flush()
			}
		})
	},
}

func init() {
	RootCmd.AddCommand(workItemsCmd)
	workItemsCmd.AddCommand(wiGetCmd, wiListCmd, wiTypesCmd, wiMyCmd, wiCommentsCmd, wiMetricsCmd)

	workItemsCmd.PersistentFlags().StringP("project", "p", "", "Project name")

	wiListCmd.Flags().StringP("query", "q", "", "WIQL query string")
	wiListCmd.Flags().Int("top", 25, "Max results to return")
}

// fieldStr extracts a string value from a work item fields map.
func fieldStr(fields map[string]interface{}, key string) string {
	v, ok := fields[key]
	if !ok {
		return "-"
	}
	switch val := v.(type) {
	case string:
		return val
	case map[string]interface{}:
		// AssignedTo is often a nested object with displayName
		if dn, ok := val["displayName"]; ok {
			return fmt.Sprintf("%v", dn)
		}
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
