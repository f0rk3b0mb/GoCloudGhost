package gcpenum

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var EnumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Enumerate GCP resources",
	Long:  `This command allows you to enumerate GCP resources and their permissions.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		project_id, _ := cmd.Flags().GetString("project-id")
		if token == "" {
			fmt.Println("Error: --token is required")
			return
		}
		if project_id == "" {
			fmt.Println("Error: --project-id is required")
			return
		}
		run(token, project_id)
	},
}

func init() {
	EnumCmd.Flags().String("token", "", "GCP OAuth2 token")
	EnumCmd.Flags().String("project-id", "", "GCP Project ID")
	EnumCmd.MarkFlagRequired("token")
	EnumCmd.MarkFlagRequired("project-id")
}

func GetPermissions(name string, url string, token string) []string {
	ctx := context.Background()
	HEADERS := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Content-Type": "application/json"}

	// make get request to url with header set

	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []string{"Error creating request: " + err.Error()}
	}

	// set headers
	for key, value := range HEADERS {
		req.Header.Set(key, value)
	}

	// make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return []string{"Error making request: " + err.Error()}
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode == http.StatusOK {
		return []string{"✅ " + name + ": OK - Permissions retrieved successfully."}
	} else if resp.StatusCode == http.StatusUnauthorized {
		return []string{"❌ " + name + ": Forbidden (needs permission) or token is expired."}
	} else if resp.StatusCode == http.StatusForbidden {
		return []string{"⚠️ " + name + ": Not found (maybe API not enabled)"}
	} else {
		return []string{fmt.Sprintf("Error: received unexpected status code %d", resp.StatusCode)}
	}

}

func GetScopePermissions(name string, url string, token string) []string {
	ctx := context.Background()
	HEADERS := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}

	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []string{"Error creating request: " + err.Error()}
	}

	for key, value := range HEADERS {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return []string{"Error making request: " + err.Error()}
	}
	defer resp.Body.Close()

	if true {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []string{"Error reading response body: " + err.Error()}
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
			return []string{"Error parsing JSON: " + err.Error()}
		}

		formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return []string{"Error formatting JSON: " + err.Error()}
		}

		// You can also return both the message and the JSON
		return []string{
			"✅ OK - Token info retrieved successfully.",
			string(formattedJSON),
		}

	} else {
		return []string{fmt.Sprintf("Error: %s returned unexpected status code %d", name, resp.StatusCode)}
	}

	//else if resp.StatusCode == http.StatusUnauthorized {
	//	return []string{"❌ " + name + ": Unauthorized (invalid or expired token)."}
	//} else if resp.StatusCode == http.StatusForbidden {
	//	return []string{"⚠️ " + name + ": Forbidden (missing permissions or API not enabled)."}
	//} else {
	//	return []string{fmt.Sprintf("Error: %s returned unexpected status code %d", name, resp.StatusCode)}
	//}
}

func run(token, projectID string) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", token)
	name := "GCP OAuth2 Token Info"
	permissions := GetPermissions(name, url, token)
	for _, perm := range permissions {
		fmt.Println(perm)
	}

	url = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/aggregated/instances", projectID)
	name = "null"
	permissions = GetScopePermissions(name, url, token)
	for _, perm := range permissions {
		fmt.Println(perm)
	}

	url = fmt.Sprintf("https://storage.googleapis.com/storage/v1/b?project=%s", projectID)
	name = "GCP Storage Buckets"
	permissions = GetPermissions(name, url, token)
	for _, perm := range permissions {
		fmt.Println(perm)
	}

	url = fmt.Sprintf("https://cloudfunctions.googleapis.com/v1/projects/%s/locations/-/functions", projectID)
	name = "GCP Cloud Functions"
	permissions = GetPermissions(name, url, token)
	for _, perm := range permissions {
		fmt.Println(perm)
	}

	url = fmt.Sprintf("https://run.googleapis.com/v1/projects/%s/locations/-/services", projectID)
	name = "GCP Cloud Run Services"
	permissions = GetPermissions(name, url, token)
	for _, perm := range permissions {
		fmt.Println(perm)
	}

	url = fmt.Sprintf("https://bigquery.googleapis.com/bigquery/v2/projects/%s/datasets", projectID)
	name = "GCP BigQuery Datasets"
	permissions = GetPermissions(name, url, token)
	for _, perm := range permissions {
		fmt.Println(perm)
	}

}
