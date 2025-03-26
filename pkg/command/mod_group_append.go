package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modGroupAppendBind struct {
	ID    string
	Group string
	Perm  string
}

var (
	modGroupAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append group to mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modGroupAppendAction)
		},
		Args: cobra.NoArgs,
	}

	modGroupAppendArgs = modGroupAppendBind{}
)

func init() {
	modGroupCmd.AddCommand(modGroupAppendCmd)

	modGroupAppendCmd.Flags().StringVarP(
		&modGroupAppendArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modGroupAppendCmd.Flags().StringVar(
		&modGroupAppendArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)

	modGroupAppendCmd.Flags().StringVar(
		&modGroupAppendArgs.Perm,
		"perm",
		"",
		"Role for the group",
	)
}

func modGroupAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modGroupAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modGroupAppendArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	body := kleister.AttachModToGroupJSONRequestBody{
		Group: modGroupAppendArgs.Group,
	}

	if modGroupAppendArgs.Perm != "" {
		body.Perm = modGroupAppendArgs.Perm
	}

	resp, err := client.AttachModToGroupWithResponse(
		ccmd.Context(),
		modGroupAppendArgs.ID,
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
