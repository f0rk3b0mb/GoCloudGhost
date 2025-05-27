package azure

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a blob from an Azure Storage container",
	Run: func(cmd *cobra.Command, args []string) {
		account, _ := cmd.Flags().GetString("account")
		key, _ := cmd.Flags().GetString("key")
		container, _ := cmd.Flags().GetString("container")
		blobName, _ := cmd.Flags().GetString("blob")
		output, _ := cmd.Flags().GetString("output")

		if account == "" || key == "" || container == "" || blobName == "" || output == "" {
			log.Fatalf("--account, --key, --container, --blob, and --output are required")
		}

		err := downloadBlob(account, key, container, blobName, output)
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
	},
}

func init() {
	downloadCmd.Flags().String("account", "", "Azure Storage account name (required)")
	downloadCmd.Flags().String("key", "", "Azure Storage account key (required)")
	downloadCmd.Flags().String("container", "", "Azure container name (required)")
	downloadCmd.Flags().String("blob", "", "Name of the blob to download (required)")
	downloadCmd.Flags().String("output", "", "Path to save the downloaded file (required)")

	downloadCmd.MarkFlagRequired("account")
	downloadCmd.MarkFlagRequired("key")
	downloadCmd.MarkFlagRequired("container")
	downloadCmd.MarkFlagRequired("blob")
	downloadCmd.MarkFlagRequired("output")

	BlobCmd.AddCommand(downloadCmd)
}

func downloadBlob(account, key, container, blobName, output string) error {
	ctx := context.Background()
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", account)

	cred, err := azblob.NewSharedKeyCredential(account, key)
	if err != nil {
		return err
	}

	serviceClient, err := azblob.NewServiceClientWithSharedKey(url, cred, nil)
	if err != nil {
		return err
	}

	blobClient := serviceClient.NewContainerClient(container).NewBlobClient(blobName)
	resp, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		return err
	}

	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
