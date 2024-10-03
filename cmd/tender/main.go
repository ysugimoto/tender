package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ysugimoto/tender"
)

func exitError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		exitError("Source template file must be provided.")
	}
	file := os.Args[1]
	if stat, err := os.Stat(file); err != nil {
		exitError("Source template file %s is not found.", file)
	} else if stat.IsDir() {
		exitError("Source template file %s is directory.", file)
	}
	fp, err := os.Open(file)
	if err != nil {
		exitError("Failed to open file %s.", file)
	}
	defer fp.Close()

	rendered, err := tender.New(fp).Render()
	if err != nil {
		exitError("Failed to execute template: %s", err.Error())
	}

	io.WriteString(os.Stdout, rendered) // nolint:errcheck
}
