package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gunnyworks/gunny/pkg/gunny"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	verbose := false
	templateContent := ""
	namedDataValues := []string{}
	stdinDataFormatString := string(gunny.DataFormatJSON)

	cmd := &cobra.Command{
		Use:   "gunny",
		Short: "Gunny helps you weave your data through templates to generate static text",
		Run: func(cmd *cobra.Command, args []string) {
			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}
			slog.SetDefault(slog.New(
				tint.NewHandler(os.Stderr, &tint.Options{
					Level:      logLevel,
					TimeFormat: time.RFC3339,
				}),
			))
			stdinDataFormat, err := gunny.DataFormatFromString(stdinDataFormatString)
			if err != nil {
				slog.Error("Invalid argument(s)", "error", err)
				os.Exit(1)
			}
			pipeline, err := gunny.NewPipeline(
				gunny.WithMustacheTemplateFromReader(strings.NewReader(templateContent)),
				gunny.WithDataFromNameValuePairs(namedDataValues),
				gunny.WithDataFromReader(os.Stdin, stdinDataFormat),
			)
			if err != nil {
				slog.Error("Failed to initialize Gunny", "error", err)
				os.Exit(1)
			}
			if err := pipeline.Render(context.Background()); err != nil {
				slog.Error("Failed to render", "error", err)
				os.Exit(1)
			}
		},
		Example: `
    # Render the value "Michael" (named "name") through the given template.
    gunny -t 'Hello {{name}}!' -d name=Michael`,
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase output logging verbosity")
	cmd.PersistentFlags().StringVarP(&templateContent, "template", "t", "", "template content")
	cmd.PersistentFlags().StringArrayVarP(&namedDataValues, "data", "d", []string{}, "specify named data values (name=value) to inject into the template")
	cmd.PersistentFlags().StringVar(&stdinDataFormatString, "stdin-format", stdinDataFormatString, "when supplying data via stdin, the format in which to interpret it")
	return cmd
}
