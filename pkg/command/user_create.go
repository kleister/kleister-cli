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

type userCreateBind struct {
	Username string
	Password string
	Email    string
	Fullname string
	Active   bool
	Inactive bool
	Admin    bool
	Regular  bool
	Format   string
}

var (
	userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userCreateAction)
		},
		Args: cobra.NoArgs,
	}

	userCreateArgs = userCreateBind{}
)

func init() {
	userCmd.AddCommand(userCreateCmd)

	userCreateCmd.Flags().StringVar(
		&userCreateArgs.Username,
		"username",
		"",
		"Username for user",
	)

	userCreateCmd.Flags().StringVar(
		&userCreateArgs.Password,
		"password",
		"",
		"Password for user",
	)

	userCreateCmd.Flags().StringVar(
		&userCreateArgs.Email,
		"email",
		"",
		"Email for user",
	)

	userCreateCmd.Flags().StringVar(
		&userCreateArgs.Fullname,
		"fullname",
		"",
		"Fullname for user",
	)

	userCreateCmd.Flags().BoolVar(
		&userCreateArgs.Active,
		"active",
		false,
		"Mark user as active",
	)

	userCreateCmd.Flags().BoolVar(
		&userCreateArgs.Inactive,
		"inactive",
		false,
		"Mark user as inactive",
	)

	userCreateCmd.Flags().BoolVar(
		&userCreateArgs.Admin,
		"admin",
		false,
		"Mark user as admin",
	)

	userCreateCmd.Flags().BoolVar(
		&userCreateArgs.Regular,
		"regular",
		false,
		"Mark user as regular",
	)
}

func userCreateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	body := kleister.CreateUserJSONRequestBody{}
	changed := false

	if val := userCreateArgs.Username; val != "" {
		body.Username = kleister.ToPtr(val)
		changed = true
	}

	if val := userCreateArgs.Password; val != "" {
		body.Password = kleister.ToPtr(val)
		changed = true
	}

	if val := userCreateArgs.Email; val != "" {
		body.Email = kleister.ToPtr(val)
		changed = true
	}

	if val := userCreateArgs.Fullname; val != "" {
		body.Fullname = kleister.ToPtr(val)
		changed = true
	}

	if val := userCreateArgs.Active; val {
		body.Active = kleister.ToPtr(true)
		changed = true
	}

	if val := userCreateArgs.Inactive; val {
		body.Active = kleister.ToPtr(false)
		changed = true
	}

	if val := userCreateArgs.Admin; val {
		body.Admin = kleister.ToPtr(true)
		changed = true
	}

	if val := userCreateArgs.Regular; val {
		body.Admin = kleister.ToPtr(false)
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
		fmt.Sprintln(userCreateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.CreateUserWithResponse(
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
