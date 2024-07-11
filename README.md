# prefixwriter
Package prefixwriter provides a writer that prefixes each line with a specified byte slice.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/prefixwriter.svg)](https://pkg.go.dev/github.com/goaux/prefixwriter)

## Features

- Adds a specified prefix to the beginning of each line
- Implements the `io.Writer` interface
- Customizable buffer size
- Tracks total bytes written

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/goaux/prefixwriter"
)

func main() {
	w := prefixwriter.New(os.Stdout, []byte("LOG: "))
	fmt.Fprintln(w, "This is a log message")
	fmt.Fprintln(w, "Another log message")
}
```

OUTPUT:

```
LOG: This is a log message
LOG: Another log message
```
