package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zach-snell/adtk/internal/devops"
	"github.com/zach-snell/adtk/internal/version"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:     "adtk",
	Version: version.Version,
	Short:   "A unified CLI and MCP server for Azure DevOps",
	Long: `adtk is a complete command-line interface and Model Context Protocol
server for Azure DevOps.

It allows you to manage work items, repositories, pull requests, pipelines,
and boards directly from your terminal, or expose these capabilities to your
AI agents via the MCP protocol.

Try running 'adtk auth' to get started!`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().Bool("json", false, "Output raw JSON instead of formatted tables")
}

// getClient creates an Azure DevOps API client from stored credentials or env vars.
func getClient() *devops.Client {
	creds, err := devops.LoadCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'adtk auth' to authenticate, or set AZURE_DEVOPS_ORG and AZURE_DEVOPS_PAT env vars.\n")
		os.Exit(1)
	}
	return devops.NewClient(creds.Organization, creds.PAT)
}
