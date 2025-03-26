package command

import (
	"github.com/spf13/cobra"
)

var (
	groupPackCmd = &cobra.Command{
		Use:   "pack",
		Short: "Pack assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	groupCmd.AddCommand(groupPackCmd)
}
