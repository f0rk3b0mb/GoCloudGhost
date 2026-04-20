// management/management.go
package management

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/f0rk3b0mb/GoCloudGhost/azure/auth"
	"github.com/f0rk3b0mb/GoCloudGhost/azure/models"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// EnumerationFlags holds all enumeration configuration
type EnumerationFlags struct {
	Token          string
	SubscriptionID string
	EnumSubs       bool
	EnumGroups     bool
	EnumRoles      bool
	EnumPolicies   bool
	EnumStorage    bool
	EnumKeyVaults  bool
}

// EnumerationTask represents a single enumeration function with its dependencies
type EnumerationTask struct {
	Name      string
	Requires  string // "token" or "subscription"
	FlagValue bool
	Fn        func(token, subscriptionID string) error
}

var MgmtCmd = &cobra.Command{
	Use:   "management",
	Short: "Enumerate Azure Management API resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse and validate flags
		flags, err := parseFlags(cmd)
		if err != nil {
			return err
		}

		// Load and validate credentials
		if err := loadCredentials(flags); err != nil {
			return err
		}

		// Validate that at least one enumeration option is selected
		if !hasAnyEnumerationFlag(flags) {
			return fmt.Errorf("no enumeration option selected. Use --help to see available options")
		}

		// Validate subscription requirement
		if err := validateSubscriptionRequirement(flags); err != nil {
			return err
		}

		// Execute enumeration tasks
		return executeEnumerationTasks(flags)
	},
}

// parseFlags extracts all command flags into a typed structure
func parseFlags(cmd *cobra.Command) (*EnumerationFlags, error) {
	flags := &EnumerationFlags{}

	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return nil, fmt.Errorf("failed to parse token flag: %w", err)
	}
	flags.Token = token

	subscriptionID, err := cmd.Flags().GetString("subscription")
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscription flag: %w", err)
	}
	flags.SubscriptionID = subscriptionID

	flags.EnumSubs, _ = cmd.Flags().GetBool("subscriptions")
	flags.EnumGroups, _ = cmd.Flags().GetBool("groups")
	flags.EnumRoles, _ = cmd.Flags().GetBool("roles")
	flags.EnumPolicies, _ = cmd.Flags().GetBool("policies")
	flags.EnumStorage, _ = cmd.Flags().GetBool("storage")
	flags.EnumKeyVaults, _ = cmd.Flags().GetBool("keyvaults")

	return flags, nil
}

// loadCredentials loads token and subscription from CLI or environment
func loadCredentials(flags *EnumerationFlags) error {
	// Load token: CLI flag -> ACCESS_TOKEN env -> .env file
	token, err := loadTokenFromMultipleSources(flags.Token)
	if err != nil {
		return err
	}
	flags.Token = token

	// Load subscription: CLI flag -> AZURE_SUBSCRIPTION_ID env
	if flags.SubscriptionID == "" {
		flags.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}

	return nil
}

// loadTokenFromMultipleSources attempts to load token from multiple sources in priority order
func loadTokenFromMultipleSources(cliToken string) (string, error) {
	// Priority 1: CLI token if provided
	if cliToken != "" {
		return cliToken, nil
	}

	// Priority 2: ACCESS_TOKEN environment variable
	if envToken := os.Getenv("ACCESS_TOKEN"); envToken != "" {
		return envToken, nil
	}

	// Priority 3: Try loading from .env file
	_ = godotenv.Load() // Ignore error if .env doesn't exist
	if envToken := os.Getenv("ACCESS_TOKEN"); envToken != "" {
		return envToken, nil
	}

	return "", fmt.Errorf("Azure access token is required. Provide it via:\n  1. --token flag\n  2. ACCESS_TOKEN environment variable\n  3. ACCESS_TOKEN in .env file")
}

// hasAnyEnumerationFlag checks if at least one enumeration option is enabled
func hasAnyEnumerationFlag(flags *EnumerationFlags) bool {
	return flags.EnumSubs || flags.EnumGroups || flags.EnumRoles ||
		flags.EnumPolicies || flags.EnumStorage || flags.EnumKeyVaults
}

// validateSubscriptionRequirement validates that subscription is provided when needed
func validateSubscriptionRequirement(flags *EnumerationFlags) error {
	subscriptionRequired := flags.EnumGroups || flags.EnumRoles ||
		flags.EnumStorage || flags.EnumKeyVaults

	if subscriptionRequired && flags.SubscriptionID == "" {
		return fmt.Errorf("--subscription is required for: groups, roles, storage, and keyvaults\nProvide via:\n  1. --subscription flag\n  2. AZURE_SUBSCRIPTION_ID environment variable \n 3. run with flag --subscriptions to enumerate subscriptions first")
	}

	return nil
}

