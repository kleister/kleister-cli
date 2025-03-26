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

type modVersionListBind struct {
	Mod    string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplModVersionList represents a row within mod version listing.
var tmplModVersionList = "Name: \x1b[33m{{ .Name }} \x1b[0m" + `
ID: {{ .Id }}
`

var (
	modVersionListCmd = &cobra.Command{
		Use:   "list",
		Short: "List mod versions",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modVersionListAction)
		},
		Args: cobra.NoArgs,
	}

	modVersionListArgs = modVersionListBind{}
)

func init() {
	modVersionCmd.AddCommand(modVersionListCmd)

	modVersionListCmd.Flags().StringVar(
		&modVersionListArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	modVersionListCmd.Flags().StringVar(
		&modVersionListArgs.Format,
		"format",
		tmplModVersionList,
		"Custom output format",
	)

	modVersionListCmd.Flags().StringVar(
		&modVersionListArgs.Search,
		"search",
		"",
		"Search query",
	)

	modVersionListCmd.Flags().StringVar(
		&modVersionListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	modVersionListCmd.Flags().StringVar(
		&modVersionListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	modVersionListCmd.Flags().IntVar(
		&modVersionListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	modVersionListCmd.Flags().IntVar(
		&modVersionListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func modVersionListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modVersionListArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or slug")
	}

	params := &kleister.ListVersionsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if modVersionListArgs.Search != "" {
		params.Search = kleister.ToPtr(modVersionListArgs.Search)
	}

	if modVersionListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(modVersionListArgs.Sort)
	}

	if modVersionListArgs.Order != "" {
		val, err := kleister.ToListVersionsParamsOrder(modVersionListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListVersionsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if modVersionListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(modVersionListArgs.Limit)
	}

	if modVersionListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(modVersionListArgs.Offset)
	}

	resp, err := client.ListVersionsWithResponse(
		ccmd.Context(),
		modVersionListArgs.Mod,
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
		fmt.Sprintln(modVersionListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Versions

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
