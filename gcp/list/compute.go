package gcplist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var ComputeCmd = &cobra.Command{
	Use:   "compute",
	Short: "List GCP Compute Resources",
	Long:  `This command allows you to list GCP Compute resources.`,
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
		ListComputeResources(token, project_id)
	},
}

func init() {
	// Add the Compute command to the List command
	//ListCmd.AddCommand(ComputeCmd)

	// Add flags for Compute command if needed
	ComputeCmd.Flags().String("token", "", "GCP OAuth2 token")
	ComputeCmd.Flags().String("project-id", "", "GCP Project ID")
	ComputeCmd.MarkFlagRequired("token")
	ComputeCmd.MarkFlagRequired("project-id")
}

func ListComputeResources(token string, projectID string) {
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/aggregated/instances", projectID)

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

		// Parse and format the JSON
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
		fmt.Printf("Failed to retrieve compute resources. Status code: %d ,  token is invalid or expired\n", resp.StatusCode)
	}
}
