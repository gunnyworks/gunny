package gunny_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gunnyworks/gunny/pkg/gunny"
	"github.com/stretchr/testify/require"
)

type pipelineConfigTestCase struct {
	name           string
	path           string
	format         gunny.DataFormat
	expectErr      bool
	expectedConfig *gunny.PipelineConfig
}

var (
	dataFormatYAML          = gunny.DataFormatYAML
	templateFilePath        = "/path/to/template.mustache"
	pipelineConfigTestCases = []pipelineConfigTestCase{
		{
			name:      "trivial JSON configuration",
			path:      "trivial.json",
			format:    gunny.DataFormatJSON,
			expectErr: false,
			expectedConfig: &gunny.PipelineConfig{
				Version: gunny.ConfigVersion1,
				DataSources: []*gunny.DataSourceConfig{
					{
						Type: gunny.DataSourceStdin,
						Config: &gunny.StdinDataSourceConfig{
							Format: gunny.DataFormatJSON,
						},
					},
				},
				Renderer:       gunny.NewRendererConfigWithDefaults(),
				OutputBasePath: ".",
				Outputs:        []*gunny.OutputConfig{gunny.NewOutputConfigWithDefaults()},
			},
		},
		{
			name:      "trivial YAML configuration",
			path:      "trivial.yaml",
			format:    gunny.DataFormatYAML,
			expectErr: false,
			expectedConfig: &gunny.PipelineConfig{
				Version: gunny.ConfigVersion1,
				DataSources: []*gunny.DataSourceConfig{
					{
						Type: gunny.DataSourceStdin,
						Config: &gunny.StdinDataSourceConfig{
							Format: gunny.DataFormatJSON,
						},
					},
				},
				Renderer:       gunny.NewRendererConfigWithDefaults(),
				OutputBasePath: ".",
				Outputs:        []*gunny.OutputConfig{gunny.NewOutputConfigWithDefaults()},
			},
		},
		{
			name:      "JSON config with multiple data sources and single output",
			path:      "multi-data-source-single-output.json",
			format:    gunny.DataFormatJSON,
			expectErr: false,
			expectedConfig: &gunny.PipelineConfig{
				Version: gunny.ConfigVersion1,
				DataSources: []*gunny.DataSourceConfig{
					{
						Type: gunny.DataSourceEnvVars,
						Config: &gunny.EnvVarsDataSourceConfig{
							Expected: []string{"ENV_VAR1", "ENV_VAR2"},
							Optional: []string{"ENV_VAR3"},
						},
					},
					{
						Type: gunny.DataSourceFile,
						Config: &gunny.FileDataSourceConfig{
							Path:   "/path/to/file.yaml",
							Format: &dataFormatYAML,
						},
					},
					{
						Type: gunny.DataSourceCLIArgs,
						Config: &gunny.CLIArgsDataSourceConfig{
							Expected: []string{"arg1", "arg2"},
							Optional: []string{"arg3", "arg4"},
						},
					},
				},
				Renderer: &gunny.RendererConfig{
					Type: gunny.RendererMustache,
					Config: &gunny.MustacheRendererConfig{
						TemplateFile: &templateFilePath,
					},
				},
				OutputBasePath: ".",
				Outputs: []*gunny.OutputConfig{
					{
						Type: gunny.OutputFile,
						Config: &gunny.FileOutputConfig{
							Path: "/path/to/output.html",
						},
					},
				},
			},
		},
		{
			name:      "YAML config with multiple data sources and single output",
			path:      "multi-data-source-single-output.yaml",
			format:    gunny.DataFormatYAML,
			expectErr: false,
			expectedConfig: &gunny.PipelineConfig{
				Version: gunny.ConfigVersion1,
				DataSources: []*gunny.DataSourceConfig{
					{
						Type: gunny.DataSourceEnvVars,
						Config: &gunny.EnvVarsDataSourceConfig{
							Expected: []string{"ENV_VAR1", "ENV_VAR2"},
							Optional: []string{"ENV_VAR3"},
						},
					},
					{
						Type: gunny.DataSourceFile,
						Config: &gunny.FileDataSourceConfig{
							Path:   "/path/to/file.yaml",
							Format: &dataFormatYAML,
						},
					},
					{
						Type: gunny.DataSourceCLIArgs,
						Config: &gunny.CLIArgsDataSourceConfig{
							Expected: []string{"arg1", "arg2"},
							Optional: []string{"arg3", "arg4"},
						},
					},
				},
				Renderer: &gunny.RendererConfig{
					Type: gunny.RendererMustache,
					Config: &gunny.MustacheRendererConfig{
						TemplateFile: &templateFilePath,
					},
				},
				OutputBasePath: ".",
				Outputs: []*gunny.OutputConfig{
					{
						Type: gunny.OutputFile,
						Config: &gunny.FileOutputConfig{
							Path: "/path/to/output.html",
						},
					},
				},
			},
		},
	}
)

func TestPipelineConfigParsing(t *testing.T) {
	for _, testCase := range pipelineConfigTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", testCase.path))
			require.NoError(t, err)
			config, err := gunny.ReadPipelineConfig(f, testCase.format)
			if testCase.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.expectedConfig, config)
		})
	}
}
