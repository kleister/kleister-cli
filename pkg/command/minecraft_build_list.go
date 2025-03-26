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

type minecraftBuildListBind struct {
	ID     string
	Format string
	Search string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

// tmplMinecraftBuildList represents a row within minecraft build listing.
var tmplMinecraftBuildList = "Name: \x1b[33m{{ .Name }} \x1b[0m" + `
Pack: {{ .Pack.Name }}
ID: {{ .Pack.Id }}/{{ .Id }}
`

var (
	minecraftBuildListCmd = &cobra.Command{
		Use:   "list",
		Short: "List assigned builds for a minecraft",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, minecraftBuildListAction)
		},
		Args: cobra.NoArgs,
	}

	minecraftBuildListArgs = minecraftBuildListBind{}
)

func init() {
	minecraftBuildCmd.AddCommand(minecraftBuildListCmd)

	minecraftBuildListCmd.Flags().StringVarP(
		&minecraftBuildListArgs.ID,
		"id",
		"i",
		"",
		"Minecraft ID or slug",
	)

	minecraftBuildListCmd.Flags().StringVar(
		&minecraftBuildListArgs.Format,
		"format",
		tmplMinecraftBuildList,
		"Custom output format",
	)

	minecraftBuildListCmd.Flags().StringVar(
		&minecraftBuildListArgs.Search,
		"search",
		"",
		"Search query",
	)

	minecraftBuildListCmd.Flags().StringVar(
		&minecraftBuildListArgs.Sort,
		"sort",
		"",
		"Sorting column",
	)

	minecraftBuildListCmd.Flags().StringVar(
		&minecraftBuildListArgs.Order,
		"order",
		"asc",
		"Sorting order",
	)

	minecraftBuildListCmd.Flags().IntVar(
		&minecraftBuildListArgs.Limit,
		"limit",
		0,
		"Paging limit",
	)

	minecraftBuildListCmd.Flags().IntVar(
		&minecraftBuildListArgs.Offset,
		"offset",
		0,
		"Paging offset",
	)
}

func minecraftBuildListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if minecraftBuildListArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	params := &kleister.ListMinecraftBuildsParams{
		Limit:  kleister.ToPtr(10000),
		Offset: kleister.ToPtr(0),
	}

	if minecraftBuildListArgs.Search != "" {
		params.Search = kleister.ToPtr(minecraftBuildListArgs.Search)
	}

	if minecraftBuildListArgs.Sort != "" {
		params.Sort = kleister.ToPtr(minecraftBuildListArgs.Sort)
	}

	if minecraftBuildListArgs.Order != "" {
		val, err := kleister.ToListMinecraftBuildsParamsOrder(minecraftBuildListArgs.Order)

		if err != nil && errors.Is(err, kleister.ErrListMinecraftBuildsParamsOrder) {
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

	resp, err := client.ListMinecraftBuildsWithResponse(
		ccmd.Context(),
		minecraftBuildListArgs.ID,
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
		fmt.Sprintln(minecraftBuildListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := resp.JSON200.Builds

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
