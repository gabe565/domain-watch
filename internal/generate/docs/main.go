package main

import (
	"fmt"
	"log"
	"os"

	"gabe565.com/domain-watch/cmd"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra/doc"
)

func main() {
	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		log.Fatal(fmt.Errorf("failed to remove existing dia: %w", err))
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		log.Fatal(fmt.Errorf("failed to mkdir: %w", err))
	}

	rootCmd := cmd.New(cobrax.WithVersion("beta"))
	if err := doc.GenMarkdownTree(rootCmd, output); err != nil {
		log.Fatal(fmt.Errorf("failed to generate markdown: %w", err))
	}
}
