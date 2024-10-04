package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type minecraftBuildRemoveBind struct {
	ID    string
	Pack  string
	Build string
}

var (
	minecraftBuildRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove build from minecraft",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, minecraftBuildRemoveAction)
		},
		Args: cobra.NoArgs,
	}

	minecraftBuildRemoveArgs = minecraftBuildRemoveBind{}
)

func init() {
	minecraftBuildCmd.AddCommand(minecraftBuildRemoveCmd)

	minecraftBuildRemoveCmd.Flags().StringVarP(
		&minecraftBuildRemoveArgs.ID,
		"id",
		"i",
		"",
		"Minecraft ID or slug",
	)

	minecraftBuildRemoveCmd.Flags().StringVar(
		&minecraftBuildRemoveArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	minecraftBuildRemoveCmd.Flags().StringVar(
		&minecraftBuildRemoveArgs.Build,
		"build",
		"",
		"Build ID or slug",
	)
}

func minecraftBuildRemoveAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if minecraftBuildRemoveArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if minecraftBuildRemoveArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	if minecraftBuildRemoveArgs.Build == "" {
		return fmt.Errorf("you must provide a build ID or a slug")
	}

	resp, err := client.DeleteMinecraftFromBuildWithResponse(
		ccmd.Context(),
		minecraftBuildRemoveArgs.ID,
		kleister.DeleteMinecraftFromBuildJSONRequestBody{
			Pack:  minecraftBuildRemoveArgs.Pack,
			Build: minecraftBuildRemoveArgs.Build,
		},
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusPreconditionFailed:
		return errors.New(kleister.FromPtr(resp.JSON412.Message))
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
