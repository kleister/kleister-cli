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

type modVersionCreateBind struct {
	Mod         string
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
	modVersionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a mod version",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modVersionCreateAction)
		},
		Args: cobra.NoArgs,
	}

	modVersionCreateArgs = modVersionCreateBind{}
)

func init() {
	modVersionCmd.AddCommand(modVersionCreateCmd)

	modVersionCreateCmd.Flags().StringVar(
		&modVersionCreateArgs.Mod,
		"mod",
		"",
		"Mod ID or slug",
	)

	modVersionCreateCmd.Flags().StringVar(
		&modVersionCreateArgs.Name,
		"name",
		"",
		"Name for mod version",
	)

	// TODO File

	modVersionCreateCmd.Flags().BoolVar(
		&modVersionCreateArgs.Private,
		"private",
		false,
		"Mark mod version as private",
	)

	modVersionCreateCmd.Flags().BoolVar(
		&modVersionCreateArgs.Public,
		"public",
		false,
		"Mark mod version as public",
	)

	modVersionCreateCmd.Flags().StringVar(
		&modVersionCreateArgs.Format,
		"format",
		tmplModVersionShow,
		"Custom output format",
	)
}

func modVersionCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modVersionCreateArgs.Mod == "" {
		return fmt.Errorf("you must provide a mod ID or slug")
	}

	body := kleister.CreateVersionJSONRequestBody{}
	changed := false

	if val := modVersionCreateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := modVersionCreateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := modVersionCreateArgs.Public; val {
		body.Public = kleister.ToPtr(true)
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
		fmt.Sprintln(modVersionCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreateVersionWithResponse(
		ccmd.Context(),
		modVersionCreateArgs.Mod,
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
