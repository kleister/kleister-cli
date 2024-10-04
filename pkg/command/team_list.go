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

// tmplTeamList represents a row within user listing.
var tmplTeamList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
`

type teamListBind struct {
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

var (
	teamListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all teams",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamListAction)
		},
		Args: cobra.NoArgs,
	}

	teamListArgs = teamListBind{}
)

func init() {
	teamCmd.AddCommand(teamListCmd)

	teamListCmd.Flags().StringVar(
		&teamListArgs.Format,
		"format",
		tmplTeamList,
		"Custom output format",
	)

	teamListCmd.Flags().StringVar(
		&teamListArgs.Search,
		"search",
		"",
		"Search query",
	)

	teamListCmd.Flags().StringVar(
		&teamListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	teamListCmd.Flags().StringVar(
		&teamListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	teamListCmd.Flags().IntVar(
		&teamListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	teamListCmd.Flags().IntVar(
		&teamListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func teamListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	params := &kleister.ListTeamsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if teamListArgs.Search != "" {
		params.Search = kleister.ToPtr(teamListArgs.Search)
	}

	if teamListArgs.Sort != "" {
		val, err := kleister.ToListTeamsParamsSort(teamListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListTeamsParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if teamListArgs.Order != "" {
		val, err := kleister.ToListTeamsParamsOrder(teamListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListTeamsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if teamListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(teamListArgs.Limit)
	}

	if teamListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(teamListArgs.Offset)
	}

	resp, err := client.ListTeamsWithResponse(
		ccmd.Context(),
		params,
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
		fmt.Sprintln(teamListArgs.Format),
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
		return errors.New(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
