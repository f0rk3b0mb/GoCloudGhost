package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GoCloudGhost",
	Short: "GoCloudGhost - Cloud Enumerator",
	Long:  `GoCloudGhost allows you to authenticate with cloud and enumerate when testing cloud security.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}
