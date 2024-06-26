package command

import (
	"github.com/spf13/cobra"
)

var (
	modCmd = &cobra.Command{
		Use:   "mod",
		Short: "Mod related sub-commands",
		Args:  cobra.NoArgs,
	}
)

func init() {
	rootCmd.AddCommand(modCmd)
}
