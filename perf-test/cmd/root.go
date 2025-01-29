package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "perf-test",
	Short: "perf-test is a tool to run tests for sending transactions",
	Long:  `Send a lot of transactions at once to test the TPS.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Invalid command, please run ./perf-test --help")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
