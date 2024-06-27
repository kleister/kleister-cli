package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamPackPermitBind struct {
	ID   string
	Pack string
	Perm string
}

var (
	teamPackPermitCmd = &cobra.Command{
		Use:   "permit",
		Short: "Permit pack for team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamPackPermitAction)
		},
		Args: cobra.NoArgs,
	}

	teamPackPermitArgs = teamPackPermitBind{}
)

func init() {
	teamPackCmd.AddCommand(teamPackPermitCmd)

	teamPackPermitCmd.Flags().StringVarP(
		&teamPackPermitArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamPackPermitCmd.Flags().StringVar(
		&teamPackPermitArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	teamPackPermitCmd.Flags().StringVar(
		&teamPackPermitArgs.Perm,
		"perm",
		"",
		"Role for the pack",
	)
}

func teamPackPermitAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamPackPermitArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamPackPermitArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	body := kleister.PermitTeamPackJSONRequestBody{
		Pack: teamPackPermitArgs.Pack,
	}

	if teamPackPermitArgs.Perm != "" {
		val, err := kleister.ToTeamPackParamsPerm(teamPackPermitArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrTeamPackParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.PermitTeamPackWithResponse(
		ccmd.Context(),
		teamPackPermitArgs.ID,
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
