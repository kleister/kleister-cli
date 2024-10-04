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

type teamModListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplTeamModList represents a row within team mod listing.
var tmplTeamModList = "Slug: \x1b[33m{{ .Mod.Slug }} \x1b[0m" + `
ID: {{ .Mod.Id }}
Name: {{ .Mod.Name }}
Perm: {{ .Perm }}
`

var (
	teamModListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned mods for a team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamModListAction)
		},
		Args: cobra.NoArgs,
	}

	teamModListArgs = teamModListBind{}
)

func init() {
	teamModCmd.AddCommand(teamModListCmd)

	teamModListCmd.Flags().StringVarP(
		&teamModListArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamModListCmd.Flags().StringVar(
		&teamModListArgs.Format,
		"format",
		tmplTeamModList,
		"Custom output format",
	)

	teamModListCmd.Flags().StringVar(
		&teamModListArgs.Search,
		"search",
		"",
		"Search query",
	)

	teamModListCmd.Flags().StringVar(
		&teamModListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	teamModListCmd.Flags().StringVar(
		&teamModListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	teamModListCmd.Flags().IntVar(
		&teamModListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	teamModListCmd.Flags().IntVar(
		&teamModListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func teamModListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamModListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListTeamModsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if teamModListArgs.Search != "" {
		params.Search = kleister.ToPtr(teamModListArgs.Search)
	}

	if teamModListArgs.Sort != "" {
		val, err := kleister.ToListTeamModsParamsSort(teamModListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListTeamModsParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if teamModListArgs.Order != "" {
		val, err := kleister.ToListTeamModsParamsOrder(teamModListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListTeamModsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if teamModListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(teamModListArgs.Limit)
	}

	if teamModListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(teamModListArgs.Offset)
	}

	resp, err := client.ListTeamModsWithResponse(
		ccmd.Context(),
		teamModListArgs.ID,
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
		fmt.Sprintln(teamModListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Mods)

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
