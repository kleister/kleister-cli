package command

import (
	"github.com/spf13/cobra"
)

var (
	modGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Group assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modCmd.AddCommand(modGroupCmd)
}
