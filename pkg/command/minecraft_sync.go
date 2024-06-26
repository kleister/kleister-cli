package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

var (
	minecraftSyncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync minecraft upstream versions",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, minecraftSyncAction)
		},
		Args: cobra.NoArgs,
	}
)

func init() {
	minecraftCmd.AddCommand(minecraftSyncCmd)
}

func minecraftSyncAction(ccmd *cobra.Command, _ []string, client *Client) error {
	resp, err := client.UpdateMinecraftWithResponse(
		ccmd.Context(),
	)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		fmt.Fprintln(os.Stderr, "successfully synchronized")
	case http.StatusForbidden:
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusServiceUnavailable:
		return fmt.Errorf(kleister.FromPtr(resp.JSON503.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
