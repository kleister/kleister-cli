package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modUpdateBind struct {
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
	modUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	modUpdateArgs = modUpdateBind{}
)

func init() {
	modCmd.AddCommand(modUpdateCmd)

	modUpdateCmd.Flags().StringVarP(
		&modUpdateArgs.ID,
		"id",
		"i",
		"",
		"Mod ID or slug",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Slug,
		"slug",
		"",
		"Slug for mod",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Name,
		"name",
		"",
		"Name for mod",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Side,
		"side",
		"",
		"Side for mod",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Description,
		"description",
		"",
		"Description for mod",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Author,
		"author",
		"",
		"Author for mod",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Website,
		"website",
		"",
		"Website for mod",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Donate,
		"donate",
		"",
		"Donate for mod",
	)

	modUpdateCmd.Flags().BoolVar(
		&modUpdateArgs.Private,
		"private",
		false,
		"Mark mod as private",
	)

	modUpdateCmd.Flags().BoolVar(
		&modUpdateArgs.Public,
		"public",
		false,
		"Mark mod as public",
	)

	modUpdateCmd.Flags().StringVar(
		&modUpdateArgs.Format,
		"format",
		tmplModShow,
		"Custom output format",
	)
}

func modUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if modShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	body := kleister.UpdateModJSONRequestBody{}
	changed := false

	if val := modUpdateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Side; val != "" {
		body.Side = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Description; val != "" {
		body.Description = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Author; val != "" {
		body.Author = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Website; val != "" {
		body.Website = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Donate; val != "" {
		body.Donate = kleister.ToPtr(val)
		changed = true
	}

	if val := modUpdateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := modUpdateArgs.Public; val {
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
		fmt.Sprintln(modUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdateModWithResponse(
		ccmd.Context(),
		modUpdateArgs.ID,
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
	case http.StatusNotFound:
		return fmt.Errorf(kleister.FromPtr(resp.JSON404.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
