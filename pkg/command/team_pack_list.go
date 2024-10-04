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

type teamPackListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplTeamPackList represents a row within team pack listing.
var tmplTeamPackList = "Slug: \x1b[33m{{ .Pack.Slug }} \x1b[0m" + `
ID: {{ .Pack.Id }}
Name: {{ .Pack.Name }}
Perm: {{ .Perm }}
`

var (
	teamPackListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned packs for a team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamPackListAction)
		},
		Args: cobra.NoArgs,
	}

	teamPackListArgs = teamPackListBind{}
)

func init() {
	teamPackCmd.AddCommand(teamPackListCmd)

	teamPackListCmd.Flags().StringVarP(
		&teamPackListArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamPackListCmd.Flags().StringVar(
		&teamPackListArgs.Format,
		"format",
		tmplTeamPackList,
		"Custom output format",
	)

	teamPackListCmd.Flags().StringVar(
		&teamPackListArgs.Search,
		"search",
		"",
		"Search query",
	)

	teamPackListCmd.Flags().StringVar(
		&teamPackListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	teamPackListCmd.Flags().StringVar(
		&teamPackListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	teamPackListCmd.Flags().IntVar(
		&teamPackListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	teamPackListCmd.Flags().IntVar(
		&teamPackListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func teamPackListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamPackListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListTeamPacksParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if teamPackListArgs.Search != "" {
		params.Search = kleister.ToPtr(teamPackListArgs.Search)
	}

	if teamPackListArgs.Sort != "" {
		val, err := kleister.ToListTeamPacksParamsSort(teamPackListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListTeamPacksParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if teamPackListArgs.Order != "" {
		val, err := kleister.ToListTeamPacksParamsOrder(teamPackListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListTeamPacksParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if teamPackListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(teamPackListArgs.Limit)
	}

	if teamPackListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(teamPackListArgs.Offset)
	}

	resp, err := client.ListTeamPacksWithResponse(
		ccmd.Context(),
		teamPackListArgs.ID,
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
		fmt.Sprintln(teamPackListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Packs)

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
	case http.StatusNotFound:
		return errors.New(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
