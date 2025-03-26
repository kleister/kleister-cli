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

type groupModListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplGroupModList represents a row within group mod listing.
var tmplGroupModList = "Slug: \x1b[33m{{ .Mod.Slug }} \x1b[0m" + `
ID: {{ .Mod.Id }}
Name: {{ .Mod.Name }}
Perm: {{ .Perm }}
`

var (
	groupModListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned mods for a group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupModListAction)
		},
		Args: cobra.NoArgs,
	}

	groupModListArgs = groupModListBind{}
)

func init() {
	groupModCmd.AddCommand(groupModListCmd)

	groupModListCmd.Flags().StringVarP(
		&groupModListArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupModListCmd.Flags().StringVar(
		&groupModListArgs.Format,
		"format",
		tmplGroupModList,
		"Custom output format",
	)

	groupModListCmd.Flags().StringVar(
		&groupModListArgs.Search,
		"search",
		"",
		"Search query",
	)

	groupModListCmd.Flags().StringVar(
		&groupModListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	groupModListCmd.Flags().StringVar(
		&groupModListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	groupModListCmd.Flags().IntVar(
		&groupModListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	groupModListCmd.Flags().IntVar(
		&groupModListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func groupModListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupModListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListGroupModsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if groupModListArgs.Search != "" {
		params.Search = kleister.ToPtr(groupModListArgs.Search)
	}

	if groupModListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(groupModListArgs.Sort)
	}

	if groupModListArgs.Order != "" {
		val, err := kleister.ToListGroupModsParamsOrder(groupModListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListGroupModsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if groupModListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(groupModListArgs.Limit)
	}

	if groupModListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(groupModListArgs.Offset)
	}

	resp, err := client.ListGroupModsWithResponse(
		ccmd.Context(),
		groupModListArgs.ID,
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
		fmt.Sprintln(groupModListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Mods

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
