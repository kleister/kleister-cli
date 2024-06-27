package command

import (
	"github.com/spf13/cobra"
)

var (
	packTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packCmd.AddCommand(packTeamCmd)
}
