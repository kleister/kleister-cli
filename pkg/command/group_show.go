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

// tmplGroupShow represents a user within details view.
var tmplGroupShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type groupShowBind struct {
	ID     string
	Format string
}

var (
	groupShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show an group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupShowAction)
		},
		Args: cobra.NoArgs,
	}

	groupShowArgs = groupShowBind{}
)

func init() {
	groupCmd.AddCommand(groupShowCmd)

	groupShowCmd.Flags().StringVarP(
		&groupShowArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupShowCmd.Flags().StringVar(
		&groupShowArgs.Format,
		"format",
		tmplGroupShow,
		"Custom output format",
	)
}

func groupShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ShowGroupWithResponse(
		ccmd.Context(),
		groupShowArgs.ID,
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
		fmt.Sprintln(groupShowArgs.Format),
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
