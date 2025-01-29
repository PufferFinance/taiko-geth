package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"os"
	"perf-test/spammer"
)

var startNumber int

// helloCmd represents the 'hello' command
var startTest = &cobra.Command{
	Use:   "start",
	Short: "Start the performance test with pre-set accounts",
	Long: "Start the performance test with pre-set accounts." +
		"Use the accounts argument to specify how many accounts you want test with",,
	"Example: ./perf-test start -a 10",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(false)
		spammer := spammer.New(url, chainID, logger, accounts, maxTxsPerAccount, prefundedAccounts[:startNumber])

		spammer.Start()
	},
}

func init() {
	// Add helloCmd as a subcommand of rootCmd
	rootCmd.AddCommand(startTest)

	// Register flags for the hello command
	//helloCmd.Flags().StringVarP(&name, "num", "n", 1, "Name to greet")
	startTest.Flags().IntVarP(&startNumber, "accounts", "a", 1, "Number of accounts to prefund")
}
