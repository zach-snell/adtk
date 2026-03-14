package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var testPlansCmd = &cobra.Command{
	Use:     "test-plans",
	Aliases: []string{"tp"},
	Short:   "Manage Azure DevOps test plans, suites, and cases",
}

var testPlansListCmd = &cobra.Command{
	Use:   "list",
	Short: "List test plans in a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		plans, err := c.ListTestPlans(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(plans)
	},
}

var testPlansSuitesCmd = &cobra.Command{
	Use:   "suites <plan-id>",
	Short: "List test suites in a plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		planID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid plan ID: %s\n", args[0])
			os.Exit(1)
		}
		suites, err := c.ListTestSuites(project, planID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(suites)
	},
}

var testPlansCasesCmd = &cobra.Command{
	Use:   "cases <plan-id> <suite-id>",
	Short: "List test cases in a suite",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		planID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid plan ID: %s\n", args[0])
			os.Exit(1)
		}
		suiteID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid suite ID: %s\n", args[1])
			os.Exit(1)
		}
		cases, err := c.ListTestCases(project, planID, suiteID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(cases)
	},
}

var testPlansResultsCmd = &cobra.Command{
	Use:   "results <build-id>",
	Short: "Get test results for a build",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		buildID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid build ID: %s\n", args[0])
			os.Exit(1)
		}
		results, err := c.GetTestResultsForBuild(project, buildID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(results)
	},
}

func init() {
	RootCmd.AddCommand(testPlansCmd)
	testPlansCmd.AddCommand(testPlansListCmd, testPlansSuitesCmd, testPlansCasesCmd, testPlansResultsCmd)

	testPlansCmd.PersistentFlags().StringP("project", "p", "", "Project name (required)")
}
