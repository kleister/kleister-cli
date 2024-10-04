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

// tmplPackShow represents a user within details view.
var tmplPackShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
{{ with .Website -}}
Website: {{ . }}
{{ end -}}
Public: {{ .Public }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type packShowBind struct {
	ID     string
	Format string
}

var (
	packShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packShowAction)
		},
		Args: cobra.NoArgs,
	}

	packShowArgs = packShowBind{}
)

func init() {
	packCmd.AddCommand(packShowCmd)

	packShowCmd.Flags().StringVarP(
		&packShowArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packShowCmd.Flags().StringVar(
		&packShowArgs.Format,
		"format",
		tmplPackShow,
		"Custom output format",
	)
}

func packShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ShowPackWithResponse(
		ccmd.Context(),
		packShowArgs.ID,
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
		fmt.Sprintln(packShowArgs.Format),
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
