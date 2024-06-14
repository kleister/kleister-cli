package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamDeleteBind struct {
	ID string
}

var (
	teamDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamDeleteAction)
		},
		Args: cobra.NoArgs,
	}

	teamDeleteArgs = teamDeleteBind{}
)

func init() {
	teamCmd.AddCommand(teamDeleteCmd)

	teamDeleteCmd.Flags().StringVarP(
		&teamDeleteArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)
}

func teamDeleteAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamDeleteArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.DeleteTeamWithResponse(
		ccmd.Context(),
		teamDeleteArgs.ID,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully delete")
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
