package prefixwriter_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/goaux/prefixwriter"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prefix   string
		expected string
	}{
		{"Empty input", "", "PREFIX: ", ""},
		{"Single line", "Hello, World!", "PREFIX: ", "PREFIX: Hello, World!"},
		{"Multiple lines", "Line 1\nLine 2\nLine 3", "-> ", "-> Line 1\n-> Line 2\n-> Line 3"},
		{"Empty lines", "\n\nContent\n\n", "# ", "# \n# \n# Content\n# \n"},
		{"Empty lines", "\n\n\n\n", "# ", "# \n# \n# \n# \n"},
		{"No newline at end", "Text without newline", "> ", "> Text without newline"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Buf interface {
				Write([]byte) (int, error)
				String() string
			}
			for _, buf := range []Buf{&bytes.Buffer{}, &strings.Builder{}} {
				w := prefixwriter.New(buf, []byte(tt.prefix))

				n, err := io.WriteString(w, tt.input)
				if err != nil {
					t.Fatalf("Write error: %v", err)
				}
				if n != len(tt.input) {
					t.Errorf("Write returned %d, want %d", n, len(tt.input))
				}

				if got := buf.String(); got != tt.expected {
					t.Errorf("Got %q, want %q", got, tt.expected)
				}

				if w.Written() != int64(len(tt.expected)) {
					t.Errorf("Written() returned %d, want %d", w.Written(), len(tt.expected))
				}
			}
		})
	}
}

func TestNewSize(t *testing.T) {
	buf := &bytes.Buffer{}
	w := prefixwriter.NewSize(buf, []byte("PREFIX: "), 1024)

	_, err := io.WriteString(w, "Test")
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if got := buf.String(); got != "PREFIX: Test" {
		t.Errorf("Got %q, want %q", got, "PREFIX: Test")
	}
}
