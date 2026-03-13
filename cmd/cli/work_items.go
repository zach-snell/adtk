package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var workItemsCmd = &cobra.Command{
	Use:     "work-items",
	Aliases: []string{"wi"},
	Short:   "Manage Azure DevOps work items",
}

var wiGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a work item by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid work item ID: %s\n", args[0])
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

func init() {
	RootCmd.AddCommand(workItemsCmd)
	workItemsCmd.AddCommand(wiGetCmd, wiListCmd, wiTypesCmd)

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

// fieldStrSlice extracts semicolon-delimited string as slice.
func fieldStrSlice(fields map[string]interface{}, key string) []string {
	v := fieldStr(fields, key)
	if v == "-" || v == "" {
		return nil
	}
	return strings.Split(v, "; ")
}
