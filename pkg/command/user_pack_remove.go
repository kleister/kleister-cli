package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userPackRemoveBind struct {
	ID   string
	Pack string
}

var (
	userPackRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove pack from user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userPackRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	userPackRemoveArgs = userPackRemoveBind{}
)

func init() {
	userPackCmd.AddCommand(userPackRemoveCmd)

	userPackRemoveCmd.Flags().StringVarP(
		&userPackRemoveArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userPackRemoveCmd.Flags().StringVar(
		&userPackRemoveArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)
}

func userPackRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userPackRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userPackRemoveArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	resp, err := client.DeleteUserFromPackWithResponse(
		ccmd.Context(),
		userPackRemoveArgs.ID,
		kleister.DeleteUserFromPackJSONRequestBody{
			Pack: userPackRemoveArgs.Pack,
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
