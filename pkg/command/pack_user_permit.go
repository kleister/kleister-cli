package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packUserPermitBind struct {
	ID   string
	User string
	Perm string
}

var (
	packUserPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit user for pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packUserPermitAction)
		},
		Args: cobra.NoArgs,
	}

	packUserPermitArgs = packUserPermitBind{}
)

func init() {
	packUserCmd.AddCommand(packUserPermitCmd)

	packUserPermitCmd.Flags().StringVarP(
		&packUserPermitArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packUserPermitCmd.Flags().StringVar(
		&packUserPermitArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	packUserPermitCmd.Flags().StringVar(
		&packUserPermitArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func packUserPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packUserPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packUserPermitArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.PermitPackUserJSONRequestBody{
		User: packUserPermitArgs.User,
	}

	if packUserPermitArgs.Perm != "" {
		val, err := kleister.ToPackUserParamsPerm(packUserPermitArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrPackUserParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.PermitPackUserWithResponse(
		ccmd.Context(),
		packUserPermitArgs.ID,
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
