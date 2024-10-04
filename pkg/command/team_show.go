package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

// tmplTeamShow represents a user within details view.
var tmplTeamShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type teamShowBind struct {
	ID     string
	Format string
}

var (
	teamShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show an team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamShowAction)
		},
		Args: cobra.NoArgs,
	}

	teamShowArgs = teamShowBind{}
)

func init() {
	teamCmd.AddCommand(teamShowCmd)

	teamShowCmd.Flags().StringVarP(
		&teamShowArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamShowCmd.Flags().StringVar(
		&teamShowArgs.Format,
		"format",
		tmplTeamShow,
		"Custom output format",
	)
}

func teamShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ShowTeamWithResponse(
		ccmd.Context(),
		teamShowArgs.ID,
	)

	if err != nil {
		return err
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		basicFuncMap,
	).Parse(
		fmt.Sprintln(teamShowArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if err := tmpl.Execute(
			os.Stdout,
			resp.JSON200,
		); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
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
