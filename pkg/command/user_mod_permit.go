package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userModPermitBind struct {
	ID   string
	Mod  string
	Perm string
}

var (
	userModPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit mod for user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userModPermitAction)
		},
		Args: cobra.NoArgs,
	}

	userModPermitArgs = userModPermitBind{}
)

func init() {
	userModCmd.AddCommand(userModPermitCmd)

	userModPermitCmd.Flags().StringVarP(
		&userModPermitArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userModPermitCmd.Flags().StringVar(
		&userModPermitArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	userModPermitCmd.Flags().StringVar(
		&userModPermitArgs.Perm,
		"perm",
		"",
		"Role for the mod",
	)
}

func userModPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userModPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userModPermitArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	body := kleister.PermitUserModJSONRequestBody{
		Mod: userModPermitArgs.Mod,
	}

	if userModPermitArgs.Perm != "" {
		body.Perm = kleister.ToPtr(userModPerm(userModPermitArgs.Perm))
	}

	resp, err := client.PermitUserModWithResponse(
		ccmd.Context(),
		userModPermitArgs.ID,
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
