package gunny

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// Pipeline encapsulates a Gunny rendering pipeline.
type Pipeline struct {
	resolvers DataResolverMap
	renderer  Renderer
	writers   []io.WriteCloser
	logger    Logger
}

type newPipelineOption func(*Pipeline) error

// WithMustacheTemplateFromReader reads all it can from the given reader,
// attempting to interpret it as a Mustache template.
func WithMustacheTemplateFromReader(r io.Reader) newPipelineOption {
	return func(pipeline *Pipeline) error {
		renderer, err := NewMustacheTemplateRenderer(r)
		if err != nil {
			return err
		}
		pipeline.renderer = renderer
		return nil
	}
}

// WithMustacheTemplateFromFile is similar to [WithMustacheTemplateFromReader],
// except that it handles opening, reading from and closing the supplied file.
func WithMustacheTemplateFromFile(filename string) newPipelineOption {
	return func(pipeline *Pipeline) error {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				pipeline.logger.Error("Failed to close template file", "filename", filename, "error", err)
			}
		}()
		renderer, err := NewMustacheTemplateRenderer(f)
		if err != nil {
			return err
		}
		pipeline.renderer = renderer
		return nil
	}
}

// WithNamedDataResolver allows one to supply a custom named data resolver to a
// Gunny rendering pipeline.
func WithNamedDataResolver(resolver DataResolverMap) newPipelineOption {
	return func(pipeline *Pipeline) error {
		pipeline.resolvers.Merge(resolver)
		return nil
	}
}

// WithDataFromNameValuePairs allows one to supply data to a Gunny pipeline
// from "name=value" strings (usually from the command line).
func WithDataFromNameValuePairs(nameValuePairStrings []string) newPipelineOption {
	return func(pipeline *Pipeline) error {
		resolver, err := NewNameValuePairsDataResolver(nameValuePairStrings)
		if err != nil {
			return err
		}
		pipeline.resolvers.Merge(resolver)
		return nil
	}
}

// WithDataFromMap trivially merges the given map into the pipeline's data
// resolvers.
func WithDataFromMap(m map[string]any) newPipelineOption {
	return func(pipeline *Pipeline) error {
		resolver, err := NewInMemoryDataResolverMap(m)
		if err != nil {
			return err
		}
		pipeline.resolvers.Merge(resolver)
		return nil
	}
}

// WithDataFromCLIArgs uses the same strategy to extract data values as
// [WithDataFromNameValuePairs], but allows one to specify values that are
// expected or optional. If an expected CLI argument is not present, an error
// will be produced.
//
// Only values whose names match the expected and optional arguments lists will
// be merged into the pipeline's input data set. All other supplied arguments
// will be filtered out.
func WithDataFromCLIArgs(cliArgsNamedDataValues []string, expectedArgs []string, optionalArgs []string) newPipelineOption {
	return func(pipeline *Pipeline) error {
		resolver, err := NewNameValuePairsDataResolver(cliArgsNamedDataValues)
		if err != nil {
			return err
		}
		filteredResolver := make(DataResolverMap)
		for _, expectedArg := range expectedArgs {
			value, ok := resolver[expectedArg]
			if !ok {
				return MissingCLIArgError{Name: expectedArg}
			}
			filteredResolver[expectedArg] = value
		}
		for _, optionalArg := range optionalArgs {
			value, ok := resolver[optionalArg]
			if ok {
				filteredResolver[optionalArg] = value
			}
		}
		pipeline.resolvers.Merge(filteredResolver)
		return nil
	}
}

// WithDataFromEnvVars injects environment variables as data into the rendering
// pipeline. If an expected environment variable is not present, an error will
// be produced. If an optional environment variable is not present, it will
// simply not be set.
func WithDataFromEnvVars(expectedEnvVars []string, optionalEnvVars []string) newPipelineOption {
	return func(pipeline *Pipeline) error {
		resolver := make(DataResolverMap)
		for _, envVar := range expectedEnvVars {
			value, exists := os.LookupEnv(envVar)
			if !exists {
				return MissingEnvVarError{Name: envVar}
			}
			resolver[envVar] = NewInMemoryDataValue(value)
		}
		for _, envVar := range optionalEnvVars {
			value, exists := os.LookupEnv(envVar)
			if exists {
				resolver[envVar] = NewInMemoryDataValue(value)
			}
		}
		pipeline.resolvers.Merge(resolver)
		return nil
	}
}

// WithDataFromStdin is similar to [WithDataFromReader], but first checks
// whether there is any data available to read from stdin. Otherwise results in
// an empty data set.
func WithDataFromStdin(format DataFormat) newPipelineOption {
	return func(pipeline *Pipeline) error {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			// Do nothing
			return nil
		}
		return WithDataFromReader(os.Stdin, format)(pipeline)
	}
}

