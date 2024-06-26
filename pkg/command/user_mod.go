package command

import (
	"github.com/spf13/cobra"
)

var (
	userModCmd = &cobra.Command{
		Use:   "mod",
		Short: "Mod assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	userCmd.AddCommand(userModCmd)
}

// TODO user mod sub commands
