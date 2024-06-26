package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	teamModCmd = &cobra.Command{
		Use:   "mod",
		Short: "Mod assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	teamCmd.AddCommand(teamModCmd)
}

func teamModPerm(val string) kleister.TeamModParamsPerm {
	switch val {
	case "owner":
		return kleister.TeamModParamsPermOwner
	case "admin":
		return kleister.TeamModParamsPermAdmin
	case "user":
		return kleister.TeamModParamsPermUser
	}

	return kleister.TeamModParamsPermUser
}
