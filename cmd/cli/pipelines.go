package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var pipelinesCmd = &cobra.Command{
	Use:     "pipelines",
	Aliases: []string{"pipe"},
	Short:   "Manage Azure DevOps CI/CD pipelines",
}

var pipelinesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pipeline definitions in a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		top, _ := cmd.Flags().GetInt("top")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		pipelines, err := c.ListPipelines(project, top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, pipelines, func() {
			t := NewTable()
			t.Header("ID", "Name", "Folder")
			for _, p := range pipelines {
				t.Row(fmt.Sprintf("%d", p.ID), p.Name, p.Folder)
			}
			t.Flush()
			fmt.Printf("\n%d pipelines\n", len(pipelines))
		})
	},
}

var pipelinesGetCmd = &cobra.Command{
	Use:   "get <pipeline-id>",
	Short: "Get pipeline details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid pipeline ID: %s\n", args[0])
			os.Exit(1)
		}
		pipeline, err := c.GetPipeline(project, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, pipeline, func() {
			KVf("ID", "%d", pipeline.ID)
			KV("Name", pipeline.Name)
			KV("Folder", pipeline.Folder)
		})
	},
}

var pipelinesRunsCmd = &cobra.Command{
	Use:   "runs <pipeline-id>",
	Short: "List runs for a pipeline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		top, _ := cmd.Flags().GetInt("top")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid pipeline ID: %s\n", args[0])
			os.Exit(1)
		}
		runs, err := c.ListPipelineRuns(project, id, top)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, runs, func() {
			t := NewTable()
			t.Header("ID", "Name", "State", "Result", "Created")
			for _, r := range runs {
				t.Row(
					fmt.Sprintf("%d", r.ID),
					r.Name,
					r.State,
					r.Result,
					r.CreatedDate.Format("2006-01-02 15:04"),
				)
			}
			t.Flush()
			fmt.Printf("\n%d runs\n", len(runs))
		})
	},
}

var pipelinesDefinitionsCmd = &cobra.Command{
	Use:   "definitions",
	Short: "List build definitions",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		defs, err := c.ListBuildDefinitions(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(defs)
	},
}

var pipelinesVarGroupsCmd = &cobra.Command{
	Use:   "var-groups",
	Short: "List variable groups in a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		groups, err := c.ListVariableGroups(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, groups, func() {
			t := NewTable()
			t.Header("ID", "Name", "Description")
			for _, g := range groups {
				id := fmt.Sprintf("%v", g["id"])
				name, _ := g["name"].(string)
				desc, _ := g["description"].(string)
				t.Row(id, name, Truncate(desc, 50))
			}
			t.Flush()
			fmt.Printf("\n%d variable groups\n", len(groups))
		})
	},
}

var pipelinesVarGroupGetCmd = &cobra.Command{
	Use:   "var-group <id>",
	Short: "Get a variable group by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid group ID: %s\n", args[0])
			os.Exit(1)
		}
		group, err := c.GetVariableGroup(project, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintJSON(group)
	},
}

var pipelinesEnvironmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "List deployment environments in a project",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Fprintln(os.Stderr, "Error: --project is required")
			os.Exit(1)
		}
		envs, err := c.ListEnvironments(project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintOrJSON(cmd, envs, func() {
			t := NewTable()
			t.Header("ID", "Name")
			for _, e := range envs {
				id := fmt.Sprintf("%v", e["id"])
				name, _ := e["name"].(string)
				t.Row(id, name)
			}
			t.Flush()
			fmt.Printf("\n%d environments\n", len(envs))
		})
	},
}

func init() {
	RootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.AddCommand(pipelinesListCmd, pipelinesGetCmd, pipelinesRunsCmd, pipelinesDefinitionsCmd, pipelinesVarGroupsCmd, pipelinesVarGroupGetCmd, pipelinesEnvironmentsCmd)

	pipelinesCmd.PersistentFlags().StringP("project", "p", "", "Project name (required)")
	pipelinesCmd.PersistentFlags().Int("top", 25, "Max results to return")
}
