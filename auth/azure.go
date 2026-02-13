package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Azure",
	Long:  "Authenticate with Azure using a service principal",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		tenantID, _ := cmd.Flags().GetString("tenant-id")

		if clientID == "" || clientSecret == "" || tenantID == "" {
			return fmt.Errorf("client-id, client-secret, and tenant-id are required")
		}

		return Authenticate(clientID, clientSecret, tenantID)
	},
}

func init() {
	AuthCmd.Flags().String("client-id", "", "Azure Client ID")
	AuthCmd.Flags().String("client-secret", "", "Azure Client Secret")
	AuthCmd.Flags().String("tenant-id", "", "Azure Tenant ID")

	_ = AuthCmd.MarkFlagRequired("client-id")
	_ = AuthCmd.MarkFlagRequired("client-secret")
	_ = AuthCmd.MarkFlagRequired("tenant-id")
}

/* =======================
   ENV HELPERS
======================= */

func loadEnv(path string) (map[string]string, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return map[string]string{}, nil // file may not exist
	}
	return env, nil
}

func saveEnv(path string, env map[string]string) error {
	return godotenv.Write(env, path)
}

/* =======================
   AUTH
======================= */

func Authenticate(clientID, clientSecret, tenantID string) error {
	scope := "https://management.azure.com/.default"
	tokenURL := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/token",
		tenantID,
	)

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("grant_type", "client_credentials")
	form.Set("scope", scope)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request failed: %s\n%s", resp.Status, body)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("no access_token in response")
	}

	env, err := loadEnv("./.env")
	if err != nil {
		return err
	}

	env["ACCESS_TOKEN"] = tokenResp.AccessToken
	env["AZURE_TENANT_ID"] = tenantID

	if err := saveEnv("./.env", env); err != nil {
		return err
	}

	fmt.Println("Access token acquired and stored")

	return enumerateSubscriptions(tokenResp.AccessToken)
}

/* =======================
   SUBSCRIPTIONS
======================= */

func enumerateSubscriptions(token string) error {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://management.azure.com/subscriptions?api-version=2020-01-01",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("subscription request failed: %s\n%s", resp.Status, body)
	}

	var result struct {
		Value []struct {
			ID   string `json:"subscriptionId"`
			Name string `json:"displayName"`
		} `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if len(result.Value) == 0 {
		return fmt.Errorf("no subscriptions found")
	}

	env, err := loadEnv("./.env")
	if err != nil {
		return err
	}

	env["AZURE_SUBSCRIPTION_ID"] = result.Value[0].ID
	env["AZURE_SUBSCRIPTION_NAME"] = result.Value[0].Name

	if err := saveEnv("./.env", env); err != nil {
		return err
	}

	fmt.Printf("Subscription selected: %s (%s)\n",
		result.Value[0].Name,
		result.Value[0].ID,
	)

	return nil
}
