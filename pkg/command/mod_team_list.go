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

type modTeamListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplModTeamList represents a row within mod team listing.
var tmplModTeamList = "Slug: \x1b[33m{{ .Team.Slug }} \x1b[0m" + `
ID: {{ .Team.Id }}
Name: {{ .Team.Name }}
Perm: {{ .Perm }}
`

var (
	modTeamListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned teams for a mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modTeamListAction)
		},
		Args: cobra.NoArgs,
	}

	modTeamListArgs = modTeamListBind{}
)

func init() {
	modTeamCmd.AddCommand(modTeamListCmd)

	modTeamListCmd.Flags().StringVarP(
		&modTeamListArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modTeamListCmd.Flags().StringVar(
		&modTeamListArgs.Format,
		"format",
		tmplModTeamList,
		"Custom output format",
	)

	modTeamListCmd.Flags().StringVar(
		&modTeamListArgs.Search,
		"search",
		"",
		"Search query",
	)

	modTeamListCmd.Flags().StringVar(
		&modTeamListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	modTeamListCmd.Flags().StringVar(
		&modTeamListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	modTeamListCmd.Flags().IntVar(
		&modTeamListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	modTeamListCmd.Flags().IntVar(
		&modTeamListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func modTeamListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modTeamListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListModTeamsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if modTeamListArgs.Search != "" {
		params.Search = kleister.ToPtr(modTeamListArgs.Search)
	}

	if modTeamListArgs.Sort != "" {
		val, err := kleister.ToListModTeamsParamsSort(modTeamListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListModTeamsParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if modTeamListArgs.Order != "" {
		val, err := kleister.ToListModTeamsParamsOrder(modTeamListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListModTeamsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if modTeamListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(modTeamListArgs.Limit)
	}

	if modTeamListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(modTeamListArgs.Offset)
	}

	resp, err := client.ListModTeamsWithResponse(
		ccmd.Context(),
		modTeamListArgs.ID,
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
		fmt.Sprintln(modTeamListArgs.Format),
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
