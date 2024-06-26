package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamPackRemoveBind struct {
	ID   string
	Pack string
}

var (
	teamPackRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove pack from team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamPackRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	teamPackRemoveArgs = teamPackRemoveBind{}
)

func init() {
	teamPackCmd.AddCommand(teamPackRemoveCmd)

	teamPackRemoveCmd.Flags().StringVarP(
		&teamPackRemoveArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamPackRemoveCmd.Flags().StringVar(
		&teamPackRemoveArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)
}

func teamPackRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamPackRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamPackRemoveArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	resp, err := client.DeleteTeamFromPackWithResponse(
		ccmd.Context(),
		teamPackRemoveArgs.ID,
		kleister.DeleteTeamFromPackJSONRequestBody{
			Pack: teamPackRemoveArgs.Pack,
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
