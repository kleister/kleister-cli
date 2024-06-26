package command

import (
	"github.com/spf13/cobra"
)

var (
	packBuildCmd = &cobra.Command{
		Use:   "build",
		Short: "Builds for the pack",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packCmd.AddCommand(packBuildCmd)
}

// TODO pack build sub commands
