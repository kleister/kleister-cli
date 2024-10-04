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

type teamUpdateBind struct {
	ID     string
	Slug   string
	Name   string
	Format string
}

var (
	teamUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an team",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, teamUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	teamUpdateArgs = teamUpdateBind{}
)

func init() {
	teamCmd.AddCommand(teamUpdateCmd)

	teamUpdateCmd.Flags().StringVarP(
		&teamUpdateArgs.ID,
		"id",
		"i",
		"",
		"Team ID or slug",
	)

	teamUpdateCmd.Flags().StringVar(
		&teamUpdateArgs.Slug,
		"slug",
		"",
		"Slug for team",
	)

	teamUpdateCmd.Flags().StringVar(
		&teamUpdateArgs.Name,
		"name",
		"",
		"Name for team",
	)

	teamUpdateCmd.Flags().StringVar(
		&teamUpdateArgs.Format,
		"format",
		tmplTeamShow,
		"Custom output format",
	)
}

func teamUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if teamShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	body := kleister.UpdateTeamJSONRequestBody{}
	changed := false

	if val := teamUpdateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := teamUpdateArgs.Name; val != "" {
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
		fmt.Sprintln(teamUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdateTeamWithResponse(
		ccmd.Context(),
		teamUpdateArgs.ID,
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
