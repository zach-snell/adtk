package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var wikiCmd = &cobra.Command{
	Use:   "wiki",
	Short: "Manage Azure DevOps wiki pages (markdown-native)",
}

var wikiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List wikis in a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		wikis, err := c.ListWikis(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, wikis, func() {
			t := NewTable()
			t.Header("ID", "Name", "Type")
			for _, w := range wikis {
				t.Row(w.ID, w.Name, w.Type)
			}
			t.Flush()
			fmt.Printf("\n%d wikis\n", len(wikis))
		})
	},
}

var wikiGetCmd = &cobra.Command{
	Use:   "get <wiki-id> <page-path>",
	Short: "Get a wiki page by path",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		page, err := c.GetWikiPage(project, args[0], args[1], true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, page, func() {
			KV("Path", page.Path)
			KVf("ID", "%d", page.ID)
			if page.Content != "" {
				fmt.Printf("\n%s\n", page.Content)
			}
		})
	},
}

func init() {
	RootCmd.AddCommand(wikiCmd)
	wikiCmd.AddCommand(wikiListCmd, wikiGetCmd)

	wikiCmd.PersistentFlags().StringP("project", "p", "", "Project name (required)")
}
