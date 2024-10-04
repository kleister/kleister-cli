package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userTeamRemoveBind struct {
	ID   string
	Team string
}

var (
	userTeamRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove team from user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userTeamRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	userTeamRemoveArgs = userTeamRemoveBind{}
)

func init() {
	userTeamCmd.AddCommand(userTeamRemoveCmd)

	userTeamRemoveCmd.Flags().StringVarP(
		&userTeamRemoveArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userTeamRemoveCmd.Flags().StringVar(
		&userTeamRemoveArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)
}

func userTeamRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userTeamRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if userTeamRemoveArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	resp, err := client.DeleteUserFromTeamWithResponse(
		ccmd.Context(),
		userTeamRemoveArgs.ID,
		kleister.DeleteUserFromTeamJSONRequestBody{
			Team: userTeamRemoveArgs.Team,
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
