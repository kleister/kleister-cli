package command

import (
	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	modUserCmd = &cobra.Command{
		Use:   "user",
		Short: "User assignments",
		Args:  cobra.NoArgs,
	}
)

func init() {
	modCmd.AddCommand(modUserCmd)
}

func modUserPerm(val string) kleister.ModUserParamsPerm {
	switch val {
	case "owner":
		return kleister.ModUserParamsPermOwner
	case "admin":
		return kleister.ModUserParamsPermAdmin
	case "user":
		return kleister.ModUserParamsPermUser
	}

	return kleister.ModUserParamsPermUser
}
