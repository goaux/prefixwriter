package prefixwriter_test

import (
	"fmt"
	"os"

	"github.com/goaux/prefixwriter"
)

func Example() {
	w := prefixwriter.New(os.Stdout, []byte("001> "))
	fmt.Fprintln(w, "hello")
	fmt.Fprintln(w, "world")
	// output:
	// 001> hello
	// 001> world
}
