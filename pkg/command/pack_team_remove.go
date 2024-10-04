package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packTeamRemoveBind struct {
	ID   string
	Team string
}

var (
	packTeamRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove team from pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packTeamRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	packTeamRemoveArgs = packTeamRemoveBind{}
)

func init() {
	packTeamCmd.AddCommand(packTeamRemoveCmd)

	packTeamRemoveCmd.Flags().StringVarP(
		&packTeamRemoveArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packTeamRemoveCmd.Flags().StringVar(
		&packTeamRemoveArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)
}

func packTeamRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packTeamRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packTeamRemoveArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	resp, err := client.DeletePackFromTeamWithResponse(
		ccmd.Context(),
		packTeamRemoveArgs.ID,
		kleister.DeletePackFromTeamJSONRequestBody{
			Team: packTeamRemoveArgs.Team,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
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
