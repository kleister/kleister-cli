package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupUserAppendBind struct {
	ID   string
	User string
	Perm string
}

var (
	groupUserAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append user to group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupUserAppendAction)
		},
		Args: cobra.NoArgs,
	}

	groupUserAppendArgs = groupUserAppendBind{}
)

func init() {
	groupUserCmd.AddCommand(groupUserAppendCmd)

	groupUserAppendCmd.Flags().StringVarP(
		&groupUserAppendArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupUserAppendCmd.Flags().StringVar(
		&groupUserAppendArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	groupUserAppendCmd.Flags().StringVar(
		&groupUserAppendArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func groupUserAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupUserAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupUserAppendArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.AttachGroupToUserJSONRequestBody{
		User: groupUserAppendArgs.User,
	}

	if groupUserAppendArgs.Perm != "" {
		body.Perm = groupUserAppendArgs.Perm
	}

	resp, err := client.AttachGroupToUserWithResponse(
		ccmd.Context(),
		groupUserAppendArgs.ID,
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
