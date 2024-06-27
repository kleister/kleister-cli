package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modTeamAppendBind struct {
	ID   string
	Team string
	Perm string
}

var (
	modTeamAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append team to mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modTeamAppendAction)
		},
		Args: cobra.NoArgs,
	}

	modTeamAppendArgs = modTeamAppendBind{}
)

func init() {
	modTeamCmd.AddCommand(modTeamAppendCmd)

	modTeamAppendCmd.Flags().StringVarP(
		&modTeamAppendArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modTeamAppendCmd.Flags().StringVar(
		&modTeamAppendArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)

	modTeamAppendCmd.Flags().StringVar(
		&modTeamAppendArgs.Perm,
		"perm",
		"",
		"Role for the team",
	)
}

func modTeamAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modTeamAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modTeamAppendArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	body := kleister.AttachModToTeamJSONRequestBody{
		Team: modTeamAppendArgs.Team,
	}

	if modTeamAppendArgs.Perm != "" {
		body.Perm = kleister.ToPtr(modTeamPerm(modTeamAppendArgs.Perm))
	}

	resp, err := client.AttachModToTeamWithResponse(
		ccmd.Context(),
		modTeamAppendArgs.ID,
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
