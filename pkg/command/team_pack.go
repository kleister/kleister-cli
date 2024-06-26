package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	teamPackCmd = &cobra.Command{
		Use:   "pack",
		Short: "Pack assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	teamCmd.AddCommand(teamPackCmd)
}

func teamPackPerm(val string) kleister.TeamPackParamsPerm {
	switch val {
	case "owner":
		return kleister.TeamPackParamsPermOwner
	case "admin":
		return kleister.TeamPackParamsPermAdmin
	case "user":
		return kleister.TeamPackParamsPermUser
	}

	return kleister.TeamPackParamsPermUser
}
