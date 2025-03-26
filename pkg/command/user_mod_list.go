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

type userModListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplUserModList represents a row within user mod listing.
var tmplUserModList = "Slug: \x1b[33m{{ .Mod.Slug }} \x1b[0m" + `
ID: {{ .Mod.Id }}
Name: {{ .Mod.Name }}
Perm: {{ .Perm }}
`

var (
	userModListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned mods for a user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userModListAction)
		},
		Args: cobra.NoArgs,
	}

	userModListArgs = userModListBind{}
)

func init() {
	userModCmd.AddCommand(userModListCmd)

	userModListCmd.Flags().StringVarP(
		&userModListArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userModListCmd.Flags().StringVar(
		&userModListArgs.Format,
		"format",
		tmplUserModList,
		"Custom output format",
	)

	userModListCmd.Flags().StringVar(
		&userModListArgs.Search,
		"search",
		"",
		"Search query",
	)

	userModListCmd.Flags().StringVar(
		&userModListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	userModListCmd.Flags().StringVar(
		&userModListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	userModListCmd.Flags().IntVar(
		&userModListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	userModListCmd.Flags().IntVar(
		&userModListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func userModListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userModListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListUserModsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if userModListArgs.Search != "" {
		params.Search = kleister.ToPtr(userModListArgs.Search)
	}

	if userModListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(userModListArgs.Sort)
	}

	if userModListArgs.Order != "" {
		val, err := kleister.ToListUserModsParamsOrder(userModListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListUserModsParamsOrder) {
			return fmt.Errorf("invalid order attribute")
		}

		params.Order = kleister.ToPtr(val)
	}

	if userModListArgs.Limit != 0 {
		params.Limit = kleister.ToPtr(userModListArgs.Limit)
	}

	if userModListArgs.Offset != 0 {
		params.Offset = kleister.ToPtr(userModListArgs.Offset)
	}

	resp, err := client.ListUserModsWithResponse(
		ccmd.Context(),
		userModListArgs.ID,
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
		fmt.Sprintln(userModListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Mods

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
