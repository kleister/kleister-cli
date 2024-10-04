package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamUserRemoveBind struct {
	ID   string
	User string
}

var (
	teamUserRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamUserRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	teamUserRemoveArgs = teamUserRemoveBind{}
)

func init() {
	teamUserCmd.AddCommand(teamUserRemoveCmd)

	teamUserRemoveCmd.Flags().StringVarP(
		&teamUserRemoveArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamUserRemoveCmd.Flags().StringVar(
		&teamUserRemoveArgs.User,
		"user",
		"",
		"User ID or slug",
	)
}

func teamUserRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamUserRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamUserRemoveArgs.User == "" {
		return fmt.Errorf("you must provide a user ID or a slug")
	}

	resp, err := client.DeleteTeamFromUserWithResponse(
		ccmd.Context(),
		teamUserRemoveArgs.ID,
		kleister.DeleteTeamFromUserJSONRequestBody{
			User: teamUserRemoveArgs.User,
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
