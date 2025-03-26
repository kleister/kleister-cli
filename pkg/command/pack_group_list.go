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

type packGroupListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplPackGroupList represents a row within pack group listing.
var tmplPackGroupList = "Slug: \x1b[33m{{ .Group.Slug }} \x1b[0m" + `
ID: {{ .Group.Id }}
Name: {{ .Group.Name }}
Perm: {{ .Perm }}
`

var (
	packGroupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned groups for a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packGroupListAction)
		},
		Args: cobra.NoArgs,
	}

	packGroupListArgs = packGroupListBind{}
)

func init() {
	packGroupCmd.AddCommand(packGroupListCmd)

	packGroupListCmd.Flags().StringVarP(
		&packGroupListArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packGroupListCmd.Flags().StringVar(
		&packGroupListArgs.Format,
		"format",
		tmplPackGroupList,
		"Custom output format",
	)

	packGroupListCmd.Flags().StringVar(
		&packGroupListArgs.Search,
		"search",
		"",
		"Search query",
	)

	packGroupListCmd.Flags().StringVar(
		&packGroupListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	packGroupListCmd.Flags().StringVar(
		&packGroupListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	packGroupListCmd.Flags().IntVar(
		&packGroupListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	packGroupListCmd.Flags().IntVar(
		&packGroupListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func packGroupListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packGroupListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListPackGroupsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if packGroupListArgs.Search != "" {
		params.Search = kleister.ToPtr(packGroupListArgs.Search)
	}

	if packGroupListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(packGroupListArgs.Sort)
	}

	if packGroupListArgs.Order != "" {
		val, err := kleister.ToListPackGroupsParamsOrder(packGroupListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListPackGroupsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if packGroupListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(packGroupListArgs.Limit)
	}

	if packGroupListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(packGroupListArgs.Offset)
	}

	resp, err := client.ListPackGroupsWithResponse(
		ccmd.Context(),
		packGroupListArgs.ID,
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
		fmt.Sprintln(packGroupListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Groups

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