// executeEnumerationTasks builds and executes all enabled enumeration tasks
func executeEnumerationTasks(flags *EnumerationFlags) error {
	tasks := buildEnumerationTasks(flags)

	for _, task := range tasks {
		if !task.FlagValue {
			continue
		}

		if err := task.Fn(flags.Token, flags.SubscriptionID); err != nil {
			return fmt.Errorf("failed to enumerate %s: %w", task.Name, err)
		}
	}

	return nil
}

// buildEnumerationTasks creates the list of tasks to execute
func buildEnumerationTasks(flags *EnumerationFlags) []EnumerationTask {
	return []EnumerationTask{
		{
			Name:      "subscriptions",
			Requires:  "token",
			FlagValue: flags.EnumSubs,
			Fn: func(token, _ string) error {
				fmt.Println("\n=== SUBSCRIPTIONS ===")
				return auth.EnumerateSubscriptions(token)
			},
		},
		{
			Name:      "resource groups",
			Requires:  "subscription",
			FlagValue: flags.EnumGroups,
			Fn: func(token, subID string) error {
				return enumerateResourceGroups(token, subID)
			},
		},
		{
			Name:      "role assignments",
			Requires:  "subscription",
			FlagValue: flags.EnumRoles,
			Fn: func(token, subID string) error {
				return enumerateRoleAssignments(token, subID)
			},
		},
		{
			Name:      "policies",
			Requires:  "token",
			FlagValue: flags.EnumPolicies,
			Fn: func(token, _ string) error {
				return enumeratePolicyDefinitions(token)
			},
		},
		{
			Name:      "storage accounts",
			Requires:  "subscription",
			FlagValue: flags.EnumStorage,
			Fn: func(token, subID string) error {
				return enumerateStorageAccounts(token, subID)
			},
		},
		{
			Name:      "key vaults",
			Requires:  "subscription",
			FlagValue: flags.EnumKeyVaults,
			Fn: func(token, subID string) error {
				return enumerateKeyVaults(token, subID)
			},
		},
	}
}

func init() {
	MgmtCmd.Flags().String("token", "", "Azure access token")
	MgmtCmd.Flags().String("subscription", "", "Azure subscription ID")
	MgmtCmd.Flags().Bool("subscriptions", false, "Enumerate subscriptions")
	MgmtCmd.Flags().Bool("groups", false, "Enumerate resource groups")
	MgmtCmd.Flags().Bool("roles", false, "Enumerate role assignments")
	MgmtCmd.Flags().Bool("policies", false, "Enumerate policy definitions")
	MgmtCmd.Flags().Bool("storage", false, "Enumerate storage accounts")
	MgmtCmd.Flags().Bool("keyvaults", false, "Enumerate key vaults")
}

func isDangerousRole(role models.RoleDefinition) bool {
	roleName := role.Properties.RoleName

	if roleName == "Owner" || roleName == "User Access Administrator" {
		return true
	}

	for _, perm := range role.Properties.Permissions {
		for _, action := range perm.Actions {
			if action == "*" ||
				action == "Microsoft.Authorization/*" ||
				action == "Microsoft.Authorization/roleAssignments/write" {
				return true
			}
		}
	}

	return false
}

func enumerateRoleDefinitions(token, subscriptionID string) (map[string]models.RoleDefinition, error) {
	ctx := context.Background()

	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions?api-version=2022-04-01",
		subscriptionID,
	)

	var result models.RoleDefinitionsResponse
	if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, url, &result); err != nil {
		return nil, err
	}

	fmt.Println("\n=== ROLE DEFINITIONS ===")

	roleMap := make(map[string]models.RoleDefinition)

	for _, role := range result.Value {
		roleMap[role.ID] = role

		level := "[INFO]"
		if isDangerousRole(role) {
			level = "[CRITICAL]"
		}

		fmt.Printf("%s Role: %-30s AssignableScopes: %d\n",
			level,
			role.Properties.RoleName,
			len(role.Properties.AssignableScopes),
		)
	}

	return roleMap, nil
}

func enumerateRoleAssignments(token, subscriptionID string) error {
	ctx := context.Background()

	roleMap, err := enumerateRoleDefinitions(token, subscriptionID)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/providers/Microsoft.Authorization/roleAssignments?api-version=2022-04-01",
		subscriptionID,
	)

	var result models.RoleAssignmentsResponse
	if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, url, &result); err != nil {
		return err
	}

	fmt.Println("\n=== ROLE ASSIGNMENTS ===")

	for _, assignment := range result.Value {
		role, exists := roleMap[assignment.Properties.RoleDefinitionID]
		if !exists {
			fmt.Printf("[WARN] Unknown role for principal %s\n", assignment.Properties.PrincipalID)
			continue
		}

		level := "[INFO]"
		if isDangerousRole(role) {
			level = "[CRITICAL]"
		}

		fmt.Printf("%s Principal: %-36s Role: %-30s Scope: %s\n",
			level,
			assignment.Properties.PrincipalID,
			role.Properties.RoleName,
			assignment.Properties.Scope,
		)
	}

	return nil
}

