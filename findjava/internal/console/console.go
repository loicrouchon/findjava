package console

import (
	"fmt"
	"io"
	"os"
)

type writer struct {
	Stdout io.Writer
	Stderr io.Writer
}

var Writer = writer{
	Stdout: os.Stdout,
	Stderr: os.Stderr,
}

func (console *writer) Printf(message string, v ...interface{}) {
	_, _ = fmt.Fprintf(console.Stdout, message, v...)
}

func (console *writer) Eprintf(message string, v ...interface{}) {
	_, _ = fmt.Fprintf(console.Stderr, message, v...)
}
