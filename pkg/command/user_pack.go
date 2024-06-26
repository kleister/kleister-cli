package command

import (
	"github.com/spf13/cobra"
)

var (
	userPackCmd = &cobra.Command{
		Use:   "pack",
		Short: "Pack assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	userCmd.AddCommand(userPackCmd)
}

// TODO user pack sub commands