// WithDataFromReader allows one to supply data to a Gunny pipeline from a
// reader. The format must be specified.
func WithDataFromReader(reader io.Reader, format DataFormat) newPipelineOption {
	return func(pipeline *Pipeline) error {
		resolver, err := NewDataResolverFromReader(reader, format)
		if err != nil {
			return err
		}
		pipeline.resolvers.Merge(resolver)
		return nil
	}
}

// WithDataFromFile is similar to [WithDataFromReader], except it handles
// opening the supplied file and attempting to auto-detect its file type.
//
// If overrideFormat is supplied (non-null), no file type autodetection will be
// performed.
func WithDataFromFile(filename string, overrideFormat *DataFormat) newPipelineOption {
	return func(pipeline *Pipeline) error {
		var format DataFormat
		if overrideFormat != nil {
			format = *overrideFormat
		} else {
			var err error
			format, err = GetFileDataFormatFromFilename(filename)
			if err != nil {
				return err
			}
		}
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				pipeline.logger.Error("Failed to close file", "filename", filename, "error", err)
			}
		}()
		resolver, err := NewDataResolverFromReader(f, format)
		if err != nil {
			return err
		}
		pipeline.resolvers.Merge(resolver)
		return nil
	}
}

// WithOutputWriter allows one to specify where to write rendered output.
func WithOutputWriter(w io.Writer) newPipelineOption {
	return func(pipeline *Pipeline) error {
		pipeline.writers = append(pipeline.writers, &writeCloserNoopClose{w: w})
		return nil
	}
}

// WithOutputWriteCloser allows one to specify where to write rendered output,
// where the output needs to eventually be closed after rendering.
func WithOutputWriteCloser(w io.WriteCloser) newPipelineOption {
	return func(pipeline *Pipeline) error {
		pipeline.writers = append(pipeline.writers, w)
		return nil
	}
}

// WithOutputFile allows one to render output to a file.
//
// [Pipeline.Close] must always be called after [Pipeline] creation if this
// option is used, regardless of whether or not rendering was executed and
// regardless of its success.
func WithOutputFile(filename string) newPipelineOption {
	return func(pipeline *Pipeline) error {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		pipeline.writers = append(pipeline.writers, f)
		return nil
	}
}

// WithLogger allows customization of the logging mechanism used. By default,
// Go's slog is used.
func WithLogger(logger Logger) newPipelineOption {
	return func(pipeline *Pipeline) error {
		pipeline.logger = logger
		return nil
	}
}

// NewPipeline constructs an empty Gunny rendering pipeline that, by default,
// renders to stdout.
func NewPipeline(opts ...newPipelineOption) (*Pipeline, error) {
	pipeline := &Pipeline{
		renderer:  &NullRenderer{},
		resolvers: make(DataResolverMap),
		writers:   make([]io.WriteCloser, 0),
		logger:    &SlogLogger{},
	}
	for _, opt := range opts {
		if err := opt(pipeline); err != nil {
			return nil, err
		}
	}
	// Default to stdout if no output writer was specified.
	if len(pipeline.writers) == 0 {
		pipeline.writers = []io.WriteCloser{os.Stdout}
	}
	return pipeline, nil
}

// Render executes the entire pipeline rendering operation fully.
func (p *Pipeline) Render(ctx context.Context) error {
	p.logger.Debug("Rendering")
	ctxWithLogger := NewContextWithLogger(ctx, p.logger)
	return p.renderer.Render(ctxWithLogger, p.resolvers, io.MultiWriter(lo.Map(p.writers, func(w io.WriteCloser, _ int) io.Writer {
		return w
	})...))
}

// Close calls the underlying writer's Close method. This should always be
// called after rendering, whether successful or not.
func (p *Pipeline) Close() error {
	errors := make([]error, 0)
	for _, writer := range p.writers {
		if err := writer.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return MultiWriterCloseFailedError{Causes: errors}
	}
	return nil
}

// RenderUsingConfig executes an entire rendering pipeline end-to-end based on
// the given configuration.
func RenderUsingConfig(
	ctx context.Context,
	config *PipelineConfig,
	cliArgsNamedDataValues []string,
) error {
	opts, err := buildOptionsFromDataSourcesConfig(config.DataSources, cliArgsNamedDataValues)
	if err != nil {
		return err
	}
	rendererOpts, err := buildOptionsFromRendererConfig(config.Renderer)
	if err != nil {
		return err
	}
	opts = append(opts, rendererOpts...)
	outputOpts, err := buildOptionsFromOutputsConfig(config.OutputBasePath, config.Outputs)
	if err != nil {
		return err
	}
	opts = append(opts, outputOpts...)
	pipeline, err := NewPipeline(opts...)
	if err != nil {
		return err
	}
	return pipeline.Render(ctx)
}

