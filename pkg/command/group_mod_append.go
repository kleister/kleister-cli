package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupModAppendBind struct {
	ID   string
	Mod  string
	Perm string
}

var (
	groupModAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append mod to group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupModAppendAction)
		},
		Args: cobra.NoArgs,
	}

	groupModAppendArgs = groupModAppendBind{}
)

func init() {
	groupModCmd.AddCommand(groupModAppendCmd)

	groupModAppendCmd.Flags().StringVarP(
		&groupModAppendArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupModAppendCmd.Flags().StringVar(
		&groupModAppendArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	groupModAppendCmd.Flags().StringVar(
		&groupModAppendArgs.Perm,
		"perm",
		"",
		"Role for the mod",
	)
}

func groupModAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupModAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupModAppendArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	body := kleister.AttachGroupToModJSONRequestBody{
		Mod: groupModAppendArgs.Mod,
	}

	if groupModAppendArgs.Perm != "" {
		body.Perm = groupModAppendArgs.Perm
	}

	resp, err := client.AttachGroupToModWithResponse(
		ccmd.Context(),
		groupModAppendArgs.ID,
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
