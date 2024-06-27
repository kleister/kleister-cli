package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modTeamRemoveBind struct {
	ID   string
	Team string
}

var (
	modTeamRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove team from mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modTeamRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	modTeamRemoveArgs = modTeamRemoveBind{}
)

func init() {
	modTeamCmd.AddCommand(modTeamRemoveCmd)

	modTeamRemoveCmd.Flags().StringVarP(
		&modTeamRemoveArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modTeamRemoveCmd.Flags().StringVar(
		&modTeamRemoveArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)
}

func modTeamRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modTeamRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if modTeamRemoveArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	resp, err := client.DeleteModFromTeamWithResponse(
		ccmd.Context(),
		modTeamRemoveArgs.ID,
		kleister.DeleteModFromTeamJSONRequestBody{
			Team: modTeamRemoveArgs.Team,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
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
