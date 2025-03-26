package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userGroupPermitBind struct {
	ID    string
	Group string
	Perm  string
}

var (
	userGroupPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit group for user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userGroupPermitAction)
		},
		Args: cobra.NoArgs,
	}

	userGroupPermitArgs = userGroupPermitBind{}
)

func init() {
	userGroupCmd.AddCommand(userGroupPermitCmd)

	userGroupPermitCmd.Flags().StringVarP(
		&userGroupPermitArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userGroupPermitCmd.Flags().StringVar(
		&userGroupPermitArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)

	userGroupPermitCmd.Flags().StringVar(
		&userGroupPermitArgs.Perm,
		"perm",
		"",
		"Role for the group",
	)
}

func userGroupPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userGroupPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userGroupPermitArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	body := kleister.PermitUserGroupJSONRequestBody{
		Group: userGroupPermitArgs.Group,
	}

	if userGroupPermitArgs.Perm != "" {
		body.Perm = userGroupPermitArgs.Perm
	}

	resp, err := client.PermitUserGroupWithResponse(
		ccmd.Context(),
		userGroupPermitArgs.ID,
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
