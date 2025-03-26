package command

import (
	"github.com/spf13/cobra"
)

var (
	groupUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	groupCmd.AddCommand(groupUserCmd)
}
