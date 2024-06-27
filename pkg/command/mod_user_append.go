package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modUserAppendBind struct {
	ID   string
	User string
	Perm string
}

var (
	modUserAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append user to mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modUserAppendAction)
		},
		Args: cobra.NoArgs,
	}

	modUserAppendArgs = modUserAppendBind{}
)

func init() {
	modUserCmd.AddCommand(modUserAppendCmd)

	modUserAppendCmd.Flags().StringVarP(
		&modUserAppendArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modUserAppendCmd.Flags().StringVar(
		&modUserAppendArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	modUserAppendCmd.Flags().StringVar(
		&modUserAppendArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func modUserAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modUserAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modUserAppendArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.AttachModToUserJSONRequestBody{
		User: modUserAppendArgs.User,
	}

	if modUserAppendArgs.Perm != "" {
		val, err := kleister.ToModUserParamsPerm(modUserAppendArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrModUserParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.AttachModToUserWithResponse(
		ccmd.Context(),
		modUserAppendArgs.ID,
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
