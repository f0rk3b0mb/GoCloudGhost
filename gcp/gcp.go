package gcp

import (
	gcpenum "github.com/f0rk3b0mb/GoCloudGhost/gcp/enum"
	gcplist "github.com/f0rk3b0mb/GoCloudGhost/gcp/list"
	"github.com/spf13/cobra"
)

var GcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Interact with GCP",
}

func init() {
	GcpCmd.AddCommand(gcpenum.EnumCmd)
	GcpCmd.AddCommand(gcplist.ListCmd)
}
