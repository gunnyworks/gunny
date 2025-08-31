package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gunnyworks/gunny/pkg/gunny"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

func newRootCmd(version string, commit string) *cobra.Command {
	showVersion := false
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
			if showVersion {
				if len(commit) > 0 {
					fmt.Printf("%s-%s\n", version, commit)
				} else {
					fmt.Println(version)
				}
				os.Exit(0)
			}
			stdinDataFormat, err := gunny.DataFormatFromString(stdinDataFormatString)
			if err != nil {
				slog.Error("Invalid argument(s)", "error", err)
				os.Exit(1)
			}
			pipeline, err := gunny.NewPipeline(
				gunny.WithMustacheTemplateFromReader(strings.NewReader(templateContent)),
				gunny.WithDataFromReader(os.Stdin, stdinDataFormat),
				gunny.WithDataFromNameValuePairs(namedDataValues),
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
    gunny -t 'Hello {{name}}!' -d name=Michael
	
	# Supply data via stdin (in JSON format, by default)
	echo '{"name": "Michael"}' | gunny -t 'Hello {{name}}!'`,
	}
	cmd.PersistentFlags().BoolVar(&showVersion, "version", false, "show the program version and exit immediately")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase output logging verbosity")
	cmd.PersistentFlags().StringVarP(&templateContent, "template", "t", "", "template content")
	cmd.PersistentFlags().StringArrayVarP(&namedDataValues, "data", "d", []string{}, "specify named data values (name=value) to inject into the template")
	cmd.PersistentFlags().StringVar(&stdinDataFormatString, "stdin-format", stdinDataFormatString, "when supplying data via stdin, the format in which to interpret it")
	return cmd
}
