package gcplist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var BucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "List GCP Storage Buckets",
	Long:  `This command allows you to list GCP Storage Buckets.`,
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
		ListBuckets(token, project_id)
	},
}

func init() {
	BucketCmd.Flags().String("token", "", "GCP OAuth2 token")
	BucketCmd.Flags().String("project-id", "", "GCP Project ID")
	BucketCmd.MarkFlagRequired("token")
	BucketCmd.MarkFlagRequired("project-id")
}

func ListBuckets(token string, projectID string) {
	url := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b?project=%s", projectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			fmt.Println("Error formatting JSON:", err)
			return
		}

		fmt.Println(string(formattedJSON))
	} else {
		fmt.Printf("Failed to retrieve storage buckets. Status code: %d, token is invalid or expired\n", resp.StatusCode)
	}
}
