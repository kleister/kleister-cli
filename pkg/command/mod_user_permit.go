package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modUserPermitBind struct {
	ID   string
	User string
	Perm string
}

var (
	modUserPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit user for mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modUserPermitAction)
		},
		Args: cobra.NoArgs,
	}

	modUserPermitArgs = modUserPermitBind{}
)

func init() {
	modUserCmd.AddCommand(modUserPermitCmd)

	modUserPermitCmd.Flags().StringVarP(
		&modUserPermitArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modUserPermitCmd.Flags().StringVar(
		&modUserPermitArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	modUserPermitCmd.Flags().StringVar(
		&modUserPermitArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func modUserPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modUserPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modUserPermitArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.PermitModUserJSONRequestBody{
		User: modUserPermitArgs.User,
	}

	if modUserPermitArgs.Perm != "" {
		body.Perm = kleister.ToPtr(modUserPerm(modUserPermitArgs.Perm))
	}

	resp, err := client.PermitModUserWithResponse(
		ccmd.Context(),
		modUserPermitArgs.ID,
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
