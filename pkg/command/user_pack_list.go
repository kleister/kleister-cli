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

type userPackListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplUserPackList represents a row within user pack listing.
var tmplUserPackList = "Slug: \x1b[33m{{ .Pack.Slug }} \x1b[0m" + `
ID: {{ .Pack.Id }}
Name: {{ .Pack.Name }}
Perm: {{ .Perm }}
`

var (
	userPackListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned packs for a user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userPackListAction)
		},
		Args: cobra.NoArgs,
	}

	userPackListArgs = userPackListBind{}
)

func init() {
	userPackCmd.AddCommand(userPackListCmd)

	userPackListCmd.Flags().StringVarP(
		&userPackListArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userPackListCmd.Flags().StringVar(
		&userPackListArgs.Format,
		"format",
		tmplUserPackList,
		"Custom output format",
	)

	userPackListCmd.Flags().StringVar(
		&userPackListArgs.Search,
		"search",
		"",
		"Search query",
	)

	userPackListCmd.Flags().StringVar(
		&userPackListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	userPackListCmd.Flags().StringVar(
		&userPackListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	userPackListCmd.Flags().IntVar(
		&userPackListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	userPackListCmd.Flags().IntVar(
		&userPackListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func userPackListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userPackListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListUserPacksParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if userPackListArgs.Search != "" {
		params.Search = kleister.ToPtr(userPackListArgs.Search)
	}

	if userPackListArgs.Sort != "" {
		val, err := kleister.ToListUserPacksParamsSort(userPackListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListUserPacksParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if userPackListArgs.Order != "" {
		val, err := kleister.ToListUserPacksParamsOrder(userPackListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListUserPacksParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if userPackListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(userPackListArgs.Limit)
	}

	if userPackListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(userPackListArgs.Offset)
	}

	resp, err := client.ListUserPacksWithResponse(
		ccmd.Context(),
		userPackListArgs.ID,
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
		fmt.Sprintln(userPackListArgs.Format),
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
	case http.StatusNotFound:
		return errors.New(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
