package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gunnyworks/gunny/pkg/gunny"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type rootCLIArgs struct {
	showVersion            bool
	verbose                bool
	configFilePath         string
	configFileFormatString string
	templateContent        string
	templateFilePath       string
	namedDataValues        []string
	stdinDataFormatString  string
	outputBasePath         string
}

func newRootCLIArgs() *rootCLIArgs {
	return &rootCLIArgs{
		configFileFormatString: string(gunny.DataFormatJSON),
		stdinDataFormatString:  string(gunny.DataFormatJSON),
	}
}

func newRootCmd(version string, commit string) *cobra.Command {
	rootArgs := newRootCLIArgs()
	cmd := &cobra.Command{
		Use:   "gunny",
		Short: "Gunny helps you weave your data through templates to generate static text",
		Run: func(cmd *cobra.Command, args []string) {
			logLevel := slog.LevelInfo
			if rootArgs.verbose {
				logLevel = slog.LevelDebug
			}
			slog.SetDefault(slog.New(
				tint.NewHandler(os.Stderr, &tint.Options{
					Level:      logLevel,
					TimeFormat: time.RFC3339,
				}),
			))

			// Handle showing of version and exiting
			if rootArgs.showVersion {
				if len(commit) > 0 {
					fmt.Printf("%s-%s\n", version, commit)
				} else {
					fmt.Println(version)
				}
				os.Exit(0)
			}

			var config *gunny.PipelineConfig
			if len(rootArgs.configFilePath) > 0 {
				var err error
				var configFileFormat *gunny.DataFormat
				if len(rootArgs.configFileFormatString) > 0 {
					format, err := gunny.DataFormatFromString(rootArgs.configFileFormatString)
					if err != nil {
						slog.Error("Invalid argument(s)", "error", err)
					}
					configFileFormat = &format
				}
				config, err = gunny.ReadPipelineConfigFromFile(rootArgs.configFilePath, configFileFormat)
				if err != nil {
					slog.Error("Failed to load configuration file", "filename", rootArgs.configFilePath, "error", err)
					os.Exit(1)
				}
			} else {
				config = gunny.NewPipelineConfigWithDefaults()
			}

			// Validate template configuration
			if len(rootArgs.templateContent) > 0 {
				config.SetTemplateContent(rootArgs.templateContent)
			} else if len(rootArgs.templateFilePath) > 0 {
				config.SetTemplateFile(rootArgs.templateFilePath)
			}

			// Validate stdin data format
			if len(rootArgs.stdinDataFormatString) > 0 {
				stdinDataFormat, err := gunny.DataFormatFromString(rootArgs.stdinDataFormatString)
				if err != nil {
					slog.Error("Invalid argument(s)", "error", err)
					os.Exit(1)
				}
				config.SetStdinDataFormat(stdinDataFormat)
			}

			config.OutputBasePath = rootArgs.outputBasePath

			if err := config.Validate(); err != nil {
				slog.Error("Invalid Gunny configuration", "error", err)
				os.Exit(1)
			}

			if err := gunny.RenderUsingConfig(context.Background(), config, rootArgs.namedDataValues); err != nil {
				slog.Error("Failed to render", "error", err)
				os.Exit(1)
			}
		},
		Example: `
    # Render the value "Michael" (named "name") through the given template.
    gunny -t 'Hello {{name}}!' -d name=Michael
    
    # Supply data via stdin (in JSON format, by default)
    echo '{"name": "Michael"}' | gunny -t 'Hello {{name}}!'

    # Drive rendering pipeline from configuration file
    gunny -c gunny.yaml`,
	}
	cmd.PersistentFlags().BoolVar(&rootArgs.showVersion, "version", false, "show the program version and exit immediately")
	cmd.PersistentFlags().StringVarP(&rootArgs.configFilePath, "config", "c", "", "configuration file defining the Gunny rendering pipeline")
	cmd.PersistentFlags().StringVar(&rootArgs.configFileFormatString, "config-file-format", "", "explicitly set configuration file format; if not set, auto-detects based on config file extension")
	cmd.PersistentFlags().BoolVarP(&rootArgs.verbose, "verbose", "v", false, "increase output logging verbosity")
	cmd.PersistentFlags().StringVarP(&rootArgs.templateContent, "template", "t", "", "template content")
	cmd.PersistentFlags().StringVar(&rootArgs.templateFilePath, "template-file", "", "file from which to read template")
	cmd.PersistentFlags().StringArrayVarP(&rootArgs.namedDataValues, "data", "d", []string{}, "specify named data values (name=value) to inject into the template")
	cmd.PersistentFlags().StringVar(&rootArgs.stdinDataFormatString, "stdin-format", rootArgs.stdinDataFormatString, "when supplying data via stdin, the format in which to interpret it")
	cmd.PersistentFlags().StringVar(&rootArgs.outputBasePath, "output-base-path", ".", "base path into which to write output files whose file name is specified as relative")
	return cmd
}
