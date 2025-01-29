package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"os"
	"perf-test/spammer"
)

var prefundNumber int

// helloCmd represents the 'hello' command
var prefundCommand = &cobra.Command{
	Use:   "prefund",
	Short: "Prefund the selected pre-set accounts",
	Long: "Prefund the selected pre-set accounts." +
		"Use the accounts argument to specify how many accounts you want prefunded",
	"Example: ./perf-test prefund -a 10",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(false)
		spammer := spammer.New(url, chainID, logger, accounts, maxTxsPerAccount, prefundedAccounts[:prefundNumber])

		spammer.PrefundAccounts()
	},
}

func init() {
	// Add helloCmd as a subcommand of rootCmd
	rootCmd.AddCommand(prefundCommand)

	// Register flags for the hello command
	//helloCmd.Flags().StringVarP(&name, "num", "n", 1, "Name to greet")
	prefundCommand.Flags().IntVarP(&prefundNumber, "accounts", "a", 1, "Number of accounts to prefund")
}
