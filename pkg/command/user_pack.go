package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	userPackCmd = &cobra.Command{
		Use:   "pack",
		Short: "Pack assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	userCmd.AddCommand(userPackCmd)
}

func userPackPerm(val string) kleister.UserPackParamsPerm {
	switch val {
	case "owner":
		return kleister.UserPackParamsPermOwner
	case "admin":
		return kleister.UserPackParamsPermAdmin
	case "user":
		return kleister.UserPackParamsPermUser
	}

	return kleister.UserPackParamsPermUser
}
