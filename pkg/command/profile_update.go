package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type profileUpdateBind struct {
	Username string
	Password string
	Email    string
	Fullname string
}

var (
	profileUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update profile details",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, profileUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	profileUpdateArgs = profileUpdateBind{}
)

func init() {
	profileCmd.AddCommand(profileUpdateCmd)

	profileUpdateCmd.Flags().StringVar(
		&profileUpdateArgs.Username,
		"username",
		"",
		"Username for your profile",
	)

	profileUpdateCmd.Flags().StringVar(
		&profileUpdateArgs.Password,
		"password",
		"",
		"Password for your profile",
	)

	profileUpdateCmd.Flags().StringVar(
		&profileUpdateArgs.Email,
		"email",
		"",
		"Email for your profile",
	)

	profileUpdateCmd.Flags().StringVar(
		&profileUpdateArgs.Fullname,
		"fullname",
		"",
		"Fullname for your profile",
	)
}

func profileUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	body := kleister.UpdateProfileJSONRequestBody{}
	changed := false

	if val := profileUpdateArgs.Username; val != "" {
		body.Username = kleister.ToPtr(val)
		changed = true
	}

	if val := profileUpdateArgs.Password; val != "" {
		body.Password = kleister.ToPtr(val)
		changed = true
	}

	if val := profileUpdateArgs.Email; val != "" {
		body.Email = kleister.ToPtr(val)
		changed = true
	}

	if val := profileUpdateArgs.Fullname; val != "" {
		body.Fullname = kleister.ToPtr(val)
		changed = true
	}

	if !changed {
		fmt.Fprintln(os.Stderr, "nothing to update...")
		return nil
	}

	resp, err := client.UpdateProfileWithResponse(
		ccmd.Context(),
		body,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully update")
	case http.StatusUnprocessableEntity:
		return validationError(resp.JSON422)
	case http.StatusForbidden:
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
