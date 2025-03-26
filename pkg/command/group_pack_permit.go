package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupPackPermitBind struct {
	ID   string
	Pack string
	Perm string
}

var (
	groupPackPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit pack for group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupPackPermitAction)
		},
		Args: cobra.NoArgs,
	}

	groupPackPermitArgs = groupPackPermitBind{}
)

func init() {
	groupPackCmd.AddCommand(groupPackPermitCmd)

	groupPackPermitCmd.Flags().StringVarP(
		&groupPackPermitArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupPackPermitCmd.Flags().StringVar(
		&groupPackPermitArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	groupPackPermitCmd.Flags().StringVar(
		&groupPackPermitArgs.Perm,
		"perm",
		"",
		"Role for the pack",
	)
}

func groupPackPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupPackPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupPackPermitArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	body := kleister.PermitGroupPackJSONRequestBody{
		Pack: groupPackPermitArgs.Pack,
	}

	if groupPackPermitArgs.Perm != "" {
		body.Perm = groupPackPermitArgs.Perm
	}

	resp, err := client.PermitGroupPackWithResponse(
		ccmd.Context(),
		groupPackPermitArgs.ID,
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
