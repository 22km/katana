package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ktn",
	Short: "Generate golang http server in one step.",
	Long: `**KATANA**
  Sharp, pure, eazy golang http server. 
  Generate golang http server in one step.
  https://github.com/22km/katana`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	initVersionCmd()
	initNewCmd()
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
