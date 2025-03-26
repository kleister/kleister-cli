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

type packBuildListBind struct {
	Pack   string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplPackBuildList represents a row within pack build listing.
var tmplPackBuildList = "Name: \x1b[33m{{ .Name }} \x1b[0m" + `
ID: {{ .Id }}
`

var (
	packBuildListCmd = &cobra.Command{
		Use:   "list",
		Short: "List pack builds",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packBuildListAction)
		},
		Args: cobra.NoArgs,
	}

	packBuildListArgs = packBuildListBind{}
)

func init() {
	packBuildCmd.AddCommand(packBuildListCmd)

	packBuildListCmd.Flags().StringVar(
		&packBuildListArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	packBuildListCmd.Flags().StringVar(
		&packBuildListArgs.Format,
		"format",
		tmplPackBuildList,
		"Custom output format",
	)

	packBuildListCmd.Flags().StringVar(
		&packBuildListArgs.Search,
		"search",
		"",
		"Search query",
	)

	packBuildListCmd.Flags().StringVar(
		&packBuildListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	packBuildListCmd.Flags().StringVar(
		&packBuildListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	packBuildListCmd.Flags().IntVar(
		&packBuildListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	packBuildListCmd.Flags().IntVar(
		&packBuildListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func packBuildListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packBuildListArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or slug")
	}

	params := &kleister.ListBuildsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if packBuildListArgs.Search != "" {
		params.Search = kleister.ToPtr(packBuildListArgs.Search)
	}

	if packBuildListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(packBuildListArgs.Sort)
	}

	if packBuildListArgs.Order != "" {
		val, err := kleister.ToListBuildsParamsOrder(packBuildListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListBuildsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if packBuildListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(packBuildListArgs.Limit)
	}

	if packBuildListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(packBuildListArgs.Offset)
	}

	resp, err := client.ListBuildsWithResponse(
		ccmd.Context(),
		packBuildListArgs.Pack,
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
		fmt.Sprintln(packBuildListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Builds

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
