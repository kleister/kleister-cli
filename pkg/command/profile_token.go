package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

// tmplProfileToken represents a permanent login token.
var tmplProfileToken = "Token: \x1b[33m{{ .Token }} \x1b[0m" + `
`

type profileTokenBind struct {
	Format string
}

var (
	profileTokenCmd = &cobra.Command{
		Use:   "token",
		Short: "Show your token",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, profileTokenAction)
		},
		Args: cobra.NoArgs,
	}

	profileTokenArgs = profileTokenBind{}
)

func init() {
	profileCmd.AddCommand(profileTokenCmd)

	profileTokenCmd.Flags().StringVar(
		&profileTokenArgs.Format,
		"format",
		tmplProfileToken,
		"Custom output format",
	)
}

func profileTokenAction(ccmd *cobra.Command, _ []string, client *Client) error {
	resp, err := client.TokenProfileWithResponse(
		ccmd.Context(),
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
		fmt.Sprintln(profileTokenArgs.Format),
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
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
