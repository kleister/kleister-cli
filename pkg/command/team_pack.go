package command

import (
	"github.com/spf13/cobra"
)

var (
	teamPackCmd = &cobra.Command{
		Use:   "pack",
		Short: "Pack assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	teamCmd.AddCommand(teamPackCmd)
}
