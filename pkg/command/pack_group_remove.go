package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packGroupRemoveBind struct {
	ID    string
	Group string
}

var (
	packGroupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packGroupRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	packGroupRemoveArgs = packGroupRemoveBind{}
)

func init() {
	packGroupCmd.AddCommand(packGroupRemoveCmd)

	packGroupRemoveCmd.Flags().StringVarP(
		&packGroupRemoveArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packGroupRemoveCmd.Flags().StringVar(
		&packGroupRemoveArgs.Group,
		"group",
		"",
		"Group ID or slug",
	)
}

func packGroupRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packGroupRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packGroupRemoveArgs.Group == "" {
		return fmt.Errorf("you must provide a group ID or a slug")
	}

	resp, err := client.DeletePackFromGroupWithResponse(
		ccmd.Context(),
		packGroupRemoveArgs.ID,
		kleister.DeletePackFromGroupJSONRequestBody{
			Group: packGroupRemoveArgs.Group,
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
