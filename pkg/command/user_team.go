package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	userTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	userCmd.AddCommand(userTeamCmd)
}

func userTeamPerm(val string) kleister.UserTeamParamsPerm {
	switch val {
	case "owner":
		return kleister.UserTeamParamsPermOwner
	case "admin":
		return kleister.UserTeamParamsPermAdmin
	case "user":
		return kleister.UserTeamParamsPermUser
	}

	return kleister.UserTeamParamsPermUser
}
