package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupPackAppendBind struct {
	ID   string
	Pack string
	Perm string
}

var (
	groupPackAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append pack to group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupPackAppendAction)
		},
		Args: cobra.NoArgs,
	}

	groupPackAppendArgs = groupPackAppendBind{}
)

func init() {
	groupPackCmd.AddCommand(groupPackAppendCmd)

	groupPackAppendCmd.Flags().StringVarP(
		&groupPackAppendArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupPackAppendCmd.Flags().StringVar(
		&groupPackAppendArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	groupPackAppendCmd.Flags().StringVar(
		&groupPackAppendArgs.Perm,
		"perm",
		"",
		"Role for the pack",
	)
}

func groupPackAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupPackAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupPackAppendArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	body := kleister.AttachGroupToPackJSONRequestBody{
		Pack: groupPackAppendArgs.Pack,
	}

	if groupPackAppendArgs.Perm != "" {
		body.Perm = groupPackAppendArgs.Perm
	}

	resp, err := client.AttachGroupToPackWithResponse(
		ccmd.Context(),
		groupPackAppendArgs.ID,
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
