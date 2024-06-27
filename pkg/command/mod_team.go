package command

import (
	"github.com/spf13/cobra"
)

var (
	modTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modCmd.AddCommand(modTeamCmd)
}
