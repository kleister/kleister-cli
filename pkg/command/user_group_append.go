package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userGroupAppendBind struct {
	ID    string
	Group string
	Perm  string
}

var (
	userGroupAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append group to user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userGroupAppendAction)
		},
		Args: cobra.NoArgs,
	}

	userGroupAppendArgs = userGroupAppendBind{}
)

func init() {
	userGroupCmd.AddCommand(userGroupAppendCmd)

	userGroupAppendCmd.Flags().StringVarP(
		&userGroupAppendArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userGroupAppendCmd.Flags().StringVar(
		&userGroupAppendArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)

	userGroupAppendCmd.Flags().StringVar(
		&userGroupAppendArgs.Perm,
		"perm",
		"",
		"Role for the group",
	)
}

func userGroupAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userGroupAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userGroupAppendArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	body := kleister.AttachUserToGroupJSONRequestBody{
		Group: userGroupAppendArgs.Group,
	}

	if userGroupAppendArgs.Perm != "" {
		body.Perm = userGroupAppendArgs.Perm
	}

	resp, err := client.AttachUserToGroupWithResponse(
		ccmd.Context(),
		userGroupAppendArgs.ID,
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
