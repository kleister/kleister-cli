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

type userGroupListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplUserGroupList represents a row within user group listing.
var tmplUserGroupList = "Slug: \x1b[33m{{ .Group.Slug }} \x1b[0m" + `
ID: {{ .Group.Id }}
Name: {{ .Group.Name }}
Perm: {{ .Perm }}
`

var (
	userGroupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned groups for a user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userGroupListAction)
		},
		Args: cobra.NoArgs,
	}

	userGroupListArgs = userGroupListBind{}
)

func init() {
	userGroupCmd.AddCommand(userGroupListCmd)

	userGroupListCmd.Flags().StringVarP(
		&userGroupListArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userGroupListCmd.Flags().StringVar(
		&userGroupListArgs.Format,
		"format",
		tmplUserGroupList,
		"Custom output format",
	)

	userGroupListCmd.Flags().StringVar(
		&userGroupListArgs.Search,
		"search",
		"",
		"Search query",
	)

	userGroupListCmd.Flags().StringVar(
		&userGroupListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	userGroupListCmd.Flags().StringVar(
		&userGroupListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	userGroupListCmd.Flags().IntVar(
		&userGroupListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	userGroupListCmd.Flags().IntVar(
		&userGroupListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func userGroupListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userGroupListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListUserGroupsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if userGroupListArgs.Search != "" {
		params.Search = kleister.ToPtr(userGroupListArgs.Search)
	}

	if userGroupListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(userGroupListArgs.Sort)
	}

	if userGroupListArgs.Order != "" {
		val, err := kleister.ToListUserGroupsParamsOrder(userGroupListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListUserGroupsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if userGroupListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(userGroupListArgs.Limit)
	}

	if userGroupListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(userGroupListArgs.Offset)
	}

	resp, err := client.ListUserGroupsWithResponse(
		ccmd.Context(),
		userGroupListArgs.ID,
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
		fmt.Sprintln(userGroupListArgs.Format),
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
