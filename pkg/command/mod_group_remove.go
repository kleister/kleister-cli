package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modGroupRemoveBind struct {
	ID    string
	Group string
}

var (
	modGroupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modGroupRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	modGroupRemoveArgs = modGroupRemoveBind{}
)

func init() {
	modGroupCmd.AddCommand(modGroupRemoveCmd)

	modGroupRemoveCmd.Flags().StringVarP(
		&modGroupRemoveArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modGroupRemoveCmd.Flags().StringVar(
		&modGroupRemoveArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)
}

func modGroupRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modGroupRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modGroupRemoveArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	resp, err := client.DeleteModFromGroupWithResponse(
		ccmd.Context(),
		modGroupRemoveArgs.ID,
		kleister.DeleteModFromGroupJSONRequestBody{
			Group: modGroupRemoveArgs.Group,
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
