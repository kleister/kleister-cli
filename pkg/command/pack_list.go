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

// tmplPackList represents a row within user listing.
var tmplPackList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
`

type packListBind struct {
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

var (
	packListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all packs",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packListAction)
		},
		Args: cobra.NoArgs,
	}

	packListArgs = packListBind{}
)

func init() {
	packCmd.AddCommand(packListCmd)

	packListCmd.Flags().StringVar(
		&packListArgs.Format,
		"format",
		tmplPackList,
		"Custom output format",
	)

	packListCmd.Flags().StringVar(
		&packListArgs.Search,
		"search",
		"",
		"Search query",
	)

	packListCmd.Flags().StringVar(
		&packListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	packListCmd.Flags().StringVar(
		&packListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	packListCmd.Flags().IntVar(
		&packListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	packListCmd.Flags().IntVar(
		&packListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func packListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	params := &kleister.ListPacksParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if packListArgs.Search != "" {
		params.Search = kleister.ToPtr(packListArgs.Search)
	}

	if packListArgs.Sort != "" {
		val, err := kleister.ToListPacksParamsSort(packListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListPacksParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if packListArgs.Order != "" {
		val, err := kleister.ToListPacksParamsOrder(packListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListPacksParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if packListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(packListArgs.Limit)
	}

	if packListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(packListArgs.Offset)
	}

	resp, err := client.ListPacksWithResponse(
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
		fmt.Sprintln(packListArgs.Format),
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
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
