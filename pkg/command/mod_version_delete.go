package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modVersionDeleteBind struct {
	Mod string
	ID  string
}

var (
	modVersionDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a mod version",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modVersionDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	modVersionDeleteArgs = modVersionDeleteBind{}
)

func init() {
	modVersionCmd.AddCommand(modVersionDeleteCmd)

	modVersionDeleteCmd.Flags().StringVar(
		&modVersionDeleteArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	modVersionDeleteCmd.Flags().StringVarP(
		&modVersionDeleteArgs.ID,
		"id",
		"i",
		"",
		"Version ID or slug",
	)
}

func modVersionDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modVersionUpdateArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or slug")
	}

	if modVersionUpdateArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or slug")
	}

	resp, err := client.DeleteVersionWithResponse(
		ccmd.Context(),
		modVersionDeleteArgs.Mod,
		modVersionDeleteArgs.ID,
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
