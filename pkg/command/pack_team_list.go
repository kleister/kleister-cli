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

type packTeamListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplPackTeamList represents a row within pack team listing.
var tmplPackTeamList = "Slug: \x1b[33m{{ .Team.Slug }} \x1b[0m" + `
ID: {{ .Team.Id }}
Name: {{ .Team.Name }}
Perm: {{ .Perm }}
`

var (
	packTeamListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned teams for a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packTeamListAction)
		},
		Args: cobra.NoArgs,
	}

	packTeamListArgs = packTeamListBind{}
)

func init() {
	packTeamCmd.AddCommand(packTeamListCmd)

	packTeamListCmd.Flags().StringVarP(
		&packTeamListArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packTeamListCmd.Flags().StringVar(
		&packTeamListArgs.Format,
		"format",
		tmplPackTeamList,
		"Custom output format",
	)

	packTeamListCmd.Flags().StringVar(
		&packTeamListArgs.Search,
		"search",
		"",
		"Search query",
	)

	packTeamListCmd.Flags().StringVar(
		&packTeamListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	packTeamListCmd.Flags().StringVar(
		&packTeamListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	packTeamListCmd.Flags().IntVar(
		&packTeamListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	packTeamListCmd.Flags().IntVar(
		&packTeamListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func packTeamListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packTeamListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListPackTeamsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if packTeamListArgs.Search != "" {
		params.Search = kleister.ToPtr(packTeamListArgs.Search)
	}

	if packTeamListArgs.Sort != "" {
		val, err := kleister.ToListPackTeamsParamsSort(packTeamListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListPackTeamsParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if packTeamListArgs.Order != "" {
		val, err := kleister.ToListPackTeamsParamsOrder(packTeamListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListPackTeamsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if packTeamListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(packTeamListArgs.Limit)
	}

	if packTeamListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(packTeamListArgs.Offset)
	}

	resp, err := client.ListPackTeamsWithResponse(
		ccmd.Context(),
		packTeamListArgs.ID,
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
		fmt.Sprintln(packTeamListArgs.Format),
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
	case http.StatusNotFound:
		return errors.New(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
