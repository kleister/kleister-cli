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

type groupUserListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplGroupUserList represents a row within group user listing.
var tmplGroupUserList = "Slug: \x1b[33m{{ .User.Username }} \x1b[0m" + `
ID: {{ .User.Id }}
Email: {{ .User.Email }}
Perm: {{ .Perm }}
`

var (
	groupUserListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned users for a group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupUserListAction)
		},
		Args: cobra.NoArgs,
	}

	groupUserListArgs = groupUserListBind{}
)

func init() {
	groupUserCmd.AddCommand(groupUserListCmd)

	groupUserListCmd.Flags().StringVarP(
		&groupUserListArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupUserListCmd.Flags().StringVar(
		&groupUserListArgs.Format,
		"format",
		tmplGroupUserList,
		"Custom output format",
	)

	groupUserListCmd.Flags().StringVar(
		&groupUserListArgs.Search,
		"search",
		"",
		"Search query",
	)

	groupUserListCmd.Flags().StringVar(
		&groupUserListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	groupUserListCmd.Flags().StringVar(
		&groupUserListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	groupUserListCmd.Flags().IntVar(
		&groupUserListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	groupUserListCmd.Flags().IntVar(
		&groupUserListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func groupUserListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupUserListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListGroupUsersParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if groupUserListArgs.Search != "" {
		params.Search = kleister.ToPtr(groupUserListArgs.Search)
	}

	if groupUserListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(groupUserListArgs.Sort)
	}

	if groupUserListArgs.Order != "" {
		val, err := kleister.ToListGroupUsersParamsOrder(groupUserListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListGroupUsersParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if groupUserListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(groupUserListArgs.Limit)
	}

	if groupUserListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(groupUserListArgs.Offset)
	}

	resp, err := client.ListGroupUsersWithResponse(
		ccmd.Context(),
		groupUserListArgs.ID,
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
		fmt.Sprintln(groupUserListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Users

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
