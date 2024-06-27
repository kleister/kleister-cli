package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packUserAppendBind struct {
	ID   string
	User string
	Perm string
}

var (
	packUserAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append user to pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packUserAppendAction)
		},
		Args: cobra.NoArgs,
	}

	packUserAppendArgs = packUserAppendBind{}
)

func init() {
	packUserCmd.AddCommand(packUserAppendCmd)

	packUserAppendCmd.Flags().StringVarP(
		&packUserAppendArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packUserAppendCmd.Flags().StringVar(
		&packUserAppendArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	packUserAppendCmd.Flags().StringVar(
		&packUserAppendArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func packUserAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packUserAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packUserAppendArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.AttachPackToUserJSONRequestBody{
		User: packUserAppendArgs.User,
	}

	if packUserAppendArgs.Perm != "" {
		body.Perm = kleister.ToPtr(packUserPerm(packUserAppendArgs.Perm))
	}

	resp, err := client.AttachPackToUserWithResponse(
		ccmd.Context(),
		packUserAppendArgs.ID,
		body,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusUnprocessableEntity:
		return validationError(resp.JSON422)
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
