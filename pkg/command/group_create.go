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

type groupCreateBind struct {
	Slug   string
	Name   string
	Format string
}

var (
	groupCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an group",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, groupCreateAction)
		},
		Args: cobra.NoArgs,
	}

	groupCreateArgs = groupCreateBind{}
)

func init() {
	groupCmd.AddCommand(groupCreateCmd)

	groupCreateCmd.Flags().StringVar(
		&groupCreateArgs.Slug,
		"slug",
		"",
		"Slug for group",
	)

	groupCreateCmd.Flags().StringVar(
		&groupCreateArgs.Name,
		"name",
		"",
		"Name for group",
	)

	groupCreateCmd.Flags().StringVar(
		&groupCreateArgs.Format,
		"format",
		tmplGroupShow,
		"Custom output format",
	)
}

func groupCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	body := kleister.CreateGroupJSONRequestBody{}
	changed := false

	if val := groupCreateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := groupCreateArgs.Name; val != "" {
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
		fmt.Sprintln(groupCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreateGroupWithResponse(
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
		return errors.New(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusInternalServerError:
		return errors.New(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
