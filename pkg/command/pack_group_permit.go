package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packGroupPermitBind struct {
	ID    string
	Group string
	Perm  string
}

var (
	packGroupPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit group for pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packGroupPermitAction)
		},
		Args: cobra.NoArgs,
	}

	packGroupPermitArgs = packGroupPermitBind{}
)

func init() {
	packGroupCmd.AddCommand(packGroupPermitCmd)

	packGroupPermitCmd.Flags().StringVarP(
		&packGroupPermitArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packGroupPermitCmd.Flags().StringVar(
		&packGroupPermitArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)

	packGroupPermitCmd.Flags().StringVar(
		&packGroupPermitArgs.Perm,
		"perm",
		"",
		"Role for the group",
	)
}

func packGroupPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packGroupPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packGroupPermitArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	body := kleister.PermitPackGroupJSONRequestBody{
		Group: packGroupPermitArgs.Group,
	}

	if packGroupPermitArgs.Perm != "" {
		body.Perm = packGroupPermitArgs.Perm
	}

	resp, err := client.PermitPackGroupWithResponse(
		ccmd.Context(),
		packGroupPermitArgs.ID,
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
