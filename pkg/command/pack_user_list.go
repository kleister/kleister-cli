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

type packUserListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplPackUserList represents a row within pack user listing.
var tmplPackUserList = "Slug: \x1b[33m{{ .User.Username }} \x1b[0m" + `
ID: {{ .User.Id }}
Email: {{ .User.Email }}
Perm: {{ .Perm }}
`

var (
	packUserListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned users for a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packUserListAction)
		},
		Args: cobra.NoArgs,
	}

	packUserListArgs = packUserListBind{}
)

func init() {
	packUserCmd.AddCommand(packUserListCmd)

	packUserListCmd.Flags().StringVarP(
		&packUserListArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packUserListCmd.Flags().StringVar(
		&packUserListArgs.Format,
		"format",
		tmplPackUserList,
		"Custom output format",
	)

	packUserListCmd.Flags().StringVar(
		&packUserListArgs.Search,
		"search",
		"",
		"Search query",
	)

	packUserListCmd.Flags().StringVar(
		&packUserListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	packUserListCmd.Flags().StringVar(
		&packUserListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	packUserListCmd.Flags().IntVar(
		&packUserListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	packUserListCmd.Flags().IntVar(
		&packUserListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func packUserListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packUserListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListPackUsersParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if minecraftBuildListArgs.Search != "" {
		params.Search = kleister.ToPtr(minecraftBuildListArgs.Search)
	}

	if minecraftBuildListArgs.Sort != "" {
		val, err := kleister.ToListPackUsersParamsSort(minecraftBuildListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListPackUsersParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if minecraftBuildListArgs.Order != "" {
		val, err := kleister.ToListPackUsersParamsOrder(minecraftBuildListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListPackUsersParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if minecraftBuildListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(minecraftBuildListArgs.Limit)
	}

	if minecraftBuildListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(minecraftBuildListArgs.Offset)
	}

	resp, err := client.ListPackUsersWithResponse(
		ccmd.Context(),
		packUserListArgs.ID,
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
		fmt.Sprintln(packUserListArgs.Format),
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
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusNotFound:
		return fmt.Errorf(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
