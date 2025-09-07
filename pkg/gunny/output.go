package gunny

import (
	"io"
)

// writeCloserNoopClose implements [io.WriteCloser], with the Close method
// being a noop.
type writeCloserNoopClose struct {
	w io.Writer
}

var _ io.WriteCloser = (*writeCloserNoopClose)(nil)

// Write implements io.WriteCloser.
func (w *writeCloserNoopClose) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// Close implements io.WriteCloser.
func (w *writeCloserNoopClose) Close() error {
	return nil
}
