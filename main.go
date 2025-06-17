package main

import (
	"log"
	"os"

	azure "github.com/f0rk3b0mb/GoCloudGhost/azure"
	gcp "github.com/f0rk3b0mb/GoCloudGhost/gcp"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GoCloudGhost",
	Short: "GoCloudGhost - Cloud Enumerator",
	Long:  `GoCloudGhost allows you to authenticate with cloud and enumerate when testing cloud security.`,
}

func init() {
	// Register blob command and its subcommands
	rootCmd.AddCommand(azure.AzureCmd)
	rootCmd.AddCommand(gcp.GcpCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

}
