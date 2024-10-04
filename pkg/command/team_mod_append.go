package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamModAppendBind struct {
	ID   string
	Mod  string
	Perm string
}

var (
	teamModAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append mod to team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamModAppendAction)
		},
		Args: cobra.NoArgs,
	}

	teamModAppendArgs = teamModAppendBind{}
)

func init() {
	teamModCmd.AddCommand(teamModAppendCmd)

	teamModAppendCmd.Flags().StringVarP(
		&teamModAppendArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamModAppendCmd.Flags().StringVar(
		&teamModAppendArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	teamModAppendCmd.Flags().StringVar(
		&teamModAppendArgs.Perm,
		"perm",
		"",
		"Role for the mod",
	)
}

func teamModAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamModAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamModAppendArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or a slug")
	}

	body := kleister.AttachTeamToModJSONRequestBody{
		Mod: teamModAppendArgs.Mod,
	}

	if teamModAppendArgs.Perm != "" {
		val, err := kleister.ToTeamModParamsPerm(teamModAppendArgs.Perm)

		if err != nil && errors.Is(err, kleister.ErrTeamModParamsPerm) {
			return fmt.Errorf("invalid perm attribute")
		}

		body.Perm = kleister.ToPtr(val)
	}

	resp, err := client.AttachTeamToModWithResponse(
		ccmd.Context(),
		teamModAppendArgs.ID,
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
