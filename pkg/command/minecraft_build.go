package command

import (
	"github.com/spf13/cobra"
)

var (
	minecraftBuildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	minecraftCmd.AddCommand(minecraftBuildCmd)
}
