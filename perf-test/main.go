package main

import (
	"os"

	"perf-test/spammer"

	"github.com/charmbracelet/log"
)


func main() {
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)

	spammer := spammer.New(url, chainID, logger, accounts, maxTxsPerAccount, prefundedAccounts)

	// Prefund accounts with ETH
	// spammer.PrefundAccounts()

	// Send n(maxTxsPerAccount) transactions per prefunded account
	spammer.Start()
	// spammer.GetBalances()
}

