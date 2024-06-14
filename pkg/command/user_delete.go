package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userDeleteBind struct {
	ID string
}

var (
	userDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	userDeleteArgs = userDeleteBind{}
)

func init() {
	userCmd.AddCommand(userDeleteCmd)

	userDeleteCmd.Flags().StringVarP(
		&userDeleteArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)
}

func userDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userDeleteArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.DeleteUserWithResponse(
		ccmd.Context(),
		userDeleteArgs.ID,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully delete")
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
