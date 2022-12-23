package command

import (
	"fmt"
	"os"
	"text/template"

	"github.com/bufbuild/connect-go"
	membersv1 "github.com/kleister/kleister-go/kleister/members/v1"
	usersv1 "github.com/kleister/kleister-go/kleister/users/v1"
	"github.com/urfave/cli/v2"
)

// tmplUserList represents a row within user listing.
var tmplUserList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Username: {{ .Username }}
`

// tmplUserShow represents a user within details view.
var tmplUserShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Username: {{ .Username }}
Email: {{ .Email }}
{{- with .Firstname }}
Firstname: {{ .Firstname }}
{{ end }}
{{- with .Lastname }}
Lastname: {{ .Lastname }}
{{ end }}
Active: {{ .Active }}
Admin: {{ .Admin }}
Created: {{ .CreatedAt.AsTime }}
Updated: {{ .UpdatedAt.AsTime }}
`

// tmplUserTeamList represents a row within user team listing.
var tmplUserTeamList = "Slug: \x1b[33m{{ .TeamSlug }} \x1b[0m" + `
ID: {{ .TeamId }}
Name: {{ .TeamName }}
`

// User provides the sub-command for the user API.
func User() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Usage: "User related sub-commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List all users",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplUserList,
						Usage:  "Custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserList)
				},
			},
			{
				Name:      "show",
				Usage:     "Show an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "User ID or slug",
					},
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplUserShow,
						Usage:  "Custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserShow)
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "Delete an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "User ID or slug",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserDelete)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "User id or slug",
					},
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "Provide a slug",
					},
					&cli.StringFlag{
						Name:  "username",
						Value: "",
						Usage: "Provide an username",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
						Usage: "Provide a password",
					},
					&cli.StringFlag{
						Name:  "email",
						Value: "",
						Usage: "Provide an email",
					},
					&cli.StringFlag{
						Name:  "firstname",
						Value: "",
						Usage: "Provide a firstname",
					},
					&cli.StringFlag{
						Name:  "lastname",
						Value: "",
						Usage: "Provide a lastname",
					},
					&cli.BoolFlag{
						Name:  "active",
						Usage: "Mark user as active",
					},
					&cli.BoolFlag{
						Name:  "admin",
						Usage: "Mark user as admin",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserUpdate)
				},
			},
			{
				Name:      "create",
				Usage:     "Create an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "Provide a slug",
					},
					&cli.StringFlag{
						Name:  "username",
						Value: "",
						Usage: "Provide an username",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
						Usage: "Provide a password",
					},
					&cli.StringFlag{
						Name:  "email",
						Value: "",
						Usage: "Provide an email",
					},
					&cli.StringFlag{
						Name:  "firstname",
						Value: "",
						Usage: "Provide a firstname",
					},
					&cli.StringFlag{
						Name:  "lastname",
						Value: "",
						Usage: "Provide a lastname",
					},
					&cli.BoolFlag{
						Name:  "active",
						Usage: "Mark user as active",
					},
					&cli.BoolFlag{
						Name:  "admin",
						Usage: "Mark user as admin",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserCreate)
				},
			},
			{
				Name:  "team",
				Usage: "Team assignments",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"ls"},
						Usage:     "List assigned teams for a user",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "User ID or slug",
							},
							&cli.StringFlag{
								Name:   "format",
								Value:  tmplUserTeamList,
								Usage:  "Custom output format",
								Hidden: true,
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamList)
						},
					},
					{
						Name:      "append",
						Usage:     "Append team to user",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "User ID or slug",
							},
							&cli.StringFlag{
								Name:  "team, t",
								Value: "",
								Usage: "Team ID or slug",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamAppend)
						},
					},
					{
						Name:      "remove",
						Aliases:   []string{"rm"},
						Usage:     "Remove team from user",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "User ID or slug",
							},
							&cli.StringFlag{
								Name:  "team, t",
								Value: "",
								Usage: "Team ID or slug",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamRemove)
						},
					},
				},
			},
		},
	}
}

// UserList provides the sub-command to list all users.
func UserList(c *cli.Context, client *Client) error {
	resp, err := client.Users.List(
		c.Context,
		&connect.Request[usersv1.ListRequest]{},
	)

	if err != nil {
		return PrettyError(err)
	}

	if len(resp.Msg.Users) == 0 {
		fmt.Fprintln(os.Stderr, "Empty result")
		return nil
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		sprigFuncMap,
	).Parse(
		fmt.Sprintln(c.String("format")),
	)

	if err != nil {
		return err
	}

	for _, record := range resp.Msg.Users {
		if err := tmpl.Execute(os.Stdout, record); err != nil {
			return err
		}
	}

	return nil
}

