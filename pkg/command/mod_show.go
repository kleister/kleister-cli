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

// tmplModShow represents a user within details view.
var tmplModShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
{{ with .Side -}}
Side: {{ . }}
{{ end -}}
{{ with .Description -}}
Description: {{ . }}
{{ end -}}
{{ with .Author -}}
Author: {{ . }}
{{ end -}}
{{ with .Website -}}
Website: {{ . }}
{{ end -}}
{{ with .Donate -}}
Donate: {{ . }}
{{ end -}}
Public: {{ .Public }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type modShowBind struct {
	ID     string
	Format string
}

var (
	modShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show an mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modShowAction)
		},
		Args: cobra.NoArgs,
	}

	modShowArgs = modShowBind{}
)

func init() {
	modCmd.AddCommand(modShowCmd)

	modShowCmd.Flags().StringVarP(
		&modShowArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modShowCmd.Flags().StringVar(
		&modShowArgs.Format,
		"format",
		tmplModShow,
		"Custom output format",
	)
}

func modShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ShowModWithResponse(
		ccmd.Context(),
		modShowArgs.ID,
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
		fmt.Sprintln(modShowArgs.Format),
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
