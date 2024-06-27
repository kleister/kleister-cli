package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modTeamPermitBind struct {
	ID   string
	Team string
	Perm string
}

var (
	modTeamPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit team for mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modTeamPermitAction)
		},
		Args: cobra.NoArgs,
	}

	modTeamPermitArgs = modTeamPermitBind{}
)

func init() {
	modTeamCmd.AddCommand(modTeamPermitCmd)

	modTeamPermitCmd.Flags().StringVarP(
		&modTeamPermitArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modTeamPermitCmd.Flags().StringVar(
		&modTeamPermitArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)

	modTeamPermitCmd.Flags().StringVar(
		&modTeamPermitArgs.Perm,
		"perm",
		"",
		"Role for the team",
	)
}

func modTeamPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modTeamPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modTeamPermitArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	body := kleister.PermitModTeamJSONRequestBody{
		Team: modTeamPermitArgs.Team,
	}

	if modTeamPermitArgs.Perm != "" {
		val, err := kleister.ToModTeamParamsPerm(modTeamPermitArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrModTeamParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.PermitModTeamWithResponse(
		ccmd.Context(),
		modTeamPermitArgs.ID,
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
