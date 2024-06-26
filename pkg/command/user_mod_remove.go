package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userModRemoveBind struct {
	ID  string
	Mod string
}

var (
	userModRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove mod from user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userModRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	userModRemoveArgs = userModRemoveBind{}
)

func init() {
	userModCmd.AddCommand(userModRemoveCmd)

	userModRemoveCmd.Flags().StringVarP(
		&userModRemoveArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userModRemoveCmd.Flags().StringVar(
		&userModRemoveArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)
}

func userModRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userModRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userModRemoveArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	resp, err := client.DeleteUserFromModWithResponse(
		ccmd.Context(),
		userModRemoveArgs.ID,
		kleister.DeleteUserFromModJSONRequestBody{
			Mod: userModRemoveArgs.Mod,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusPreconditionFailed:
		return fmt.Errorf(kleister.FromPtr(resp.JSON412.Message))
	case http.StatusForbidden:
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusNotFound:
		return fmt.Errorf(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
