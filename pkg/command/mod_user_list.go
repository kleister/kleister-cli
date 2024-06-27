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

type modUserListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplModUserList represents a row within mod user listing.
var tmplModUserList = "Slug: \x1b[33m{{ .User.Username }} \x1b[0m" + `
ID: {{ .User.Id }}
Email: {{ .User.Email }}
Perm: {{ .Perm }}
`

var (
	modUserListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned users for a mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modUserListAction)
		},
		Args: cobra.NoArgs,
	}

	modUserListArgs = modUserListBind{}
)

func init() {
	modUserCmd.AddCommand(modUserListCmd)

	modUserListCmd.Flags().StringVarP(
		&modUserListArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modUserListCmd.Flags().StringVar(
		&modUserListArgs.Format,
		"format",
		tmplModUserList,
		"Custom output format",
	)

	modUserListCmd.Flags().StringVar(
		&modUserListArgs.Search,
		"search",
		"",
		"Search query",
	)

	modUserListCmd.Flags().StringVar(
		&modUserListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	modUserListCmd.Flags().StringVar(
		&modUserListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	modUserListCmd.Flags().IntVar(
		&modUserListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	modUserListCmd.Flags().IntVar(
		&modUserListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func modUserListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modUserListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListModUsersParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if minecraftBuildListArgs.Search != "" {
		params.Search = kleister.ToPtr(minecraftBuildListArgs.Search)
	}

	if minecraftBuildListArgs.Sort != "" {
		val, err := kleister.ToListModUsersParamsSort(minecraftBuildListArgs.Sort)

		if err != nil && errors.Is(err, kleister.ErrListModUsersParamsSort) {
			return fmt.Errorf("invalid sort attribute")
		}

		params.Sort = kleister.ToPtr(val)
	}

	if minecraftBuildListArgs.Order != "" {
		val, err := kleister.ToListModUsersParamsOrder(minecraftBuildListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListModUsersParamsOrder) {
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

	resp, err := client.ListModUsersWithResponse(
		ccmd.Context(),
		modUserListArgs.ID,
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
		fmt.Sprintln(modUserListArgs.Format),
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
