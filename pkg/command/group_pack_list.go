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

type groupPackListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplGroupPackList represents a row within group pack listing.
var tmplGroupPackList = "Slug: \x1b[33m{{ .Pack.Slug }} \x1b[0m" + `
ID: {{ .Pack.Id }}
Name: {{ .Pack.Name }}
Perm: {{ .Perm }}
`

var (
	groupPackListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned packs for a group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupPackListAction)
		},
		Args: cobra.NoArgs,
	}

	groupPackListArgs = groupPackListBind{}
)

func init() {
	groupPackCmd.AddCommand(groupPackListCmd)

	groupPackListCmd.Flags().StringVarP(
		&groupPackListArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupPackListCmd.Flags().StringVar(
		&groupPackListArgs.Format,
		"format",
		tmplGroupPackList,
		"Custom output format",
	)

	groupPackListCmd.Flags().StringVar(
		&groupPackListArgs.Search,
		"search",
		"",
		"Search query",
	)

	groupPackListCmd.Flags().StringVar(
		&groupPackListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	groupPackListCmd.Flags().StringVar(
		&groupPackListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	groupPackListCmd.Flags().IntVar(
		&groupPackListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	groupPackListCmd.Flags().IntVar(
		&groupPackListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func groupPackListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupPackListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListGroupPacksParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if groupPackListArgs.Search != "" {
		params.Search = kleister.ToPtr(groupPackListArgs.Search)
	}

	if groupPackListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(groupPackListArgs.Sort)
	}

	if groupPackListArgs.Order != "" {
		val, err := kleister.ToListGroupPacksParamsOrder(groupPackListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListGroupPacksParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if groupPackListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(groupPackListArgs.Limit)
	}

	if groupPackListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(groupPackListArgs.Offset)
	}

	resp, err := client.ListGroupPacksWithResponse(
		ccmd.Context(),
		groupPackListArgs.ID,
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
		fmt.Sprintln(groupPackListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Packs

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
