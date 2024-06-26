package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userTeamPermitBind struct {
	ID   string
	Team string
	Perm string
}

var (
	userTeamPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit team for user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userTeamPermitAction)
		},
		Args: cobra.NoArgs,
	}

	userTeamPermitArgs = userTeamPermitBind{}
)

func init() {
	userTeamCmd.AddCommand(userTeamPermitCmd)

	userTeamPermitCmd.Flags().StringVarP(
		&userTeamPermitArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userTeamPermitCmd.Flags().StringVar(
		&userTeamPermitArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)

	userTeamPermitCmd.Flags().StringVar(
		&userTeamPermitArgs.Perm,
		"perm",
		"",
		"Role for the team",
	)
}

func userTeamPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userTeamPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userTeamPermitArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	body := kleister.PermitUserTeamJSONRequestBody{
		Team: userTeamPermitArgs.Team,
	}

	if userTeamPermitArgs.Perm != "" {
		body.Perm = kleister.ToPtr(userTeamPerm(userTeamPermitArgs.Perm))
	}

	resp, err := client.PermitUserTeamWithResponse(
		ccmd.Context(),
		userTeamPermitArgs.ID,
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
