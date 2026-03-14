package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Manage Azure DevOps Advanced Security alerts",
}

var securityAlertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "List security alerts for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		repo, _ := cmd.Flags().GetString("repo")
		states, _ := cmd.Flags().GetString("states")
		severities, _ := cmd.Flags().GetString("severities")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		if repo == "" {
			fmt.Fprintln(os.Stderr, "Error: --repo is required")
			os.Exit(1)
		}
		alerts, err := c.GetSecurityAlerts(project, repo, states, severities)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(alerts)
	},
}

var securityAlertCmd = &cobra.Command{
	Use:   "alert <alert-id>",
	Short: "Get details for a specific security alert",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		repo, _ := cmd.Flags().GetString("repo")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		if repo == "" {
			fmt.Fprintln(os.Stderr, "Error: --repo is required")
			os.Exit(1)
		}
		alertID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid alert ID: %s\n", args[0])
			os.Exit(1)
		}
		alert, err := c.GetSecurityAlertDetails(project, repo, alertID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(alert)
	},
}

func init() {
	RootCmd.AddCommand(securityCmd)
	securityCmd.AddCommand(securityAlertsCmd, securityAlertCmd)

	securityCmd.PersistentFlags().StringP("project", "p", "", "Project name (required)")
	securityCmd.PersistentFlags().StringP("repo", "r", "", "Repository name or ID (required)")
	securityAlertsCmd.Flags().String("states", "", "Filter by alert states")
	securityAlertsCmd.Flags().String("severities", "", "Filter by severities")
}
