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

type packCreateBind struct {
	Slug    string
	Name    string
	Website string
	Private bool
	Public  bool
	Format  string
}

var (
	packCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a pack",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packCreateAction)
		},
		Args: cobra.NoArgs,
	}

	packCreateArgs = packCreateBind{}
)

func init() {
	packCmd.AddCommand(packCreateCmd)

	packCreateCmd.Flags().StringVar(
		&packCreateArgs.Slug,
		"slug",
		"",
		"Slug for pack",
	)

	packCreateCmd.Flags().StringVar(
		&packCreateArgs.Name,
		"name",
		"",
		"Name for pack",
	)

	packCreateCmd.Flags().StringVar(
		&packCreateArgs.Website,
		"website",
		"",
		"Website for pack",
	)

	packCreateCmd.Flags().BoolVar(
		&packCreateArgs.Private,
		"private",
		false,
		"Mark pack as private",
	)

	packCreateCmd.Flags().BoolVar(
		&packCreateArgs.Public,
		"public",
		false,
		"Mark pack as public",
	)

	// TODO Back
	// TODO Icon
	// TODO Logo

	packCreateCmd.Flags().StringVar(
		&packCreateArgs.Format,
		"format",
		tmplPackShow,
		"Custom output format",
	)
}

func packCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	body := kleister.CreatePackJSONRequestBody{}
	changed := false

	if val := packCreateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := packCreateArgs.Name; val != "" {
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
		fmt.Sprintln(packCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreatePackWithResponse(
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
