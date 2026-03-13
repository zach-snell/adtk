package devops

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Credentials holds persisted Azure DevOps authentication data.
type Credentials struct {
	Organization string    `json:"organization"` // e.g., "myorg" (for dev.azure.com/myorg)
	PAT          string    `json:"pat"`
	SavedAt      time.Time `json:"saved_at"`
}

// CredentialsPath returns the path to the credentials file.
func CredentialsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home dir: %w", err)
	}
	return filepath.Join(home, ".config", "adtk", "credentials.json"), nil
}

// SaveCredentials persists credentials to disk.
func SaveCredentials(creds *Credentials) error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing credentials file: %w", err)
	}

	return nil
}

// LoadCredentials loads credentials from disk or environment variables.
// Priority: env vars > stored file.
func LoadCredentials() (*Credentials, error) {
	// 1. Check environment variables first
	org := os.Getenv("AZURE_DEVOPS_ORG")
	pat := os.Getenv("AZURE_DEVOPS_PAT")

	if org != "" && pat != "" {
		return &Credentials{
			Organization: org,
			PAT:          pat,
		}, nil
	}

	// 2. Fall back to stored credentials
	path, err := CredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading credentials file: %w (run 'adtk auth' to authenticate)", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("parsing credentials file: %w", err)
	}

	return &creds, nil
}

// RemoveCredentials deletes the stored credentials file.
func RemoveCredentials() error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// VerifyPAT tests whether the PAT is valid by calling the connection data endpoint.
func VerifyPAT(org, pat string) error {
	client := NewClient(org, pat)
	_, err := client.GetIdentity("/connectionData", nil)
	if err != nil {
		return fmt.Errorf("PAT verification failed: %w", err)
	}
	return nil
}

// InteractiveLogin prompts the user for Azure DevOps credentials and stores them.
func InteractiveLogin() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("Azure DevOps PAT Authentication")
	fmt.Println("================================")
	fmt.Println()
	fmt.Println("Create a Personal Access Token (PAT) at:")
	fmt.Println("  https://dev.azure.com/{org}/_usersSettings/tokens")
	fmt.Println()
	fmt.Println("Required scopes:")
	fmt.Println("  - Work Items: Read & Write")
	fmt.Println("  - Code: Read & Write")
	fmt.Println("  - Build: Read & Execute")
	fmt.Println("  - Release: Read")
	fmt.Println("  - Wiki: Read & Write")
	fmt.Println("  - Project and Team: Read")
	fmt.Println("  - Identity: Read")
	fmt.Println()
	fmt.Println("For read-only access, omit the Write scopes.")
	fmt.Println()

	fmt.Print("Organization name (e.g., 'myorg' for dev.azure.com/myorg): ")
	org, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading organization: %w", err)
	}
	org = strings.TrimSpace(org)
	if org == "" {
		return fmt.Errorf("organization is required")
	}
	// Strip full URL if provided
	org = strings.TrimPrefix(org, "https://dev.azure.com/")
	org = strings.TrimSuffix(org, "/")

	fmt.Print("Personal Access Token (PAT): ")
	pat, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading PAT: %w", err)
	}
	pat = strings.TrimSpace(pat)
	if pat == "" {
		return fmt.Errorf("PAT is required")
	}

	fmt.Println("\nVerifying credentials...")

	// Verify by hitting the connection data endpoint
	tmpClient := NewClient(org, pat)
	resp, err := tmpClient.do(http.MethodGet, fmt.Sprintf("https://%s/%s/_apis/connectionData?api-version=%s", HostMain, org, DefaultAPIVersion), nil, "")
	if err != nil {
		return fmt.Errorf("credential verification failed: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed: invalid PAT or organization")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	fmt.Println("Credentials verified successfully!")

	creds := &Credentials{
		Organization: org,
		PAT:          pat,
		SavedAt:      time.Now(),
	}

	if err := SaveCredentials(creds); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	path, _ := CredentialsPath()
	fmt.Printf("\nCredentials saved to: %s\n", path)
	fmt.Println("You can now use the Azure DevOps CLI and MCP server.")
	return nil
}
