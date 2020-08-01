package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of Katana",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func initVersionCmd() {
	rootCmd.AddCommand(versionCmd)
}
