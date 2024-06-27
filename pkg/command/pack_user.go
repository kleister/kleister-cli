package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	packUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	packCmd.AddCommand(packUserCmd)
}

func packUserPerm(val string) kleister.PackUserParamsPerm {
	switch val {
	case "owner":
		return kleister.PackUserParamsPermOwner
	case "admin":
		return kleister.PackUserParamsPermAdmin
	case "user":
		return kleister.PackUserParamsPermUser
	}

	return kleister.PackUserParamsPermUser
}
