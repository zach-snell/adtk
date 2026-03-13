package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zach-snell/adtk/internal/devops"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Azure DevOps",
	Long: `Set up credentials for accessing Azure DevOps.

You will need:
  1. Your Azure DevOps organization name (e.g., 'myorg' for dev.azure.com/myorg)
  2. A Personal Access Token (PAT) from https://dev.azure.com/{org}/_usersSettings/tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := devops.InteractiveLogin(); err != nil {
			fmt.Fprintf(os.Stderr, "auth failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		runStatus()
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and remove stored credentials",
	Run: func(cmd *cobra.Command, args []string) {
		runLogout()
	},
}

func init() {
	RootCmd.AddCommand(authCmd)
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(logoutCmd)
}

func runStatus() {
	// Check env vars first
	org := os.Getenv("AZURE_DEVOPS_ORG")
	pat := os.Getenv("AZURE_DEVOPS_PAT")

	if org != "" && pat != "" {
		fmt.Println("Authenticated via environment variables")
		fmt.Printf("  Organization: %s\n", org)
		fmt.Printf("  URL:          https://dev.azure.com/%s\n", org)
		return
	}

	creds, err := devops.LoadCredentials()
	if err != nil {
		fmt.Println("Not authenticated. Run: adtk auth")
		return
	}

	path, _ := devops.CredentialsPath()
	fmt.Println("Authenticated via stored credentials")
	KV("Organization", creds.Organization)
	KV("URL", fmt.Sprintf("https://dev.azure.com/%s", creds.Organization))
	if len(creds.PAT) > 8 {
		KVf("PAT", "%s...%s", creds.PAT[:4], creds.PAT[len(creds.PAT)-4:])
	} else {
		KV("PAT", "****")
	}
	KV("Saved", creds.SavedAt.Format("2006-01-02 15:04:05"))
	KV("File", path)
}

func runLogout() {
	if err := devops.RemoveCredentials(); err != nil {
		fmt.Fprintf(os.Stderr, "error removing credentials: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Logged out. Credentials removed.")
}
