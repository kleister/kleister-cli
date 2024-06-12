package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

// tmplUserShow represents a user within details view.
var tmplUserShow = "Username: \x1b[33m{{ .Username }} \x1b[0m" + `
ID: {{ .Id }}
Email: {{ .Email }}
Fullname: {{ .Fullname }}
Active: {{ .Active }}
Admin: {{ .Admin }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

type userShowBind struct {
	ID     string
	Format string
}

var (
	userShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show an user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userShowAction)
		},
		Args: cobra.NoArgs,
	}

	userShowArgs = userShowBind{}
)

func init() {
	userCmd.AddCommand(userShowCmd)

	userShowCmd.Flags().StringVarP(
		&userShowArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userShowCmd.Flags().StringVar(
		&userShowArgs.Format,
		"format",
		tmplUserShow,
		"Custom output format",
	)
}

func userShowAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	resp, err := client.ShowUserWithResponse(
		ccmd.Context(),
		userShowArgs.ID,
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
		fmt.Sprintln(userShowArgs.Format),
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
