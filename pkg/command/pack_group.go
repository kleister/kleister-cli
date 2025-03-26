package command

import (
	"github.com/spf13/cobra"
)

var (
	packGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Group assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packCmd.AddCommand(packGroupCmd)
}
