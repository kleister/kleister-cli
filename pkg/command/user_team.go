package command

import (
	"github.com/spf13/cobra"
)

var (
	userTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	userCmd.AddCommand(userTeamCmd)
}
