package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modDeleteBind struct {
	ID string
}

var (
	modDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	modDeleteArgs = modDeleteBind{}
)

func init() {
	modCmd.AddCommand(modDeleteCmd)

	modDeleteCmd.Flags().StringVarP(
		&modDeleteArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)
}

func modDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modDeleteArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.DeleteModWithResponse(
		ccmd.Context(),
		modDeleteArgs.ID,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully deleted")
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
