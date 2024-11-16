package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"gabe565.com/domain-watch/cmd"
	"gabe565.com/utils/cobrax"
)

func main() {
	if err := os.RemoveAll("completions"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("completions", 0o777); err != nil {
		panic(err)
	}

	rootCmd := cmd.New()
	name := rootCmd.Name()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	for _, shell := range []cobrax.Shell{cobrax.Bash, cobrax.Zsh, cobrax.Fish} {
		if err := cobrax.GenCompletion(rootCmd, shell); err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join("completions", name+"."+string(shell)))
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(f, &buf); err != nil {
			panic(err)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}
