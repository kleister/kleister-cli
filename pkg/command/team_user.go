package command

import (
	"github.com/spf13/cobra"
)

var (
	teamUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	teamCmd.AddCommand(teamUserCmd)
}