// UserShow provides the sub-command to show user details.
func UserShow(c *cli.Context, client *Client) error {
	resp, err := client.Users.Show(
		c.Context,
		&connect.Request[usersv1.ShowRequest]{
			Msg: &usersv1.ShowRequest{
				Id: GetIdentifierParam(c),
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		sprigFuncMap,
	).Parse(
		fmt.Sprintln(c.String("format")),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, resp.Msg.User)
}

// UserDelete provides the sub-command to delete a user.
func UserDelete(c *cli.Context, client *Client) error {
	resp, err := client.Users.Delete(
		c.Context,
		&connect.Request[usersv1.DeleteRequest]{
			Msg: &usersv1.DeleteRequest{
				Id: GetIdentifierParam(c),
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	fmt.Fprintln(os.Stderr, resp.Msg.Message)
	return nil
}

// UserUpdate provides the sub-command to update a user.
func UserUpdate(c *cli.Context, client *Client) error {
	resp, err := client.Users.Show(
		c.Context,
		&connect.Request[usersv1.ShowRequest]{
			Msg: &usersv1.ShowRequest{
				Id: GetIdentifierParam(c),
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	record := resp.Msg.User
	req := &usersv1.UpdateUser{}
	changed := false

	if val := c.String("slug"); c.IsSet("slug") && val != record.Slug {
		req.Slug = &val
		changed = true
	}

	if val := c.String("username"); c.IsSet("username") && val != record.Username {
		req.Username = &val
		changed = true
	}

	if val := c.String("password"); c.IsSet("password") {
		req.Password = &val
		changed = true
	}

	if val := c.String("email"); c.IsSet("email") && val != record.Email {
		req.Email = &val
		changed = true
	}

	if val := c.String("firstname"); c.IsSet("firstname") && val != record.Firstname {
		req.Firstname = &val
		changed = true
	}

	if val := c.String("lastname"); c.IsSet("lastname") && val != record.Lastname {
		req.Lastname = &val
		changed = true
	}

	if val := c.Bool("active"); c.IsSet("active") {
		req.Active = &val
		changed = true
	}

	if val := c.Bool("admin"); c.IsSet("admin") {
		req.Admin = &val
		changed = true
	}

	if changed {
		_, err := client.Users.Update(
			c.Context,
			&connect.Request[usersv1.UpdateRequest]{
				Msg: &usersv1.UpdateRequest{
					Id:   GetIdentifierParam(c),
					User: req,
				},
			},
		)

		if err != nil {
			return PrettyError(err)
		}

		fmt.Fprintln(os.Stderr, "Successfully update")
	} else {
		fmt.Fprintln(os.Stderr, "Nothing to update...")
	}

	return nil
}

// UserCreate provides the sub-command to create a user.
func UserCreate(c *cli.Context, client *Client) error {
	record := &usersv1.CreateUser{}

	if val := c.String("slug"); c.IsSet("slug") && val != "" {
		record.Slug = val
	}

	if val := c.String("username"); c.IsSet("username") && val != "" {
		record.Username = val
	} else {
		return fmt.Errorf("you must provide an username")
	}

	if val := c.String("password"); c.IsSet("password") && val != "" {
		record.Password = val
	} else {
		return fmt.Errorf("you must provide a password")
	}

	if val := c.String("email"); c.IsSet("email") && val != "" {
		record.Email = val
	} else {
		return fmt.Errorf("you must provide an email")
	}

	if val := c.String("firstname"); c.IsSet("firstname") && val != "" {
		record.Firstname = val
	}

	if val := c.String("lastname"); c.IsSet("lastname") && val != "" {
		record.Lastname = val
	}

	if c.IsSet("active") {
		record.Active = c.Bool("active")
	}

	if c.IsSet("admin") {
		record.Admin = c.Bool("admin")
	}

	_, err := client.Users.Create(
		c.Context,
		&connect.Request[usersv1.CreateRequest]{
			Msg: &usersv1.CreateRequest{
				User: record,
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	fmt.Fprintln(os.Stderr, "Successfully created")
	return nil
}

// UserTeamList provides the sub-command to list teams of the user.
func UserTeamList(c *cli.Context, client *Client) error {
	resp, err := client.Members.List(
		c.Context,
		&connect.Request[membersv1.ListRequest]{
			Msg: &membersv1.ListRequest{
				User: GetIdentifierParam(c),
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	if len(resp.Msg.Members) == 0 {
		fmt.Fprintln(os.Stderr, "Empty result")
		return nil
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		sprigFuncMap,
	).Parse(
		fmt.Sprintln(c.String("format")),
	)

	if err != nil {
		return err
	}

	for _, record := range resp.Msg.Members {
		if err := tmpl.Execute(os.Stdout, record); err != nil {
			return err
		}
	}

	return nil
}

// UserTeamAppend provides the sub-command to append a team to the user.
func UserTeamAppend(c *cli.Context, client *Client) error {
	resp, err := client.Members.Append(
		c.Context,
		&connect.Request[membersv1.AppendRequest]{
			Msg: &membersv1.AppendRequest{
				Member: &membersv1.AppendMember{
					User: GetIdentifierParam(c),
					Team: GetTeamParam(c),
				},
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	fmt.Fprintln(os.Stderr, resp.Msg.Message)
	return nil
}

// UserTeamRemove provides the sub-command to remove a team from the user.
func UserTeamRemove(c *cli.Context, client *Client) error {
	resp, err := client.Members.Drop(
		c.Context,
		&connect.Request[membersv1.DropRequest]{
			Msg: &membersv1.DropRequest{
				Member: &membersv1.DropMember{
					User: GetIdentifierParam(c),
					Team: GetTeamParam(c),
				},
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	fmt.Fprintln(os.Stderr, resp.Msg.Message)
	return nil
}
