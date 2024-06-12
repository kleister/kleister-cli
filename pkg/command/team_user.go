package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	teamUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	teamCmd.AddCommand(teamUserCmd)
}

func teamUserPerm(val string) kleister.TeamUserParamsPerm {
	switch val {
	case "owner":
		return kleister.TeamUserParamsPermOwner
	case "admin":
		return kleister.TeamUserParamsPermAdmin
	case "user":
		return kleister.TeamUserParamsPermUser
	}

	return kleister.TeamUserParamsPermUser
}
