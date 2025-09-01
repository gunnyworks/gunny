package gunny

import (
	"context"
	"io"

	"github.com/alexkappa/mustache"
)

// Renderer provides a way to render data to a specific output writer.
type Renderer interface {
	Render(ctx context.Context, data DataResolver, w io.Writer) error
}

type NullRenderer struct{}

var _ Renderer = (*NullRenderer)(nil)

// Render implements Renderer.
func (n *NullRenderer) Render(ctx context.Context, data DataResolver, w io.Writer) error {
	logger := LoggerFromContext(ctx)
	logger.Debug("Using null renderer")
	// Do nothing
	return nil
}

// MustacheTemplateRenderer renders data through [Mustache] templates.
//
// [Mustache]: http://mustache.github.io
type MustacheTemplateRenderer struct {
	tmpl *mustache.Template
}

var _ Renderer = (*MustacheTemplateRenderer)(nil)

// NewMustacheTemplateRenderer creates a new [MustacheTemplateRenderer] whose
// template is loaded from the given reader.
func NewMustacheTemplateRenderer(r io.Reader) (*MustacheTemplateRenderer, error) {
	tmpl, err := mustache.Parse(r)
	if err != nil {
		return nil, ErrTemplateRead{Cause: err}
	}
	return &MustacheTemplateRenderer{
		tmpl: tmpl,
	}, nil
}

// Render implements Renderer.
func (r *MustacheTemplateRenderer) Render(ctx context.Context, data DataResolver, w io.Writer) error {
	logger := LoggerFromContext(ctx)
	logger.Debug("Using Mustache template renderer")
	// We only attempt to resolve data as we are about to render.
	resolvedData, err := data.Resolve(ctx)
	if err != nil {
		return err
	}
	logger.Debug("Resolved data", "data", resolvedData)
	return r.tmpl.Render(w, resolvedData)
}
