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

func newRootCmd() *cobra.Command {
	verbose := false
	templateContent := ""
	namedDataValues := []string{}

	cmd := &cobra.Command{
		Use:   "gunny",
		Short: "Gunny helps you weave your data through templates to generate static text",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			pipeline, err := gunny.NewPipeline(
				gunny.WithMustacheTemplateFromReader(strings.NewReader(templateContent)),
				gunny.WithDataFromNameValuePairs(namedDataValues),
			)
			if err != nil {
				return fmt.Errorf("failed to initialize Gunny: %w", err)
			}
			if err := pipeline.Render(context.Background()); err != nil {
				return fmt.Errorf("failed to render: %w", err)
			}
			return nil
		},
		Example: `
    # Render the value "Michael" (named "name") through the given template.
    gunny -t 'Hello {{name}}!' -d name=Michael`,
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase output logging verbosity")
	cmd.PersistentFlags().StringVarP(&templateContent, "template", "t", "", "template content")
	cmd.PersistentFlags().StringArrayVarP(&namedDataValues, "data", "d", []string{}, "specify named data values (name=value) to inject into the template")
	return cmd
}
