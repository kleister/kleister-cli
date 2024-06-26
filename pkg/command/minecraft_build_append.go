package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

type minecraftBuildAppendBind struct {
	ID    string
	Pack  string
	Build string
}

var (
	minecraftBuildAppendCmd = &cobra.Command{
		Use:   "append",
		Short: "Append build to minecraft",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, minecraftBuildAppendAction)
		},
		Args: cobra.NoArgs,
	}

	minecraftBuildAppendArgs = minecraftBuildAppendBind{}
)

func init() {
	minecraftBuildCmd.AddCommand(minecraftBuildAppendCmd)

	minecraftBuildAppendCmd.Flags().StringVarP(
		&minecraftBuildAppendArgs.ID,
		"id",
		"i",
		"",
		"Minecraft ID or slug",
	)

	minecraftBuildAppendCmd.Flags().StringVar(
		&minecraftBuildAppendArgs.Pack,
		"pack",
		"",
		"Pack ID or slug",
	)

	minecraftBuildAppendCmd.Flags().StringVar(
		&minecraftBuildAppendArgs.Build,
		"build",
		"",
		"Build ID or slug",
	)
}

func minecraftBuildAppendAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if minecraftBuildAppendArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	if minecraftBuildAppendArgs.Pack == "" {
		return fmt.Errorf("you must provide a pack ID or a slug")
	}

	if minecraftBuildAppendArgs.Build == "" {
		return fmt.Errorf("you must provide a build ID or a slug")
	}

	body := kleister.AttachMinecraftToBuildJSONRequestBody{
		Pack:  minecraftBuildAppendArgs.Pack,
		Build: minecraftBuildAppendArgs.Build,
	}

	resp, err := client.AttachMinecraftToBuildWithResponse(
		ccmd.Context(),
		minecraftBuildAppendArgs.ID,
		body,
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, kleister.FromPtr(resp.JSON200.Message))
	case http.StatusUnprocessableEntity:
		return validationError(resp.JSON422)
	case http.StatusPreconditionFailed:
		return fmt.Errorf(kleister.FromPtr(resp.JSON412.Message))
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
