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

// tmplModList represents a row within user listing.
var tmplModList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
`

type modListBind struct {
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

var (
	modListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all mods",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modListAction)
		},
		Args: cobra.NoArgs,
	}

	modListArgs = modListBind{}
)

func init() {
	modCmd.AddCommand(modListCmd)

	modListCmd.Flags().StringVar(
		&modListArgs.Format,
		"format",
		tmplModList,
		"Custom output format",
	)

	modListCmd.Flags().StringVar(
		&modListArgs.Search,
		"search",
		"",
		"Search query",
	)

	modListCmd.Flags().StringVar(
		&modListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	modListCmd.Flags().StringVar(
		&modListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	modListCmd.Flags().IntVar(
		&modListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	modListCmd.Flags().IntVar(
		&modListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func modListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	params := &kleister.ListModsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if modListArgs.Search != "" {
		params.Search = kleister.ToPtr(modListArgs.Search)
	}

	if modListArgs.Sort != "" {
		val, err := kleister.ToListModsParamsSort(modListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListModsParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if modListArgs.Order != "" {
		val, err := kleister.ToListModsParamsOrder(modListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListModsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if modListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(modListArgs.Limit)
	}

	if modListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(modListArgs.Offset)
	}

	resp, err := client.ListModsWithResponse(
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
		fmt.Sprintln(modListArgs.Format),
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
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
