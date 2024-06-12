package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamUserListBind struct {
	ID     string
	Format string
}

// tmplTeamUserList represents a row within team user listing.
var tmplTeamUserList = "Slug: \x1b[33m{{ .User.Username }} \x1b[0m" + `
ID: {{ .User.Id }}
Email: {{ .User.Email }}
Perm: {{ .Perm }}
`

var (
	teamUserListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned users for a team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamUserListAction)
		},
		Args: cobra.NoArgs,
	}

	teamUserListArgs = teamUserListBind{}
)

func init() {
	teamUserCmd.AddCommand(teamUserListCmd)

	teamUserListCmd.Flags().StringVarP(
		&teamUserListArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamUserListCmd.Flags().StringVar(
		&teamUserListArgs.Format,
		"format",
		tmplTeamUserList,
		"Custom output format",
	)
}

func teamUserListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamUserListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ListTeamUsersWithResponse(
		ccmd.Context(),
		teamUserListArgs.ID,
		&kleister.ListTeamUsersParams{
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
		fmt.Sprintln(teamUserListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Users)

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
