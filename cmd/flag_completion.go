package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var completionFlag string

func init() {
	Command.Flags().StringVar(&completionFlag, "completion", "", "Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.")
	err := Command.RegisterFlagCompletionFunc("completion", completionCompletion)
	if err != nil {
		panic(err)
	}
}

func completionCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
}

var completionWriter io.Writer = os.Stdout

func completion(cmd *cobra.Command, args []string) error {
	switch completionFlag {
	case "bash":
		if err := cmd.Root().GenBashCompletion(completionWriter); err != nil {
			return err
		}
	case "zsh":
		if err := cmd.Root().GenZshCompletion(completionWriter); err != nil {
			return err
		}
	case "fish":
		if err := cmd.Root().GenFishCompletion(completionWriter, true); err != nil {
			return err
		}
	case "powershell":
		if err := cmd.Root().GenPowerShellCompletionWithDesc(completionWriter); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%v: invalid shell", completionFlag)
	}
	return nil
}
