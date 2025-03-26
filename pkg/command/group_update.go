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

type groupUpdateBind struct {
	ID     string
	Slug   string
	Name   string
	Format string
}

var (
	groupUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	groupUpdateArgs = groupUpdateBind{}
)

func init() {
	groupCmd.AddCommand(groupUpdateCmd)

	groupUpdateCmd.Flags().StringVarP(
		&groupUpdateArgs.ID,
		"id",
		"i",
		"",
		"Group ID or slug",
	)

	groupUpdateCmd.Flags().StringVar(
		&groupUpdateArgs.Slug,
		"slug",
		"",
		"Slug for group",
	)

	groupUpdateCmd.Flags().StringVar(
		&groupUpdateArgs.Name,
		"name",
		"",
		"Name for group",
	)

	groupUpdateCmd.Flags().StringVar(
		&groupUpdateArgs.Format,
		"format",
		tmplGroupShow,
		"Custom output format",
	)
}

func groupUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if groupShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	body := kleister.UpdateGroupJSONRequestBody{}
	changed := false

	if val := groupUpdateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := groupUpdateArgs.Name; val != "" {
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
		fmt.Sprintln(groupUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdateGroupWithResponse(
		ccmd.Context(),
		groupUpdateArgs.ID,
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
