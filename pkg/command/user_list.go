package command

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/kleister/kleister-go/kleister"
	"github.com/spf13/cobra"
)

// tmplUserList represents a row within user listing.
var tmplUserList = "Username: \x1b[33m{{ .Username }} \x1b[0m" + `
ID: {{ .Id }}
Email: {{ .Email }}
`

type userListBind struct {
	Format string
}

var (
	userListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userListAction)
		},
		Args: cobra.NoArgs,
	}

	userListArgs = userListBind{}
)

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().StringVar(
		&userListArgs.Format,
		"format",
		tmplUserList,
		"Custom output format",
	)
}

func userListAction(ccmd *cobra.Command, _ []string, client *Client) error {
	resp, err := client.ListUsersWithResponse(
		ccmd.Context(),
		&kleister.ListUsersParams{
			Limit:  kleister.ToPtr(10000),
			Offset: kleister.ToPtr(0),
		},
	)

	if err != nil {
		return err
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		basicFuncMap,
	).Parse(
		fmt.Sprintln(userListArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		records := kleister.FromPtr(resp.JSON200.Users)

		if len(records) == 0 {
			fmt.Fprintln(os.Stderr, "Empty result")
			return nil
		}

		for _, record := range records {
			if err := tmpl.Execute(
				os.Stdout,
				record,
			); err != nil {
				return fmt.Errorf("failed to render template: %w", err)
			}
		}
	case http.StatusForbidden:
		return fmt.Errorf(kleister.FromPtr(resp.JSON403.Message))
	case http.StatusInternalServerError:
		return fmt.Errorf(kleister.FromPtr(resp.JSON500.Message))
	default:
		return fmt.Errorf("unknown api response")
	}

	return nil
}
