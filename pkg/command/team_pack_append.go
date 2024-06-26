package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamPackAppendBind struct {
	ID   string
	Pack string
	Perm string
}

var (
	teamPackAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append pack to team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamPackAppendAction)
		},
		Args: cobra.NoArgs,
	}

	teamPackAppendArgs = teamPackAppendBind{}
)

func init() {
	teamPackCmd.AddCommand(teamPackAppendCmd)

	teamPackAppendCmd.Flags().StringVarP(
		&teamPackAppendArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamPackAppendCmd.Flags().StringVar(
		&teamPackAppendArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	teamPackAppendCmd.Flags().StringVar(
		&teamPackAppendArgs.Perm,
		"perm",
		"",
		"Role for the pack",
	)
}

func teamPackAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamPackAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if teamPackAppendArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	body := kleister.AttachTeamToPackJSONRequestBody{
		Pack: teamPackAppendArgs.Pack,
	}

	if teamPackAppendArgs.Perm != "" {
		body.Perm = kleister.ToPtr(teamPackPerm(teamPackAppendArgs.Perm))
	}

	resp, err := client.AttachTeamToPackWithResponse(
		ccmd.Context(),
		teamPackAppendArgs.ID,
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
