package command

import (
	"github.com/spf13/cobra"
)

var (
	userGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Group assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	userCmd.AddCommand(userGroupCmd)
}
