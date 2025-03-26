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

type modGroupListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplModGroupList represents a row within mod group listing.
var tmplModGroupList = "Slug: \x1b[33m{{ .Group.Slug }} \x1b[0m" + `
ID: {{ .Group.Id }}
Name: {{ .Group.Name }}
Perm: {{ .Perm }}
`

var (
	modGroupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned groups for a mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modGroupListAction)
		},
		Args: cobra.NoArgs,
	}

	modGroupListArgs = modGroupListBind{}
)

func init() {
	modGroupCmd.AddCommand(modGroupListCmd)

	modGroupListCmd.Flags().StringVarP(
		&modGroupListArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modGroupListCmd.Flags().StringVar(
		&modGroupListArgs.Format,
		"format",
		tmplModGroupList,
		"Custom output format",
	)

	modGroupListCmd.Flags().StringVar(
		&modGroupListArgs.Search,
		"search",
		"",
		"Search query",
	)

	modGroupListCmd.Flags().StringVar(
		&modGroupListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	modGroupListCmd.Flags().StringVar(
		&modGroupListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	modGroupListCmd.Flags().IntVar(
		&modGroupListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	modGroupListCmd.Flags().IntVar(
		&modGroupListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func modGroupListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modGroupListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListModGroupsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if modGroupListArgs.Search != "" {
		params.Search = kleister.ToPtr(modGroupListArgs.Search)
	}

	if modGroupListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(modGroupListArgs.Sort)
	}

	if modGroupListArgs.Order != "" {
		val, err := kleister.ToListModGroupsParamsOrder(modGroupListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListModGroupsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if modGroupListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(modGroupListArgs.Limit)
	}

	if modGroupListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(modGroupListArgs.Offset)
	}

	resp, err := client.ListModGroupsWithResponse(
		ccmd.Context(),
		modGroupListArgs.ID,
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
		fmt.Sprintln(modGroupListArgs.Format),
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
