package main

import (
	"log/slog"
	"os"
)

var (
	version = "dev"
	commit  = ""
)

func main() {
	if err := newRootCmd(version, commit).Execute(); err != nil {
		slog.Error("Failed to run Gunny", "error", err)
		os.Exit(1)
	}
}
