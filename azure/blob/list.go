package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/cobra"
)

// listCmd is the blob enumeration command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Enumerate blobs in an Azure Storage container",
	Long:  "The blob command lists blobs in a given Azure Storage container using access key authentication.",
	Run: func(cmd *cobra.Command, args []string) {
		account, _ := cmd.Flags().GetString("account")
		key, _ := cmd.Flags().GetString("key")
		container, _ := cmd.Flags().GetString("container")

		if account == "" || key == "" || container == "" {
			log.Fatalf("Error: --account, --key, and --container are required")
		}

		err := listBlobs(account, key, container)
		if err != nil {
			log.Fatalf("Error listing blobs: %v", err)
		}
	},
}

func init() {
	listCmd.Flags().StringP("account", "a", "", "Azure Storage account name (required)")
	listCmd.Flags().StringP("key", "k", "", "Azure Storage account key (required)")
	listCmd.Flags().StringP("container", "c", "", "Azure container name (required)")

	listCmd.MarkFlagRequired("account")
	listCmd.MarkFlagRequired("key")
	listCmd.MarkFlagRequired("container")
}

func listBlobs(account, key, container string) error {
	ctx := context.Background()
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", account)

	cred, err := azblob.NewSharedKeyCredential(account, key)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	serviceClient, err := azblob.NewClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create service client: %w", err)
	}

	containerClient := serviceClient.ServiceClient().NewContainerClient(container)

	pager := containerClient.NewListBlobsFlatPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list blobs: %w", err)
		}
		for _, blob := range resp.Segment.BlobItems {
			fmt.Println("Blob Name:", *blob.Name)
		}
	}

	return nil
}
