package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type modCreateBind struct {
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
	modCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an mod",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, modCreateAction)
		},
		Args: cobra.NoArgs,
	}

	modCreateArgs = modCreateBind{}
)

func init() {
	modCmd.AddCommand(modCreateCmd)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Slug,
		"slug",
		"",
		"Slug for mod",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Name,
		"name",
		"",
		"Name for mod",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Side,
		"side",
		"",
		"Side for mod",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Description,
		"description",
		"",
		"Description for mod",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Author,
		"author",
		"",
		"Author for mod",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Website,
		"website",
		"",
		"Website for mod",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Donate,
		"donate",
		"",
		"Donate for mod",
	)

	modCreateCmd.Flags().BoolVar(
		&modCreateArgs.Private,
		"private",
		false,
		"Mark mod as private",
	)

	modCreateCmd.Flags().BoolVar(
		&modCreateArgs.Public,
		"public",
		false,
		"Mark mod as public",
	)

	modCreateCmd.Flags().StringVar(
		&modCreateArgs.Format,
		"format",
		tmplModShow,
		"Custom output format",
	)
}

func modCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	body := kleister.CreateModJSONRequestBody{}
	changed := false

	if val := modCreateArgs.Slug; val != "" {
		body.Slug = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Side; val != "" {
		body.Side = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Description; val != "" {
		body.Description = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Author; val != "" {
		body.Author = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Website; val != "" {
		body.Website = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Donate; val != "" {
		body.Donate = kleister.ToPtr(val)
		changed = true
	}

	if val := modCreateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := modCreateArgs.Public; val {
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
		fmt.Sprintln(modCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreateModWithResponse(
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
