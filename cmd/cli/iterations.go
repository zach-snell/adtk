package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var iterationsCmd = &cobra.Command{
	Use:     "iterations",
	Aliases: []string{"iter", "sprints"},
	Short:   "Manage Azure DevOps iterations (sprints)",
}

var iterationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List iterations for a project/team",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		team, _ := cmd.Flags().GetString("team")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		iterations, err := c.ListIterations(project, team)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, iterations, func() {
			t := NewTable()
			t.Header("Name", "Path", "Time Frame", "Start", "End")
			for _, iter := range iterations {
				tf := "-"
				start := "-"
				end := "-"
				if iter.Attributes != nil {
					tf = iter.Attributes.TimeFrame
					if iter.Attributes.StartDate != nil {
						start = iter.Attributes.StartDate.Format("2006-01-02")
					}
					if iter.Attributes.FinishDate != nil {
						end = iter.Attributes.FinishDate.Format("2006-01-02")
					}
				}
				t.Row(iter.Name, iter.Path, tf, start, end)
			}
			t.Flush()
			fmt.Printf("\n%d iterations\n", len(iterations))
		})
	},
}

var iterationsCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current active iteration",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		team, _ := cmd.Flags().GetString("team")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		iter, err := c.GetCurrentIteration(project, team)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, iter, func() {
			KV("Name", iter.Name)
			KV("Path", iter.Path)
			if iter.Attributes != nil {
				if iter.Attributes.StartDate != nil {
					KV("Start", iter.Attributes.StartDate.Format("2006-01-02"))
				}
				if iter.Attributes.FinishDate != nil {
					KV("End", iter.Attributes.FinishDate.Format("2006-01-02"))
				}
			}
		})
	},
}

var iterationsSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Get team iteration settings",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		team, _ := cmd.Flags().GetString("team")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		settings, err := c.GetTeamSettings(project, team)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(settings)
	},
}

func init() {
	RootCmd.AddCommand(iterationsCmd)
	iterationsCmd.AddCommand(iterationsListCmd, iterationsCurrentCmd, iterationsSettingsCmd)

	iterationsCmd.PersistentFlags().StringP("project", "p", "", "Project name (required)")
	iterationsCmd.PersistentFlags().StringP("team", "t", "", "Team name (optional)")
}
