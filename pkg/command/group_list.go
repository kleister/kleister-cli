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

// tmplGroupList represents a row within user listing.
var tmplGroupList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
`

type groupListBind struct {
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

var (
	groupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupListAction)
		},
		Args: cobra.NoArgs,
	}

	groupListArgs = groupListBind{}
)

func init() {
	groupCmd.AddCommand(groupListCmd)

	groupListCmd.Flags().StringVar(
		&groupListArgs.Format,
		"format",
		tmplGroupList,
		"Custom output format",
	)

	groupListCmd.Flags().StringVar(
		&groupListArgs.Search,
		"search",
		"",
		"Search query",
	)

	groupListCmd.Flags().StringVar(
		&groupListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	groupListCmd.Flags().StringVar(
		&groupListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	groupListCmd.Flags().IntVar(
		&groupListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	groupListCmd.Flags().IntVar(
		&groupListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func groupListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	params := &kleister.ListGroupsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if groupListArgs.Search != "" {
		params.Search = kleister.ToPtr(groupListArgs.Search)
	}

	if groupListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(groupListArgs.Sort)
	}

	if groupListArgs.Order != "" {
		val, err := kleister.ToListGroupsParamsOrder(groupListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListGroupsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if groupListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(groupListArgs.Limit)
	}

	if groupListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(groupListArgs.Offset)
	}

	resp, err := client.ListGroupsWithResponse(
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
		fmt.Sprintln(groupListArgs.Format),
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
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
