package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modUserRemoveBind struct {
	ID   string
	User string
}

var (
	modUserRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modUserRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	modUserRemoveArgs = modUserRemoveBind{}
)

func init() {
	modUserCmd.AddCommand(modUserRemoveCmd)

	modUserRemoveCmd.Flags().StringVarP(
		&modUserRemoveArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modUserRemoveCmd.Flags().StringVar(
		&modUserRemoveArgs.User,
		"user",
		"",
		"User ID or slug",
	)
}

func modUserRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modUserRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modUserRemoveArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	resp, err := client.DeleteModFromUserWithResponse(
		ccmd.Context(),
		modUserRemoveArgs.ID,
		kleister.DeleteModFromUserJSONRequestBody{
			User: modUserRemoveArgs.User,
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
