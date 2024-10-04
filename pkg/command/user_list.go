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

// tmplUserList represents a row within user listing.
var tmplUserList = "Username: \x1b[33m{{ .Username }} \x1b[0m" + `
ID: {{ .Id }}
Email: {{ .Email }}
`

type userListBind struct {
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

var (
	userListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userListAction)
		},
		Args: cobra.NoArgs,
	}

	userListArgs = userListBind{}
)

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().StringVar(
		&userListArgs.Format,
		"format",
		tmplUserList,
		"Custom output format",
	)

	userListCmd.Flags().StringVar(
		&userListArgs.Search,
		"search",
		"",
		"Search query",
	)

	userListCmd.Flags().StringVar(
		&userListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	userListCmd.Flags().StringVar(
		&userListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	userListCmd.Flags().IntVar(
		&userListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	userListCmd.Flags().IntVar(
		&userListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func userListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	params := &kleister.ListUsersParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if userListArgs.Search != "" {
		params.Search = kleister.ToPtr(userListArgs.Search)
	}

	if userListArgs.Sort != "" {
		val, err := kleister.ToListUsersParamsSort(userListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListUsersParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if userListArgs.Order != "" {
		val, err := kleister.ToListUsersParamsOrder(userListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListUsersParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if userListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(userListArgs.Limit)
	}

	if userListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(userListArgs.Offset)
	}

	resp, err := client.ListUsersWithResponse(
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
		fmt.Sprintln(userListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Users)

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
