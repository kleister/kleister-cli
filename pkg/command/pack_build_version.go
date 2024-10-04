package command

import (
	"github.com/spf13/cobra"
)

var (
	packBuildVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packBuildCmd.AddCommand(packBuildVersionCmd)
}

// TODO pack build version sub commands
