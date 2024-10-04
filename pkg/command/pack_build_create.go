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

type packBuildCreateBind struct {
	Pack      string
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
	packBuildCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a pack build",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, packBuildCreateAction)
		},
		Args: cobra.NoArgs,
	}

	packBuildCreateArgs = packBuildCreateBind{}
)

func init() {
	packBuildCmd.AddCommand(packBuildCreateCmd)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Name,
		"name",
		"",
		"Name for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Minecraft,
		"minecraft",
		"",
		"Minecraft for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Forge,
		"forge",
		"",
		"Forge for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Neoforge,
		"neoforge",
		"",
		"Neoforge for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Quilt,
		"quilt",
		"",
		"Quilt for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Fabric,
		"fabric",
		"",
		"Fabric for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Java,
		"java",
		"",
		"Java for pack build",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Memory,
		"memory",
		"",
		"Memory for pack build",
	)

	packBuildCreateCmd.Flags().BoolVar(
		&packBuildCreateArgs.Private,
		"private",
		false,
		"Mark pack build as private",
	)

	packBuildCreateCmd.Flags().BoolVar(
		&packBuildCreateArgs.Public,
		"public",
		false,
		"Mark pack build as public",
	)

	packBuildCreateCmd.Flags().StringVar(
		&packBuildCreateArgs.Format,
		"format",
		tmplPackBuildShow,
		"Custom output format",
	)
}

func packBuildCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if packBuildCreateArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or slug")
	}

	body := kleister.CreateBuildJSONRequestBody{}
	changed := false

	if val := packBuildCreateArgs.Name; val != "" {
		body.Name = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Minecraft; val != "" {
		body.MinecraftId = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Forge; val != "" {
		body.ForgeId = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Neoforge; val != "" {
		body.NeoforgeId = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Quilt; val != "" {
		body.QuiltId = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Fabric; val != "" {
		body.FabricId = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Java; val != "" {
		body.Java = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Memory; val != "" {
		body.Memory = kleister.ToPtr(val)
		changed = true
	}

	if val := packBuildCreateArgs.Private; val {
		body.Public = kleister.ToPtr(false)
		changed = true
	}

	if val := packBuildCreateArgs.Public; val {
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
		fmt.Sprintln(packBuildCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreateBuildWithResponse(
		ccmd.Context(),
		packBuildCreateArgs.Pack,
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
