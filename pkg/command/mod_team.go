package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	modTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modCmd.AddCommand(modTeamCmd)
}

func modTeamPerm(val string) kleister.ModTeamParamsPerm {
	switch val {
	case "owner":
		return kleister.ModTeamParamsPermOwner
	case "admin":
		return kleister.ModTeamParamsPermAdmin
	case "user":
		return kleister.ModTeamParamsPermUser
	}

	return kleister.ModTeamParamsPermUser
}