func buildOptionsFromDataSourcesConfig(config []*DataSourceConfig, cliArgsNamedDataValues []string) ([]newPipelineOption, error) {
	opts := []newPipelineOption{}
	filteredCLIArgs := false
	for _, dataSource := range config {
		switch dataSource.Type {
		case DataSourceEnvVars:
			dataSourceConfig, ok := dataSource.Config.(*EnvVarsDataSourceConfig)
			if !ok {
				return nil, InvalidConfigTypeError{
					Expected: EnvVarsDataSourceConfig{},
					Actual:   dataSource.Config,
				}
			}
			opts = append(opts, WithDataFromEnvVars(dataSourceConfig.Expected, dataSourceConfig.Optional))
		case DataSourceCLIArgs:
			dataSourceConfig, ok := dataSource.Config.(*CLIArgsDataSourceConfig)
			if !ok {
				return nil, InvalidConfigTypeError{
					Expected: CLIArgsDataSourceConfig{},
					Actual:   dataSource.Config,
				}
			}
			opts = append(opts, WithDataFromCLIArgs(cliArgsNamedDataValues, dataSourceConfig.Expected, dataSourceConfig.Optional))
			filteredCLIArgs = true
		case DataSourceStdin:
			dataSourceConfig, ok := dataSource.Config.(*StdinDataSourceConfig)
			if !ok {
				return nil, InvalidConfigTypeError{
					Expected: StdinDataSourceConfig{},
					Actual:   dataSource.Config,
				}
			}
			opts = append(opts, WithDataFromStdin(dataSourceConfig.Format))
		case DataSourceFile:
			dataSourceConfig, ok := dataSource.Config.(*FileDataSourceConfig)
			if !ok {
				return nil, InvalidConfigTypeError{
					Expected: FileDataSourceConfig{},
					Actual:   dataSource.Config,
				}
			}
			opts = append(opts, WithDataFromFile(dataSourceConfig.Path, dataSourceConfig.Format))
		case DataSourceMap:
			dataSourceConfig, ok := dataSource.Config.(map[string]any)
			if !ok {
				return nil, InvalidConfigTypeError{
					Expected: map[string]any{},
					Actual:   dataSource.Config,
				}
			}
			opts = append(opts, WithDataFromMap(dataSourceConfig))
		default:
			return nil, InvalidDataSourceTypeError{
				Supplied: string(dataSource.Type),
				ValidValues: lo.Map(validDataSourceTypeValues, func(t DataSourceType, _ int) string {
					return string(t)
				}),
			}
		}
	}
	// If we don't have any explicit filtering of data values from CLI
	// arguments but CLI argument-based data values are supplied, add them as
	// our final data source to override whatever other named data values we've
	// built up.
	if !filteredCLIArgs && len(cliArgsNamedDataValues) > 0 {
		opts = append(opts, WithDataFromNameValuePairs(cliArgsNamedDataValues))
	}
	return opts, nil
}

func buildOptionsFromRendererConfig(config *RendererConfig) ([]newPipelineOption, error) {
	opts := []newPipelineOption{}
	switch config.Type {
	case RendererMustache:
		rendererConfig, ok := config.Config.(*MustacheRendererConfig)
		if !ok {
			return nil, InvalidConfigTypeError{
				Expected: MustacheRendererConfig{},
				Actual:   config.Config,
			}
		}
		switch {
		case rendererConfig.Template != nil:
			opts = append(opts, WithMustacheTemplateFromReader(strings.NewReader(*rendererConfig.Template)))
		case rendererConfig.TemplateFile != nil:
			opts = append(opts, WithMustacheTemplateFromFile(*rendererConfig.TemplateFile))
		default:
			return nil, ErrMissingTemplate
		}
	default:
		return nil, InvalidRendererTypeError{
			Supplied: string(config.Type),
			ValidValues: lo.Map(validRendererTypes, func(t RendererType, _ int) string {
				return string(t)
			}),
		}
	}
	return opts, nil
}

func buildOptionsFromOutputsConfig(outputBasePath string, configs []*OutputConfig) ([]newPipelineOption, error) {
	opts := []newPipelineOption{}
	for _, config := range configs {
		switch config.Type {
		case OutputStdout:
			opts = append(opts, WithOutputWriter(os.Stdout))
		case OutputFile:
			outputConfig, ok := config.Config.(*FileOutputConfig)
			if !ok {
				return nil, InvalidConfigTypeError{
					Expected: FileOutputConfig{},
					Actual:   config.Config,
				}
			}
			outputPath := outputConfig.Path
			if !filepath.IsAbs(outputPath) {
				outputPath = filepath.Join(outputBasePath, outputPath)
			}
			opts = append(opts, WithOutputFile(outputPath))
		default:
			return nil, InvalidOutputTypeError{
				Supplied: string(config.Type),
				ValidValues: lo.Map(validOutputTypes, func(t OutputType, _ int) string {
					return string(t)
				}),
			}
		}
	}
	return opts, nil
}
