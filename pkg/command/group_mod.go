package command

import (
	"github.com/spf13/cobra"
)

var (
	groupModCmd = &cobra.Command{
		Use:   "mod",
		Short: "Mod assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	groupCmd.AddCommand(groupModCmd)
}
