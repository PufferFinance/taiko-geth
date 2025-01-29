package main

import (
	"github.com/charmbracelet/log"
	"os"
	"perf-test/cmd"
)

func main() {
	cmd.Execute()

	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)
}
