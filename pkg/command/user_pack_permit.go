package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userPackPermitBind struct {
	ID   string
	Pack string
	Perm string
}

var (
	userPackPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit pack for user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userPackPermitAction)
		},
		Args: cobra.NoArgs,
	}

	userPackPermitArgs = userPackPermitBind{}
)

func init() {
	userPackCmd.AddCommand(userPackPermitCmd)

	userPackPermitCmd.Flags().StringVarP(
		&userPackPermitArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userPackPermitCmd.Flags().StringVar(
		&userPackPermitArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	userPackPermitCmd.Flags().StringVar(
		&userPackPermitArgs.Perm,
		"perm",
		"",
		"Role for the pack",
	)
}

func userPackPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userPackPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userPackPermitArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	body := kleister.PermitUserPackJSONRequestBody{
		Pack: userPackPermitArgs.Pack,
	}

	if userPackPermitArgs.Perm != "" {
		val, err := kleister.ToUserPackParamsPerm(userPackPermitArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrUserPackParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.PermitUserPackWithResponse(
		ccmd.Context(),
		userPackPermitArgs.ID,
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
