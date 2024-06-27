package command

import (
	"github.com/spf13/cobra"
)

var (
	teamModCmd = &cobra.Command{
		Use:   "mod",
		Short: "Mod assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	teamCmd.AddCommand(teamModCmd)
}
