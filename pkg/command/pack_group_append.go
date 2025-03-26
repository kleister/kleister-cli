package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packGroupAppendBind struct {
	ID    string
	Group string
	Perm  string
}

var (
	packGroupAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append group to pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packGroupAppendAction)
		},
		Args: cobra.NoArgs,
	}

	packGroupAppendArgs = packGroupAppendBind{}
)

func init() {
	packGroupCmd.AddCommand(packGroupAppendCmd)

	packGroupAppendCmd.Flags().StringVarP(
		&packGroupAppendArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packGroupAppendCmd.Flags().StringVar(
		&packGroupAppendArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)

	packGroupAppendCmd.Flags().StringVar(
		&packGroupAppendArgs.Perm,
		"perm",
		"",
		"Role for the group",
	)
}

func packGroupAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packGroupAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packGroupAppendArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	body := kleister.AttachPackToGroupJSONRequestBody{
		Group: packGroupAppendArgs.Group,
	}

	if packGroupAppendArgs.Perm != "" {
		body.Perm = packGroupAppendArgs.Perm
	}

	resp, err := client.AttachPackToGroupWithResponse(
		ccmd.Context(),
		packGroupAppendArgs.ID,
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
