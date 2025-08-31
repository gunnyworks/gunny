package main

import (
	"log/slog"
	"os"
)

var (
	gunnyVersion = "dev"
	gunnyBuild   = ""
)

func main() {
	if err := newRootCmd(gunnyVersion, gunnyBuild).Execute(); err != nil {
		slog.Error("Failed to run Gunny", "error", err)
		os.Exit(1)
	}
}
