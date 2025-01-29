package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"os"
	"perf-test/spammer"
)

var balancesNumber int

// helloCmd represents the 'hello' command
var balancesCommand = &cobra.Command{
	Use:   "balances",
	Short: "Get the balances of the pre-set accounts",
	Long: "Get the balances of the selected pre-set accounts." +
		"Use the accounts argument to specify how many accounts you want prefunded" +
		"Example: ./perf-test balances -a 10",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(false)
		spammer := spammer.New(url, chainID, logger, accounts, maxTxsPerAccount, prefundedAccounts[:balancesNumber])

		spammer.GetBalances()
	},
}

func init() {
	// Add helloCmd as a subcommand of rootCmd
	rootCmd.AddCommand(balancesCommand)

	// Register flags for the hello command
	//helloCmd.Flags().StringVarP(&name, "num", "n", 1, "Name to greet")
	balancesCommand.Flags().IntVarP(&balancesNumber, "accounts", "a", 1, "Number of accounts to prefund")
}
