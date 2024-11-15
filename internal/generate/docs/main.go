package main

import (
	"os"

	"gabe565.com/domain-watch/cmd"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra/doc"
)

func main() {
	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		panic(err)
	}

	rootCmd := cmd.New(cobrax.WithVersion("beta"))
	if err := doc.GenMarkdownTree(rootCmd, output); err != nil {
		panic(err)
	}
}
