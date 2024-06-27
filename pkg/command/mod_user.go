package command

import (
	"github.com/spf13/cobra"
)

var (
	modUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modCmd.AddCommand(modUserCmd)
}
