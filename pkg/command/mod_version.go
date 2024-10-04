package command

import (
	"github.com/spf13/cobra"
)

var (
	modVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Versions for the mod",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modCmd.AddCommand(modVersionCmd)
}
