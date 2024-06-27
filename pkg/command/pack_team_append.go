package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type packTeamAppendBind struct {
	ID   string
	Team string
	Perm string
}

var (
	packTeamAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append team to pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packTeamAppendAction)
		},
		Args: cobra.NoArgs,
	}

	packTeamAppendArgs = packTeamAppendBind{}
)

func init() {
	packTeamCmd.AddCommand(packTeamAppendCmd)

	packTeamAppendCmd.Flags().StringVarP(
		&packTeamAppendArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packTeamAppendCmd.Flags().StringVar(
		&packTeamAppendArgs.Team,
		"team",
		"",
		"Team ID or slug",
	)

	packTeamAppendCmd.Flags().StringVar(
		&packTeamAppendArgs.Perm,
		"perm",
		"",
		"Role for the team",
	)
}

func packTeamAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packTeamAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if packTeamAppendArgs.Team == "" {
		return fmt.Errorf("you must provide a team ID or a slug")
	}

	body := kleister.AttachPackToTeamJSONRequestBody{
		Team: packTeamAppendArgs.Team,
	}

	if packTeamAppendArgs.Perm != "" {
		val, err := kleister.ToPackTeamParamsPerm(packTeamAppendArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrPackTeamParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.AttachPackToTeamWithResponse(
		ccmd.Context(),
		packTeamAppendArgs.ID,
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
