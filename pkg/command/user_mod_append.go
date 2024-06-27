package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userModAppendBind struct {
	ID   string
	Mod  string
	Perm string
}

var (
	userModAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append mod to user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userModAppendAction)
		},
		Args: cobra.NoArgs,
	}

	userModAppendArgs = userModAppendBind{}
)

func init() {
	userModCmd.AddCommand(userModAppendCmd)

	userModAppendCmd.Flags().StringVarP(
		&userModAppendArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userModAppendCmd.Flags().StringVar(
		&userModAppendArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	userModAppendCmd.Flags().StringVar(
		&userModAppendArgs.Perm,
		"perm",
		"",
		"Role for the mod",
	)
}

func userModAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userModAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userModAppendArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	body := kleister.AttachUserToModJSONRequestBody{
		Mod: userModAppendArgs.Mod,
	}

	if userModAppendArgs.Perm != "" {
		val, err := kleister.ToUserModParamsPerm(userModAppendArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrUserModParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.AttachUserToModWithResponse(
		ccmd.Context(),
		userModAppendArgs.ID,
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
