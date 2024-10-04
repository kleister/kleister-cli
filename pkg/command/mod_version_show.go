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

// tmplModVersionShow represents a user within details view.
var tmplModVersionShow = "Name: \x1b[33m{{ .Name }} \x1b[0m" + `
ID: {{ .Id }}
File: {{ .File }}
Public: {{ .Public }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type modVersionShowBind struct {
	Mod    string
	ID     string
	Format string
}

var (
	modVersionShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a mod version",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modVersionShowAction)
		},
		Args: cobra.NoArgs,
	}

	modVersionShowArgs = modVersionShowBind{}
)

func init() {
	modVersionCmd.AddCommand(modVersionShowCmd)

	modVersionShowCmd.Flags().StringVar(
		&modVersionShowArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	modVersionShowCmd.Flags().StringVarP(
		&modVersionShowArgs.ID,
		"id",
		"i",
		"",
		"Version ID or slug",
	)

	modVersionShowCmd.Flags().StringVar(
		&modVersionShowArgs.Format,
		"format",
		tmplModVersionShow,
		"Custom output format",
	)
}

func modVersionShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modVersionShowArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or slug")
	}

	if modVersionShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or slug")
	}

	resp, err := client.ShowVersionWithResponse(
		ccmd.Context(),
		modVersionShowArgs.Mod,
		modVersionShowArgs.ID,
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
		fmt.Sprintln(modVersionShowArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if err := tmpl.Execute(
			os.Stdout,
			resp.JSON200,
		); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
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
