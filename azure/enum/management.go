// management/management.go
package management

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var MgmtCmd = &cobra.Command{
	Use:   "management",
	Short: "Enumerate Azure Management API resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		subscriptionID, _ := cmd.Flags().GetString("subscription")
		enumSubs, _ := cmd.Flags().GetBool("subscriptions")
		enumGroups, _ := cmd.Flags().GetBool("groups")
		enumRoles, _ := cmd.Flags().GetBool("roles")
		enumPolicies, _ := cmd.Flags().GetBool("policies")
		enumStorage, _ := cmd.Flags().GetBool("storage")

		if token == "" {
			return fmt.Errorf("--token is required")
		}

		if enumSubs {
			fmt.Println("Enumerating subscriptions...")
			if err := enumerateSubscriptions(token); err != nil {
				return err
			}
		}

		if subscriptionID != "" {
			if enumGroups {
				fmt.Println("Enumerating resource groups...")
				if err := enumerateResourceGroups(token, subscriptionID); err != nil {
					return err
				}
			}
			if enumRoles {
				fmt.Println("Enumerating role assignments...")
				if err := enumerateRoleAssignments(token, subscriptionID); err != nil {
					return err
				}
			}
			if enumStorage {
				fmt.Println("Enumerating storage accounts...")
				if err := enumerateStorageAccounts(token, subscriptionID); err != nil {
					return err
				}
			}
		} else {
			if enumGroups || enumRoles || enumStorage {
				return fmt.Errorf("--subscription is required for --groups and --roles and --storage")
			}
		}

		if enumPolicies {
			fmt.Println("Enumerating policy definitions...")
			if err := enumeratePolicyDefinitions(token); err != nil {
				return err
			}
		}

		if !enumSubs && !enumGroups && !enumRoles && !enumPolicies && !enumStorage {
			return fmt.Errorf("No enumeration option selected. Use --subscriptions, --groups, --roles, or --policies")
		}

		return nil
	},
}

func init() {
	MgmtCmd.Flags().String("token", "", "Azure access token (required)")
	MgmtCmd.Flags().String("subscription", "", "Azure subscription ID (optional)")
	MgmtCmd.Flags().Bool("subscriptions", false, "Enumerate subscriptions")
	MgmtCmd.Flags().Bool("groups", false, "Enumerate resource groups (requires --subscription)")
	MgmtCmd.Flags().Bool("roles", false, "Enumerate role assignments (requires --subscription)")
	MgmtCmd.Flags().Bool("policies", false, "Enumerate policy definitions")
	MgmtCmd.Flags().Bool("storage", false, "Enumerate Storage accounts")
	MgmtCmd.MarkFlagRequired("token")
}

func enumerateSubscriptions(token string) error {
	ctx := context.Background()
	url := "https://management.azure.com/subscriptions?api-version=2020-01-01"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(pretty))
	return nil
}

func enumerateStorageAccounts(token, subscriptionID string) error {
	ctx := context.Background()
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/providers/Microsoft.Storage/storageAccounts?api-version=2022-09-01", subscriptionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	accounts := result["value"].([]interface{})
	for _, acc := range accounts {
		accMap := acc.(map[string]interface{})
		name := accMap["name"].(string)
		id := accMap["id"].(string)
		resourceGroup := extractResourceGroupFromID(id)
		fmt.Printf("\nStorage Account: %s\n", name)

		keyURL := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/listKeys?api-version=2022-09-01", subscriptionID, resourceGroup, name)
		keyReq, err := http.NewRequestWithContext(ctx, "POST", keyURL, nil)
		if err != nil {
			fmt.Printf("  [!] Failed to build key request for %s: %v\n", name, err)
			continue
		}
		keyReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		keyResp, err := http.DefaultClient.Do(keyReq)
		if err != nil {
			fmt.Printf("  [!] Key request failed for %s: %v\n", name, err)
			continue
		}
		defer keyResp.Body.Close()

		if keyResp.StatusCode != http.StatusOK {
			fmt.Printf("  [!] Access denied to list keys for %s: %s\n", name, keyResp.Status)
			continue
		}

		body, _ := ioutil.ReadAll(keyResp.Body)
		fmt.Printf("  Keys: %s\n", string(body))
	}

	return nil
}

func extractResourceGroupFromID(id string) string {
	parts := strings.Split(id, "/")
	for i, part := range parts {
		if part == "resourceGroups" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func enumerateResourceGroups(token, subscriptionID string) error {
	ctx := context.Background()
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourcegroups?api-version=2021-04-01", subscriptionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println("Resource Groups:")
	fmt.Println(string(pretty))
	return nil
}

func enumerateRoleAssignments(token, subscriptionID string) error {
	ctx := context.Background()
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/providers/Microsoft.Authorization/roleAssignments?api-version=2022-04-01", subscriptionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println("Role Assignments:")
	fmt.Println(string(pretty))
	return nil
}

func enumeratePolicyDefinitions(token string) error {
	ctx := context.Background()
	url := "https://management.azure.com/providers/Microsoft.Authorization/policyDefinitions?api-version=2021-06-01"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println("Policy Definitions:")
	fmt.Println(string(pretty))
	return nil
}
