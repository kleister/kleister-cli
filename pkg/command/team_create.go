package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type teamCreateBind struct {
	Slug   string
	Name   string
	Format string
}

var (
	teamCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamCreateAction)
		},
		Args: cobra.NoArgs,
	}

	teamCreateArgs = teamCreateBind{}
)

func init() {
	teamCmd.AddCommand(teamCreateCmd)

	teamCreateCmd.Flags().StringVar(
		&teamCreateArgs.Slug,
		"slug",
		"",
		"Slug for team",
	)

	teamCreateCmd.Flags().StringVar(
		&teamCreateArgs.Name,
		"name",
		"",
		"Name for team",
	)

	teamCreateCmd.Flags().StringVar(
		&teamCreateArgs.Format,
		"format",
		tmplTeamShow,
		"Custom output format",
	)
}

func teamCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	body := kleister.CreateTeamJSONRequestBody{}
	changed := false

	if val := teamCreateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := teamCreateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if !changed {
		fmt.Fprintln(os.Stderr, "nothing to create...")
		return nil
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		basicFuncMap,
	).Parse(
		fmt.Sprintln(teamCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreateTeamWithResponse(
		ccmd.Context(),
		body,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if err := tmpl.Execute(
			os.Stdout,
			resp.JSON200,
		); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
	case http.StatusUnprocessableEntity:
		return validationError(resp.JSON422)
	case http.StatusForbidden:
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
