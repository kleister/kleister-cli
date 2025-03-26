package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupUserPermitBind struct {
	ID   string
	User string
	Perm string
}

var (
	groupUserPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit user for group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupUserPermitAction)
		},
		Args: cobra.NoArgs,
	}

	groupUserPermitArgs = groupUserPermitBind{}
)

func init() {
	groupUserCmd.AddCommand(groupUserPermitCmd)

	groupUserPermitCmd.Flags().StringVarP(
		&groupUserPermitArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupUserPermitCmd.Flags().StringVar(
		&groupUserPermitArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	groupUserPermitCmd.Flags().StringVar(
		&groupUserPermitArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func groupUserPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupUserPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupUserPermitArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.PermitGroupUserJSONRequestBody{
		User: groupUserPermitArgs.User,
	}

	if groupUserPermitArgs.Perm != "" {
		body.Perm = groupUserPermitArgs.Perm
	}

	resp, err := client.PermitGroupUserWithResponse(
		ccmd.Context(),
		groupUserPermitArgs.ID,
		body,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusUnprocessableEntity:
		return validationError(resp.JSON422)
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
