package command

import (
	"github.com/spf13/cobra"
)

var (
	packCmd = &cobra.Command{
		Use:   "pack",
		Short: "Pack related sub-commands",
		Args:  cobra.NoArgs,
	}
)

func init() {
	rootCmd.AddCommand(packCmd)
}
