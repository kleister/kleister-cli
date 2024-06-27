package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	packTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packCmd.AddCommand(packTeamCmd)
}

func packTeamPerm(val string) kleister.PackTeamParamsPerm {
	switch val {
	case "owner":
		return kleister.PackTeamParamsPermOwner
	case "admin":
		return kleister.PackTeamParamsPermAdmin
	case "user":
		return kleister.PackTeamParamsPermUser
	}

	return kleister.PackTeamParamsPermUser
}
