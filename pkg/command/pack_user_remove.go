package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packUserRemoveBind struct {
	ID   string
	User string
}

var (
	packUserRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packUserRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	packUserRemoveArgs = packUserRemoveBind{}
)

func init() {
	packUserCmd.AddCommand(packUserRemoveCmd)

	packUserRemoveCmd.Flags().StringVarP(
		&packUserRemoveArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packUserRemoveCmd.Flags().StringVar(
		&packUserRemoveArgs.User,
		"user",
		"",
		"User ID or slug",
	)
}

func packUserRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packUserRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packUserRemoveArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	resp, err := client.DeletePackFromUserWithResponse(
		ccmd.Context(),
		packUserRemoveArgs.ID,
		kleister.DeletePackFromUserJSONRequestBody{
			User: packUserRemoveArgs.User,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusPreconditionFailed:
		return fmt.Errorf(kleister.FromPtr(resp.JSON412.Message))
	case http.StatusForbidden:
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusNotFound:
		return fmt.Errorf(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
