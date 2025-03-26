package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userGroupRemoveBind struct {
	ID    string
	Group string
}

var (
	userGroupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userGroupRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	userGroupRemoveArgs = userGroupRemoveBind{}
)

func init() {
	userGroupCmd.AddCommand(userGroupRemoveCmd)

	userGroupRemoveCmd.Flags().StringVarP(
		&userGroupRemoveArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userGroupRemoveCmd.Flags().StringVar(
		&userGroupRemoveArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)
}

func userGroupRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userGroupRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userGroupRemoveArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	resp, err := client.DeleteUserFromGroupWithResponse(
		ccmd.Context(),
		userGroupRemoveArgs.ID,
		kleister.DeleteUserFromGroupJSONRequestBody{
			Group: userGroupRemoveArgs.Group,
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
