package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		slog.Error("Failed to run Gunny", "error", err)
		os.Exit(1)
	}
}
