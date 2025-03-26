package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupModPermitBind struct {
	ID   string
	Mod  string
	Perm string
}

var (
	groupModPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit mod for group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupModPermitAction)
		},
		Args: cobra.NoArgs,
	}

	groupModPermitArgs = groupModPermitBind{}
)

func init() {
	groupModCmd.AddCommand(groupModPermitCmd)

	groupModPermitCmd.Flags().StringVarP(
		&groupModPermitArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupModPermitCmd.Flags().StringVar(
		&groupModPermitArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	groupModPermitCmd.Flags().StringVar(
		&groupModPermitArgs.Perm,
		"perm",
		"",
		"Role for the mod",
	)
}

func groupModPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupModPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupModPermitArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	body := kleister.PermitGroupModJSONRequestBody{
		Mod: groupModPermitArgs.Mod,
	}

	if groupModPermitArgs.Perm != "" {
		body.Perm = groupModPermitArgs.Perm
	}

	resp, err := client.PermitGroupModWithResponse(
		ccmd.Context(),
		groupModPermitArgs.ID,
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
