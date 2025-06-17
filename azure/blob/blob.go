package blob

import (
	"github.com/spf13/cobra"
)

var BlobCmd = &cobra.Command{
	Use:   "blob",
	Short: "Interact with Azure Blob Storage",
}

func init() {
	BlobCmd.AddCommand(listCmd)
	BlobCmd.AddCommand(downloadCmd)
}
