package command

import (
	"github.com/spf13/cobra"
)

var (
	minecraftCmd = &cobra.Command{
		Use:   "minecraft",
		Short: "Minecraft related sub-commands",
		Args:  cobra.NoArgs,
	}
)

func init() {
	rootCmd.AddCommand(minecraftCmd)
}
