package command

import (
	"github.com/spf13/cobra"
)

var (
	packUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packCmd.AddCommand(packUserCmd)
}

// TODO pack user sub commands
