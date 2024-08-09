package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func completion(cmd *cobra.Command, shell string) error {
	switch shell {
	case "bash":
		return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
	case "zsh":
		return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
	case "fish":
		return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
	default:
		return fmt.Errorf("%v: invalid shell", shell)
	}
}
