package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packTeamPermitBind struct {
	ID   string
	Team string
	Perm string
}

var (
	packTeamPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit team for pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packTeamPermitAction)
		},
		Args: cobra.NoArgs,
	}

	packTeamPermitArgs = packTeamPermitBind{}
)

func init() {
	packTeamCmd.AddCommand(packTeamPermitCmd)

	packTeamPermitCmd.Flags().StringVarP(
		&packTeamPermitArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packTeamPermitCmd.Flags().StringVar(
		&packTeamPermitArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)

	packTeamPermitCmd.Flags().StringVar(
		&packTeamPermitArgs.Perm,
		"perm",
		"",
		"Role for the team",
	)
}

func packTeamPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packTeamPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packTeamPermitArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	body := kleister.PermitPackTeamJSONRequestBody{
		Team: packTeamPermitArgs.Team,
	}

	if packTeamPermitArgs.Perm != "" {
		body.Perm = kleister.ToPtr(packTeamPerm(packTeamPermitArgs.Perm))
	}

	resp, err := client.PermitPackTeamWithResponse(
		ccmd.Context(),
		packTeamPermitArgs.ID,
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
