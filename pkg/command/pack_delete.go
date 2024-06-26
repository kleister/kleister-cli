package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packDeleteBind struct {
	ID string
}

var (
	packDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	packDeleteArgs = packDeleteBind{}
)

func init() {
	packCmd.AddCommand(packDeleteCmd)

	packDeleteCmd.Flags().StringVarP(
		&packDeleteArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)
}

func packDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packDeleteArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.DeletePackWithResponse(
		ccmd.Context(),
		packDeleteArgs.ID,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully deleted")
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
