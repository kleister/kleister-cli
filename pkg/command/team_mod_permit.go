package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamModPermitBind struct {
	ID   string
	Mod  string
	Perm string
}

var (
	teamModPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit mod for team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamModPermitAction)
		},
		Args: cobra.NoArgs,
	}

	teamModPermitArgs = teamModPermitBind{}
)

func init() {
	teamModCmd.AddCommand(teamModPermitCmd)

	teamModPermitCmd.Flags().StringVarP(
		&teamModPermitArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamModPermitCmd.Flags().StringVar(
		&teamModPermitArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	teamModPermitCmd.Flags().StringVar(
		&teamModPermitArgs.Perm,
		"perm",
		"",
		"Role for the mod",
	)
}

func teamModPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamModPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamModPermitArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	body := kleister.PermitTeamModJSONRequestBody{
		Mod: teamModPermitArgs.Mod,
	}

	if teamModPermitArgs.Perm != "" {
		val, err := kleister.ToTeamModParamsPerm(teamModPermitArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrTeamModParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.PermitTeamModWithResponse(
		ccmd.Context(),
		teamModPermitArgs.ID,
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
