package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupModRemoveBind struct {
	ID  string
	Mod string
}

var (
	groupModRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove mod from group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupModRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	groupModRemoveArgs = groupModRemoveBind{}
)

func init() {
	groupModCmd.AddCommand(groupModRemoveCmd)

	groupModRemoveCmd.Flags().StringVarP(
		&groupModRemoveArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupModRemoveCmd.Flags().StringVar(
		&groupModRemoveArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)
}

func groupModRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupModRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if groupModRemoveArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	resp, err := client.DeleteGroupFromModWithResponse(
		ccmd.Context(),
		groupModRemoveArgs.ID,
		kleister.DeleteGroupFromModJSONRequestBody{
			Mod: groupModRemoveArgs.Mod,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
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
