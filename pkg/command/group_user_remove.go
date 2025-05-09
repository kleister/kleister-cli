package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupUserRemoveBind struct {
	ID   string
	User string
}

var (
	groupUserRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupUserRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	groupUserRemoveArgs = groupUserRemoveBind{}
)

func init() {
	groupUserCmd.AddCommand(groupUserRemoveCmd)

	groupUserRemoveCmd.Flags().StringVarP(
		&groupUserRemoveArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupUserRemoveCmd.Flags().StringVar(
		&groupUserRemoveArgs.User,
		"user",
		"",
		"User ID or slug",
	)
}

func groupUserRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupUserRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupUserRemoveArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	resp, err := client.DeleteGroupFromUserWithResponse(
		ccmd.Context(),
		groupUserRemoveArgs.ID,
		kleister.DeleteGroupFromUserJSONRequestBody{
			User: groupUserRemoveArgs.User,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusPreconditionFailed:
		return errors.New(kleister.FromPtr(resp.JSON412.Message))
	case http.StatusForbidden:
		return errors.New(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusNotFound:
		return errors.New(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
