package azure

import (
	"github.com/spf13/cobra"
)

var BlobCmd = &cobra.Command{
	Use:   "blob",
	Short: "Interact with Azure Blob Storage",
}

func init() {
	BlobCmd.AddCommand(downloadCmd)
	BlobCmd.AddCommand(listCmd)
}
