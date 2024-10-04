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

type modVersionUpdateBind struct {
	Mod         string
	ID          string
	Slug        string
	Name        string
	Side        string
	Description string
	Author      string
	Website     string
	Donate      string
	Private     bool
	Public      bool
	Format      string
}

var (
	modVersionUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a mod version",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modVersionUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	modVersionUpdateArgs = modVersionUpdateBind{}
)

func init() {
	modVersionCmd.AddCommand(modVersionUpdateCmd)

	modVersionUpdateCmd.Flags().StringVar(
		&modVersionUpdateArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	modVersionUpdateCmd.Flags().StringVarP(
		&modVersionUpdateArgs.ID,
		"id",
		"i",
		"",
		"Version ID or slug",
	)

	modVersionUpdateCmd.Flags().StringVar(
		&modVersionUpdateArgs.Name,
		"name",
		"",
		"Name for mod version",
	)

	// TODO File

	modVersionUpdateCmd.Flags().BoolVar(
		&modVersionUpdateArgs.Private,
		"private",
		false,
		"Mark mod version as private",
	)

	modVersionUpdateCmd.Flags().BoolVar(
		&modVersionUpdateArgs.Public,
		"public",
		false,
		"Mark mod version as public",
	)

	modVersionUpdateCmd.Flags().StringVar(
		&modVersionUpdateArgs.Format,
		"format",
		tmplModVersionShow,
		"Custom output format",
	)
}

func modVersionUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modVersionUpdateArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or slug")
	}

	if modVersionUpdateArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or slug")
	}

	body := kleister.UpdateVersionJSONRequestBody{}
	changed := false

	if val := modVersionUpdateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := modVersionUpdateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := modVersionUpdateArgs.Public; val {
		body.Public = kleister.ToPtr(true)
		changed = true
	}

	if !changed {
		fmt.Fprintln(os.Stderr, "nothing to update...")
		return nil
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		basicFuncMap,
	).Parse(
		fmt.Sprintln(modVersionUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdateVersionWithResponse(
		ccmd.Context(),
		modVersionUpdateArgs.Mod,
		modVersionUpdateArgs.ID,
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
