package gcplist

import (
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List GCP Resources",
}

func init() {
	ListCmd.AddCommand(ComputeCmd)
	ListCmd.AddCommand(BucketCmd)
}
