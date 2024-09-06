// Package prefixwriter provides a writer that prefixes each line with a specified byte slice.
package prefixwriter

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"strings"
)

// Writer is a custom writer that prefixes each line with a specified byte slice.
// It implements the io.Writer interface.
type Writer struct {
	w writeFlusher

	head    bool
	prefix  []byte
	written int64
}

// Ensure Writer implements io.Writer interface.
var _ io.Writer = (*Writer)(nil)

// New creates a new Writer that wraps the given io.Writer and prefixes each line
// with the specified prefix. It uses the default buffer size.
//
// If w is already a *bufio.Writer with a large enough buffer, it is used directly.
// If w is a *bytes.Buffer or *strings.Builder, it is used without additional buffering.
func New(w io.Writer, prefix []byte) *Writer {
	return NewSize(w, prefix, 0)
}

// NewSize creates a new Writer that wraps the given io.Writer and prefixes each line
// with the specified prefix. It uses a buffer of at least the specified size.
//
// If w is already a *bufio.Writer with a buffer at least as large as the specified size,
// it is used directly. If w is a *bytes.Buffer or *strings.Builder, it is used without
// additional buffering, ignoring the size parameter. If size is <= 0, the default buffer
// size is used.
func NewSize(w io.Writer, prefix []byte, size int) *Writer {
	return &Writer{
		w:      wrapWriter(w, size),
		head:   true,
		prefix: slices.Clone(prefix),
	}
}

func wrapWriter(w io.Writer, size int) writeFlusher {
	switch w := w.(type) {
	case *bytes.Buffer:
		return nopFlusher{Writer: w}
	case *strings.Builder:
		return nopFlusher{Writer: w}
	}
	return bufio.NewWriterSize(w, size)
}

type writeFlusher interface {
	Write([]byte) (int, error)
	Flush() error
}

var _ writeFlusher = (*bufio.Writer)(nil)
var _ writeFlusher = nopFlusher{}

type nopFlusher struct {
	io.Writer
}

func (nopFlusher) Flush() error { return nil }

// Written returns the total number of bytes written, including prefixes.
func (w *Writer) Written() int64 {
	return w.written
}

// Write implements the io.Writer interface. It writes the given byte slice to the
// underlying writer, adding the prefix at the beginning of each line, including empty lines.
//
// The returned int n represents the number of bytes from the input slice p that were
// processed, not including any added prefixes. This means that n <= len(p), even though
// the actual number of bytes written to the underlying writer may be larger due to the
// added prefixes.
//
// If p contains no data (len(p) == 0), Write will not perform any operation and will
// return n = 0 and a nil error.
//
// A prefix is written before each line, including empty lines within p, but excluding
// a potential empty line at the end of p.
//
// An error is returned if the underlying writer returns an error, or if the Write
// operation cannot be completed fully.
func (w *Writer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	n := 0
	for len(p) > 0 {
		if w.head {
			pn, err := w.w.Write(w.prefix)
			w.written += int64(pn)
			if err != nil {
				return n, err
			}
			w.head = false
		}

		i := bytes.IndexByte(p, '\n')
		if i == -1 {
			m, err := w.w.Write(p)
			n += m
			w.written += int64(m)
			if err == nil {
				err = w.w.Flush()
			}
			return n, err
		}

		m, err := w.w.Write(p[:i+1])
		n += m
		w.written += int64(m)
		if err != nil {
			return n, err
		}

		p = p[i+1:]
		w.head = true
	}

	return n, w.w.Flush()
}
