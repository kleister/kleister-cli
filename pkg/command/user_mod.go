package command

import (
	"github.com/kleister/kleister-go/kleister"
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

func userModPerm(val string) kleister.UserModParamsPerm {
	switch val {
	case "owner":
		return kleister.UserModParamsPermOwner
	case "admin":
		return kleister.UserModParamsPermAdmin
	case "user":
		return kleister.UserModParamsPermUser
	}

	return kleister.UserModParamsPermUser
}
