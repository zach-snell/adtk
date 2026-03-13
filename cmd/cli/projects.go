package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zach-snell/adtk/internal/devops"
)

var projectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"proj"},
	Short:   "Manage Azure DevOps projects and teams",
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in the organization",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		result, err := devops.GetJSON[devops.ProjectList](c, "", "/projects", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, result, func() {
			t := NewTable()
			t.Header("Name", "State", "Visibility", "Last Updated")
			for _, p := range result.Value {
				t.Row(p.Name, p.State, p.Visibility, p.LastUpdateTime.Format("2006-01-02"))
			}
			t.Flush()
			fmt.Printf("\n%d projects\n", result.Count)
		})
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get <project>",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		path := fmt.Sprintf("/projects/%s", args[0])
		result, err := devops.GetJSON[devops.Project](c, "", path, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, result, func() {
			KV("ID", result.ID)
			KV("Name", result.Name)
			KV("State", result.State)
			KV("Visibility", result.Visibility)
			KV("Description", result.Description)
			KV("Last Updated", result.LastUpdateTime.Format("2006-01-02 15:04"))
		})
	},
}

var teamsListCmd = &cobra.Command{
	Use:   "teams <project>",
	Short: "List teams in a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		path := fmt.Sprintf("/projects/%s/teams", args[0])
		result, err := devops.GetJSON[devops.TeamList](c, "", path, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, result, func() {
			t := NewTable()
			t.Header("Name", "Description")
			for _, team := range result.Value {
				t.Row(team.Name, Truncate(team.Description, 50))
			}
			t.Flush()
			fmt.Printf("\n%d teams\n", result.Count)
		})
	},
}

func init() {
	RootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(projectsListCmd, projectsGetCmd, teamsListCmd)
}