func enumeratePolicyDefinitions(token string) error {
	ctx := context.Background()

	url := "https://management.azure.com/providers/Microsoft.Authorization/policyDefinitions?api-version=2021-06-01"

	var result map[string]interface{}
	if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, url, &result); err != nil {
		return err
	}

	fmt.Println("\n=== POLICY DEFINITIONS ===")

	for _, p := range result["value"].([]interface{}) {
		policy := p.(map[string]interface{})
		props := policy["properties"].(map[string]interface{})

		fmt.Printf("[INFO] Policy: %-40s Type: %s\n",
			props["displayName"],
			props["policyType"],
		)
	}

	return nil
}

func enumerateResourceGroups(token, subscriptionID string) error {
	ctx := context.Background()

	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/resourcegroups?api-version=2021-04-01",
		subscriptionID,
	)

	var result map[string]interface{}
	if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, url, &result); err != nil {
		return err
	}

	fmt.Println("\n=== RESOURCE GROUPS ===")

	for _, rg := range result["value"].([]interface{}) {
		group := rg.(map[string]interface{})
		fmt.Printf("[INFO] Resource Group: %-25s Location: %s\n",
			group["name"],
			group["location"],
		)
	}

	return nil
}

func enumerateStorageAccounts(token, subscriptionID string) error {
	ctx := context.Background()

	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/providers/Microsoft.Storage/storageAccounts?api-version=2022-09-01",
		subscriptionID,
	)

	var result map[string]interface{}
	if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, url, &result); err != nil {
		return err
	}

	fmt.Println("\n=== STORAGE ACCOUNTS ===")

	accounts := result["value"].([]interface{})
	for _, acc := range accounts {
		accMap := acc.(map[string]interface{})
		name := accMap["name"].(string)
		id := accMap["id"].(string)
		resourceGroup := extractResourceGroupFromID(id)

		fmt.Printf("[INFO] Storage Account: %-25s Resource Group: %s\n", name, resourceGroup)

		keyURL := fmt.Sprintf(
			"https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/listKeys?api-version=2022-09-01",
			subscriptionID,
			resourceGroup,
			name,
		)

		var keyResult map[string]interface{}
		if err := makeAuthenticatedRequest(ctx, token, http.MethodPost, keyURL, &keyResult); err != nil {
			fmt.Printf("[WARN] Key request failed for %s: %v\n", name, err)
			continue
		}

		fmt.Printf("[CRITICAL] Keys accessible for %s: %v\n", name, keyResult)
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

func getResponseArray(result map[string]interface{}, fieldName string) ([]interface{}, error) {
	raw, ok := result[fieldName]
	if !ok || raw == nil {
		return nil, nil
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format for %s: %T", fieldName, raw)
	}

	return arr, nil
}

func enumerateKeyVaults(token, subscriptionID string) error {
	ctx := context.Background()
	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/providers/Microsoft.KeyVault/vaults?api-version=2021-10-01",
		subscriptionID,
	)

	var result map[string]interface{}
	if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, url, &result); err != nil {
		return err
	}

	fmt.Println("\n=== KEY VAULTS ===")

	vaults, err := getResponseArray(result, "value")
	if err != nil {
		return err
	}
	if len(vaults) == 0 {
		fmt.Println("[INFO] No key vaults found.")
		return nil
	}

	for _, kv := range vaults {
		kvMap, ok := kv.(map[string]interface{})
		if !ok {
			fmt.Printf("[WARN] Skipping unexpected key vault entry: %T\n", kv)
			continue
		}

		name, _ := kvMap["name"].(string)
		id, _ := kvMap["id"].(string)
		resourceGroup := extractResourceGroupFromID(id)

		fmt.Printf("[INFO] Key Vault: %-25s Resource Group: %s\n", name, resourceGroup)

		secretURL := fmt.Sprintf(
			"https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/secrets?api-version=2021-10-01",
			subscriptionID,
			resourceGroup,
			name,
		)

		var secretResult map[string]interface{}
		if err := makeAuthenticatedRequest(ctx, token, http.MethodGet, secretURL, &secretResult); err != nil {
			fmt.Printf("[WARN] Secret request failed for %s: %v\n", name, err)
			continue
		}

		fmt.Printf("[CRITICAL] Secrets accessible for %s: %v\n", name, secretResult)
	}

	return nil
}

// makeAuthenticatedRequest performs an authenticated HTTP request and decodes JSON response
func makeAuthenticatedRequest(ctx context.Context, token, method, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
