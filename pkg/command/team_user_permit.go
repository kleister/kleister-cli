package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamUserPermitBind struct {
	ID   string
	User string
	Perm string
}

var (
	teamUserPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit user for team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamUserPermitAction)
		},
		Args: cobra.NoArgs,
	}

	teamUserPermitArgs = teamUserPermitBind{}
)

func init() {
	teamUserCmd.AddCommand(teamUserPermitCmd)

	teamUserPermitCmd.Flags().StringVarP(
		&teamUserPermitArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamUserPermitCmd.Flags().StringVar(
		&teamUserPermitArgs.User,
		"user",
		"",
		"User ID or slug",
	)

	teamUserPermitCmd.Flags().StringVar(
		&teamUserPermitArgs.Perm,
		"perm",
		"",
		"Role for the user",
	)
}

func teamUserPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamUserPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamUserPermitArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	body := kleister.PermitTeamUserJSONRequestBody{
		User: teamUserPermitArgs.User,
	}

	if teamUserPermitArgs.Perm != "" {
		val, err := kleister.ToTeamUserParamsPerm(teamUserPermitArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrTeamUserParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.PermitTeamUserWithResponse(
		ccmd.Context(),
		teamUserPermitArgs.ID,
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
