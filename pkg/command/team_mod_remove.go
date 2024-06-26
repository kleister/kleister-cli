package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamModRemoveBind struct {
	ID  string
	Mod string
}

var (
	teamModRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove mod from team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamModRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	teamModRemoveArgs = teamModRemoveBind{}
)

func init() {
	teamModCmd.AddCommand(teamModRemoveCmd)

	teamModRemoveCmd.Flags().StringVarP(
		&teamModRemoveArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamModRemoveCmd.Flags().StringVar(
		&teamModRemoveArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)
}

func teamModRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamModRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamModRemoveArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	resp, err := client.DeleteTeamFromModWithResponse(
		ccmd.Context(),
		teamModRemoveArgs.ID,
		kleister.DeleteTeamFromModJSONRequestBody{
			Mod: teamModRemoveArgs.Mod,
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
