package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packBuildDeleteBind struct {
	Pack string
	ID   string
}

var (
	packBuildDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a pack build",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packBuildDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	packBuildDeleteArgs = packBuildDeleteBind{}
)

func init() {
	packBuildCmd.AddCommand(packBuildDeleteCmd)

	packBuildDeleteCmd.Flags().StringVar(
		&packBuildDeleteArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	packBuildDeleteCmd.Flags().StringVarP(
		&packBuildDeleteArgs.ID,
		"id",
		"i",
		"",
		"Build ID or slug",
	)
}

func packBuildDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packBuildUpdateArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or slug")
	}

	if packBuildUpdateArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or slug")
	}

	resp, err := client.DeleteBuildWithResponse(
		ccmd.Context(),
		packBuildDeleteArgs.Pack,
		packBuildDeleteArgs.ID,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully deleted")
	case http.StatusForbidden:
		return errors.New(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusBadRequest:
		return errors.New(kleister.FromPtr(resp.JSON400.Message))
	case http.StatusNotFound:
		return errors.New(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
