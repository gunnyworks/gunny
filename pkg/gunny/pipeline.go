package gunny

import (
	"context"
	"io"
	"os"
)

// Pipeline encapsulates the configuration of a Gunny rendering pipeline.
type Pipeline struct {
	renderer  Renderer
	resolvers DataResolverMap
	writer    io.Writer
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

// WithOutputWriter allows one to specify where to write rendered output.
func WithOutputWriter(w io.Writer) newPipelineOption {
	return func(pipeline *Pipeline) error {
		pipeline.writer = w
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
		writer:    os.Stdout,
		logger:    &SlogLogger{},
	}
	for _, opt := range opts {
		if err := opt(pipeline); err != nil {
			return nil, err
		}
	}
	return pipeline, nil
}

// Render executes the entire pipeline rendering operation fully.
func (p *Pipeline) Render(ctx context.Context) error {
	p.logger.Debug("Rendering")
	ctxWithLogger := NewContextWithLogger(ctx, p.logger)
	return p.renderer.Render(ctxWithLogger, p.resolvers, p.writer)
}
