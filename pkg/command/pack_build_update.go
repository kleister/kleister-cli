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

type packBuildUpdateBind struct {
	Pack      string
	ID        string
	Name      string
	Minecraft string
	Forge     string
	Neoforge  string
	Quilt     string
	Fabric    string
	Java      string
	Memory    string
	Private   bool
	Public    bool
	Format    string
}

var (
	packBuildUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a mod version",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packBuildUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	packBuildUpdateArgs = packBuildUpdateBind{}
)

func init() {
	packBuildCmd.AddCommand(packBuildUpdateCmd)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	packBuildUpdateCmd.Flags().StringVarP(
		&packBuildUpdateArgs.ID,
		"id",
		"i",
		"",
		"Build ID or slug",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Name,
		"name",
		"",
		"Name for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Minecraft,
		"minecraft",
		"",
		"Minecraft for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Forge,
		"forge",
		"",
		"Forge for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Neoforge,
		"neoforge",
		"",
		"Neoforge for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Quilt,
		"quilt",
		"",
		"Quilt for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Fabric,
		"fabric",
		"",
		"Fabric for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Java,
		"java",
		"",
		"Java for pack build",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Memory,
		"memory",
		"",
		"Memory for pack build",
	)

	packBuildUpdateCmd.Flags().BoolVar(
		&packBuildUpdateArgs.Private,
		"private",
		false,
		"Mark pack build as private",
	)

	packBuildUpdateCmd.Flags().BoolVar(
		&packBuildUpdateArgs.Public,
		"public",
		false,
		"Mark pack build as public",
	)

	packBuildUpdateCmd.Flags().StringVar(
		&packBuildUpdateArgs.Format,
		"format",
		tmplPackBuildShow,
		"Custom output format",
	)
}

func packBuildUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packBuildUpdateArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or slug")
	}

	if packBuildUpdateArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or slug")
	}

	body := kleister.UpdateBuildJSONRequestBody{}
	changed := false

	if val := packBuildUpdateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Minecraft; val != "" {
		body.MinecraftID = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Forge; val != "" {
		body.ForgeID = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Neoforge; val != "" {
		body.NeoforgeID = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Quilt; val != "" {
		body.QuiltID = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Fabric; val != "" {
		body.FabricID = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Java; val != "" {
		body.Java = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Memory; val != "" {
		body.Memory = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildUpdateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := packBuildUpdateArgs.Public; val {
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
		fmt.Sprintln(packBuildUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdateBuildWithResponse(
		ccmd.Context(),
		packBuildUpdateArgs.Pack,
		packBuildUpdateArgs.ID,
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
