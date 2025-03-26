package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modGroupPermitBind struct {
	ID    string
	Group string
	Perm  string
}

var (
	modGroupPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit group for mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modGroupPermitAction)
		},
		Args: cobra.NoArgs,
	}

	modGroupPermitArgs = modGroupPermitBind{}
)

func init() {
	modGroupCmd.AddCommand(modGroupPermitCmd)

	modGroupPermitCmd.Flags().StringVarP(
		&modGroupPermitArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modGroupPermitCmd.Flags().StringVar(
		&modGroupPermitArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)

	modGroupPermitCmd.Flags().StringVar(
		&modGroupPermitArgs.Perm,
		"perm",
		"",
		"Role for the group",
	)
}

func modGroupPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modGroupPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modGroupPermitArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	body := kleister.PermitModGroupJSONRequestBody{
		Group: modGroupPermitArgs.Group,
	}

	if modGroupPermitArgs.Perm != "" {
		body.Perm = modGroupPermitArgs.Perm
	}

	resp, err := client.PermitModGroupWithResponse(
		ccmd.Context(),
		modGroupPermitArgs.ID,
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
