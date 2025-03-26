package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type groupDeleteBind struct {
	ID string
}

var (
	groupDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	groupDeleteArgs = groupDeleteBind{}
)

func init() {
	groupCmd.AddCommand(groupDeleteCmd)

	groupDeleteCmd.Flags().StringVarP(
		&groupDeleteArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)
}

func groupDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupDeleteArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.DeleteGroupWithResponse(
		ccmd.Context(),
		groupDeleteArgs.ID,
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
