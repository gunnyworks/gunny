package main

import (
	"log/slog"
	"os"
)

var (
	gunnyVersion = "dev"
	gunnyCommit  = ""
)

func main() {
	if err := newRootCmd(gunnyVersion, gunnyCommit).Execute(); err != nil {
		slog.Error("Failed to run Gunny", "error", err)
		os.Exit(1)
	}
}
