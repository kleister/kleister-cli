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

type userUpdateBind struct {
	ID       string
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
	userUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an user",
		Run: func(ccmd *cobra.Command, args []string) {
			Handle(ccmd, args, userUpdateAction)
		},
		Args: cobra.NoArgs,
	}

	userUpdateArgs = userUpdateBind{}
)

func init() {
	userCmd.AddCommand(userUpdateCmd)

	userUpdateCmd.Flags().StringVarP(
		&userUpdateArgs.ID,
		"id",
		"i",
		"",
		"User ID or slug",
	)

	userUpdateCmd.Flags().StringVar(
		&userUpdateArgs.Username,
		"username",
		"",
		"Username for user",
	)

	userUpdateCmd.Flags().StringVar(
		&userUpdateArgs.Password,
		"password",
		"",
		"Password for user",
	)

	userUpdateCmd.Flags().StringVar(
		&userUpdateArgs.Email,
		"email",
		"",
		"Email for user",
	)

	userUpdateCmd.Flags().StringVar(
		&userUpdateArgs.Fullname,
		"fullname",
		"",
		"Fullname for user",
	)

	userUpdateCmd.Flags().BoolVar(
		&userUpdateArgs.Active,
		"active",
		false,
		"Mark user as active",
	)

	userUpdateCmd.Flags().BoolVar(
		&userUpdateArgs.Inactive,
		"inactive",
		false,
		"Mark user as inactive",
	)

	userUpdateCmd.Flags().BoolVar(
		&userUpdateArgs.Admin,
		"admin",
		false,
		"Mark user as admin",
	)

	userUpdateCmd.Flags().BoolVar(
		&userUpdateArgs.Regular,
		"regular",
		false,
		"Mark user as regular",
	)
}

func userUpdateAction(ccmd *cobra.Command, _ []string, client *Client) error {
	if userShowArgs.ID == "" {
		return fmt.Errorf("you must provide an ID or a slug")
	}

	body := kleister.UpdateUserJSONRequestBody{}
	changed := false

	if val := userUpdateArgs.Username; val != "" {
		body.Username = kleister.ToPtr(val)
		changed = true
	}

	if val := userUpdateArgs.Password; val != "" {
		body.Password = kleister.ToPtr(val)
		changed = true
	}

	if val := userUpdateArgs.Email; val != "" {
		body.Email = kleister.ToPtr(val)
		changed = true
	}

	if val := userUpdateArgs.Fullname; val != "" {
		body.Fullname = kleister.ToPtr(val)
		changed = true
	}

	if val := userUpdateArgs.Active; val {
		body.Active = kleister.ToPtr(true)
		changed = true
	}

	if val := userUpdateArgs.Inactive; val {
		body.Active = kleister.ToPtr(false)
		changed = true
	}

	if val := userUpdateArgs.Admin; val {
		body.Admin = kleister.ToPtr(true)
		changed = true
	}

	if val := userUpdateArgs.Regular; val {
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
		fmt.Sprintln(userUpdateArgs.Format),
	)

	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	resp, err := client.UpdateUserWithResponse(
		ccmd.Context(),
		userUpdateArgs.ID,
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
