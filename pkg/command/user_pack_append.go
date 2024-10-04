package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userPackAppendBind struct {
	ID   string
	Pack string
	Perm string
}

var (
	userPackAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append pack to user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userPackAppendAction)
		},
		Args: cobra.NoArgs,
	}

	userPackAppendArgs = userPackAppendBind{}
)

func init() {
	userPackCmd.AddCommand(userPackAppendCmd)

	userPackAppendCmd.Flags().StringVarP(
		&userPackAppendArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userPackAppendCmd.Flags().StringVar(
		&userPackAppendArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	userPackAppendCmd.Flags().StringVar(
		&userPackAppendArgs.Perm,
		"perm",
		"",
		"Role for the pack",
	)
}

func userPackAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userPackAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userPackAppendArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	body := kleister.AttachUserToPackJSONRequestBody{
		Pack: userPackAppendArgs.Pack,
	}

	if userPackAppendArgs.Perm != "" {
		val, err := kleister.ToUserPackParamsPerm(userPackAppendArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrUserPackParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.AttachUserToPackWithResponse(
		ccmd.Context(),
		userPackAppendArgs.ID,
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
