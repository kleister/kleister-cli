package command

import (
	"github.com/spf13/cobra"
)

var (
	modVersionBuildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modVersionCmd.AddCommand(modVersionBuildCmd)
}

// TODO mod version build sub commands
