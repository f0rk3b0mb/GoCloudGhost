package gcplist

// TODO
// add impersonation support for GCP

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var TokenCmd = &cobra.Command{
	Use:   "impersonate",
	Short: "Check for token impersonation permissions",
	Long:  `This command allows you to check for user impersonation permissions for your token`,
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		projectID, _ := cmd.Flags().GetString("project-id")
		if token == "" || projectID == "" {
			fmt.Println("Error: both --token and --project-id are required")
			return
		}
		ListServiceAccounts(token, projectID)
	},
}

func init() {
	TokenCmd.Flags().String("token", "", "GCP OAuth2 token")
	TokenCmd.Flags().String("project-id", "", "GCP Project ID")
	TokenCmd.MarkFlagRequired("token")
	TokenCmd.MarkFlagRequired("project-id")
}

func ListServiceAccounts(token string, projectID string) {
	iamURL := fmt.Sprintf("https://iam.googleapis.com/v1/projects/%s/serviceAccounts", projectID)

	req, err := http.NewRequest("GET", iamURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[-] Failed to list service accounts (status: %d). Token may be invalid.\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	accounts, ok := jsonData["accounts"].([]interface{})
	if !ok || len(accounts) == 0 {
		// fallback to correct field
		accounts, ok = jsonData["accounts"].([]interface{})
		if !ok || len(accounts) == 0 {
			fmt.Println("[-] No service accounts found or response format invalid")
			return
		}
	}

	fmt.Println("\n[*] Listing Service Accounts:")
	for _, acct := range accounts {
		acctMap, ok := acct.(map[string]interface{})
		if !ok {
			continue
		}
		if email, ok := acctMap["email"].(string); ok {
			fmt.Println(" -", email)
		}
	}

	fmt.Println("\n[*] Checking for impersonation permissions...")

	for _, acct := range accounts {
		acctMap, ok := acct.(map[string]interface{})
		if !ok {
			continue
		}
		email, ok := acctMap["email"].(string)
		if !ok {
			continue
		}
		fmt.Println("[*] Trying to impersonate:", email)
		TokenImpersonate(email, token, projectID)
	}
}

func TokenImpersonate(saEmail, token, projectID string) *string {
	url := fmt.Sprintf("https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/%s:generateAccessToken", saEmail)

	payload := map[string]interface{}{
		"scope": []string{"https://www.googleapis.com/auth/cloud-platform"},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[-] Failed to create request:", err)
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[-] Request failed for %s: %v\n", saEmail, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			AccessToken string `json:"accessToken"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("[-] JSON decode failed for %s: %v\n", saEmail, err)
			return nil
		}
		fmt.Printf("[+] Impersonation SUCCESS for %s\n", saEmail)
		fmt.Printf("[+] Access Token: %s\n", result.AccessToken)
		return &result.AccessToken
	}

	fmt.Printf("[-] Impersonation FAILED for %s: %d\n", saEmail, resp.StatusCode)
	return nil
}
