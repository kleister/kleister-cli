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

// tmplMinecraftList represents a row within user listing.
var tmplMinecraftList = "Name: \x1b[33m{{ .Name }} \x1b[0m" + `
ID: {{ .Id }}
Type: {{ .Type }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type minecraftListBind struct {
	Format string
	Search string
}

var (
	minecraftListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all minecrafts",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, minecraftListAction)
		},
		Args: cobra.NoArgs,
	}

	minecraftListArgs = minecraftListBind{}
)

func init() {
	minecraftCmd.AddCommand(minecraftListCmd)

	minecraftListCmd.Flags().StringVar(
		&minecraftListArgs.Format,
		"format",
		tmplMinecraftList,
		"Custom output format",
	)

	minecraftListCmd.Flags().StringVar(
		&minecraftListArgs.Search,
		"search",
		"",
		"Search term for versions",
	)
}

func minecraftListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	params := &kleister.ListMinecraftsParams{}

	if minecraftListArgs.Search != "" {
		params.Search = kleister.ToPtr(minecraftListArgs.Search)
	}

	resp, err := client.ListMinecraftsWithResponse(
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
		fmt.Sprintln(minecraftListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Versions)

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
