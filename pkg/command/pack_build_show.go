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

// tmplPackBuildShow represents a user within details view.
var tmplPackBuildShow = "Name: \x1b[33m{{ .Name }} \x1b[0m" + `
ID: {{ .Id }}
{{ with .Minecraft -}}
Minecraft: {{ .Name }}
{{ end -}}
{{ with .Forge -}}
Forge: {{ .Name }}
{{ end -}}
{{ with .Neoforge -}}
Neoforge: {{ .Name }}
{{ end -}}
{{ with .Quilt -}}
Quilt: {{ .Name }}
{{ end -}}
{{ with .Fabric -}}
Fabric: {{ .Name }}
{{ end -}}
{{ with .Java -}}
Java: {{ . }}
{{ end -}}
{{ with .Memory -}}
Memory: {{ . }}
{{ end -}}
{{ with .Latest -}}
Latest: {{ . }}
{{ end -}}
{{ with .Recommended -}}
Recommended: {{ . }}
{{ end -}}
Public: {{ .Public }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type packBuildShowBind struct {
	Pack   string
	ID     string
	Format string
}

var (
	packBuildShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a pack build",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packBuildShowAction)
		},
		Args: cobra.NoArgs,
	}

	packBuildShowArgs = packBuildShowBind{}
)

func init() {
	packBuildCmd.AddCommand(packBuildShowCmd)

	packBuildShowCmd.Flags().StringVar(
		&packBuildShowArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	packBuildShowCmd.Flags().StringVarP(
		&packBuildShowArgs.ID,
		"id",
		"i",
		"",
		"Build ID or slug",
	)

	packBuildShowCmd.Flags().StringVar(
		&packBuildShowArgs.Format,
		"format",
		tmplPackBuildShow,
		"Custom output format",
	)
}

func packBuildShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packBuildShowArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or slug")
	}

	if packBuildShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or slug")
	}

	resp, err := client.ShowBuildWithResponse(
		ccmd.Context(),
		packBuildShowArgs.Pack,
		packBuildShowArgs.ID,
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
		fmt.Sprintln(packBuildShowArgs.Format),
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
