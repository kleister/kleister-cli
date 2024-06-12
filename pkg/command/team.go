package command

import (
	"github.com/spf13/cobra"
)

var (
	teamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team related sub-commands",
		Args:  cobra.NoArgs,
	}
)

func init() {
	rootCmd.AddCommand(teamCmd)
}
