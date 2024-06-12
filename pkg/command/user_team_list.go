package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type userTeamListBind struct {
	ID     string
	Format string
}

// tmplUserTeamList represents a row within user team listing.
var tmplUserTeamList = "Slug: \x1b[33m{{ .Team.Slug }} \x1b[0m" + `
ID: {{ .Team.Id }}
Name: {{ .Team.Name }}
Perm: {{ .Perm }}
`

var (
	userTeamListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned teams for a user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userTeamListAction)
		},
		Args: cobra.NoArgs,
	}

	userTeamListArgs = userTeamListBind{}
)

func init() {
	userTeamCmd.AddCommand(userTeamListCmd)

	userTeamListCmd.Flags().StringVarP(
		&userTeamListArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userTeamListCmd.Flags().StringVar(
		&userTeamListArgs.Format,
		"format",
		tmplUserTeamList,
		"Custom output format",
	)
}

func userTeamListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userTeamListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ListUserTeamsWithResponse(
		ccmd.Context(),
		userTeamListArgs.ID,
		&kleister.ListUserTeamsParams{
			Limit:  kleister.ToPtr(10000),
			Offset: kleister.ToPtr(0),
		},
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
		fmt.Sprintln(userTeamListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Teams)

		if len(records) == 0 {
			fmt.Fprintln(os.Stderr, "Empty result")
			return nil
		}

		for _, record := range records {
			if err := tmpl.Execute(
				os.Stdout,
				record,
			); err != nil {
				return fmt.Errorf("failed to render template: %w", err)
			}
		}
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
