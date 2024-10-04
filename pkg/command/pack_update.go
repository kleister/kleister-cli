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

type packUpdateBind struct {
	ID      string
	Slug    string
	Name    string
	Website string
	Private bool
	Public  bool
	Format  string
}

var (
	packUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	packUpdateArgs = packUpdateBind{}
)

func init() {
	packCmd.AddCommand(packUpdateCmd)

	packUpdateCmd.Flags().StringVarP(
		&packUpdateArgs.ID,
		"id",
		"i",
		"",
		"Pack ID or slug",
	)

	packUpdateCmd.Flags().StringVar(
		&packUpdateArgs.Slug,
		"slug",
		"",
		"Slug for pack",
	)

	packUpdateCmd.Flags().StringVar(
		&packUpdateArgs.Name,
		"name",
		"",
		"Name for pack",
	)

	packUpdateCmd.Flags().StringVar(
		&packUpdateArgs.Website,
		"website",
		"",
		"Website for pack",
	)

	packUpdateCmd.Flags().BoolVar(
		&packUpdateArgs.Private,
		"private",
		false,
		"Mark pack as private",
	)

	packUpdateCmd.Flags().BoolVar(
		&packUpdateArgs.Public,
		"public",
		false,
		"Mark pack as public",
	)

	// TODO Back
	// TODO Icon
	// TODO Logo

	packUpdateCmd.Flags().StringVar(
		&packUpdateArgs.Format,
		"format",
		tmplPackShow,
		"Custom output format",
	)
}

func packUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	body := kleister.UpdatePackJSONRequestBody{}
	changed := false

	if val := packUpdateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := packUpdateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := packCreateArgs.Website; val != "" {
		body.Website = kleister.ToPtr(val)
		changed = true
	}

	if val := packCreateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := packCreateArgs.Public; val {
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
		fmt.Sprintln(packUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdatePackWithResponse(
		ccmd.Context(),
		packUpdateArgs.ID,
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
