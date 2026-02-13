package azure

import (
	"github.com/f0rk3b0mb/GoCloudGhost/auth"
	blob "github.com/f0rk3b0mb/GoCloudGhost/azure/blob"
	management "github.com/f0rk3b0mb/GoCloudGhost/azure/enum"
	"github.com/spf13/cobra"
)

var AzureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Interact with Azure",
}

func init() {
	AzureCmd.AddCommand(blob.BlobCmd)
	AzureCmd.AddCommand(management.MgmtCmd)
	AzureCmd.AddCommand(auth.AuthCmd)
}
