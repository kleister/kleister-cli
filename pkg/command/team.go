package command

import (
	"fmt"
	"os"
	"text/template"

	"github.com/bufbuild/connect-go"
	membersv1 "github.com/kleister/kleister-go/kleister/members/v1"
	teamsv1 "github.com/kleister/kleister-go/kleister/teams/v1"
	"github.com/urfave/cli/v2"
)

// tmplTeamList represents a row within user listing.
var tmplTeamList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
`

// tmplTeamShow represents a user within details view.
var tmplTeamShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .Id }}
Name: {{ .Name }}
Created: {{ .CreatedAt.AsTime }}
Updated: {{ .UpdatedAt.AsTime }}
`

// tmplTeamUserList represents a row within team user listing.
var tmplTeamUserList = "Slug: \x1b[33m{{ .UserSlug }} \x1b[0m" + `
ID: {{ .UserId }}
Username: {{ .UserName }}
`

// Team provides the sub-command for the team API.
func Team() *cli.Command {
	return &cli.Command{
		Name:  "team",
		Usage: "Team related sub-commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List all teams",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplTeamList,
						Usage:  "Custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamList)
				},
			},
			{
				Name:      "show",
				Usage:     "Show a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Team ID or slug",
					},
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplTeamShow,
						Usage:  "Custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamShow)
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "Delete a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Team ID or slug",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamDelete)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Team ID or slug",
					},
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "Provide a slug",
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "",
						Usage: "Provide a name",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamUpdate)
				},
			},
			{
				Name:      "create",
				Usage:     "Create a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "Provide a slug",
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "",
						Usage: "Provide a name",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamCreate)
				},
			},
			{
				Name:  "user",
				Usage: "User assignments",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"ls"},
						Usage:     "List assigned users for a team",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "Team ID or slug",
							},
							&cli.StringFlag{
								Name:   "format",
								Value:  tmplTeamUserList,
								Usage:  "Custom output format",
								Hidden: true,
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserList)
						},
					},
					{
						Name:      "append",
						Usage:     "Append user to team",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "Team ID or slug",
							},
							&cli.StringFlag{
								Name:  "user, u",
								Value: "",
								Usage: "User ID or slug",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserAppend)
						},
					},
					{
						Name:      "remove",
						Aliases:   []string{"rm"},
						Usage:     "Remove user from a team",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "Team ID or slug",
							},
							&cli.StringFlag{
								Name:  "user, u",
								Value: "",
								Usage: "User ID or slug",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserRemove)
						},
					},
				},
			},
		},
	}
}

// TeamList provides the sub-command to list all teams.
func TeamList(c *cli.Context, client *Client) error {
	resp, err := client.Teams.List(
		c.Context,
		&connect.Request[teamsv1.ListRequest]{},
	)

	if err != nil {
		return PrettyError(err)
	}

	if len(resp.Msg.Teams) == 0 {
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

	for _, record := range resp.Msg.Teams {
		if err := tmpl.Execute(os.Stdout, record); err != nil {
			return err
		}
	}

	return nil
}

// TeamShow provides the sub-command to show team details.
func TeamShow(c *cli.Context, client *Client) error {
	resp, err := client.Teams.Show(
		c.Context,
		&connect.Request[teamsv1.ShowRequest]{
			Msg: &teamsv1.ShowRequest{
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

	return tmpl.Execute(os.Stdout, resp.Msg.Team)
}

// TeamDelete provides the sub-command to delete a team.
func TeamDelete(c *cli.Context, client *Client) error {
	resp, err := client.Teams.Delete(
		c.Context,
		&connect.Request[teamsv1.DeleteRequest]{
			Msg: &teamsv1.DeleteRequest{
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

// TeamUpdate provides the sub-command to update a team.
func TeamUpdate(c *cli.Context, client *Client) error {
	resp, err := client.Teams.Show(
		c.Context,
		&connect.Request[teamsv1.ShowRequest]{
			Msg: &teamsv1.ShowRequest{
				Id: GetIdentifierParam(c),
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	record := resp.Msg.Team
	req := &teamsv1.UpdateTeam{}
	changed := false

	if val := c.String("slug"); c.IsSet("slug") && val != record.Slug {
		req.Slug = &val
		changed = true
	}

	if val := c.String("name"); c.IsSet("name") && val != record.Name {
		req.Name = &val
		changed = true
	}

	if changed {
		_, err := client.Teams.Update(
			c.Context,
			&connect.Request[teamsv1.UpdateRequest]{
				Msg: &teamsv1.UpdateRequest{
					Id:   GetIdentifierParam(c),
					Team: req,
				},
			},
		)

		if err != nil {
			return PrettyError(err)
		}

		fmt.Fprintln(os.Stderr, "Successfully updated")
	} else {
		fmt.Fprintln(os.Stderr, "Nothing to update...")
	}

	return nil
}

// TeamCreate provides the sub-command to create a team.
func TeamCreate(c *cli.Context, client *Client) error {
	record := &teamsv1.CreateTeam{}

	if val := c.String("slug"); c.IsSet("slug") && val != "" {
		record.Slug = val
	}

	if val := c.String("name"); c.IsSet("name") && val != "" {
		record.Name = val
	} else {
		return fmt.Errorf("you must provide a name")
	}

	_, err := client.Teams.Create(
		c.Context,
		&connect.Request[teamsv1.CreateRequest]{
			Msg: &teamsv1.CreateRequest{
				Team: record,
			},
		},
	)

	if err != nil {
		return PrettyError(err)
	}

	fmt.Fprintln(os.Stderr, "Successfully created")
	return nil
}

// TeamUserList provides the sub-command to list users of the team.
func TeamUserList(c *cli.Context, client *Client) error {
	resp, err := client.Members.List(
		c.Context,
		&connect.Request[membersv1.ListRequest]{
			Msg: &membersv1.ListRequest{
				Team: GetIdentifierParam(c),
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

// TeamUserAppend provides the sub-command to append a user to the team.
func TeamUserAppend(c *cli.Context, client *Client) error {
	resp, err := client.Members.Append(
		c.Context,
		&connect.Request[membersv1.AppendRequest]{
			Msg: &membersv1.AppendRequest{
				Member: &membersv1.AppendMember{
					Team: GetIdentifierParam(c),
					User: GetUserParam(c),
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

// TeamUserRemove provides the sub-command to remove a user from the team.
func TeamUserRemove(c *cli.Context, client *Client) error {
	resp, err := client.Members.Drop(
		c.Context,
		&connect.Request[membersv1.DropRequest]{
			Msg: &membersv1.DropRequest{
				Member: &membersv1.DropMember{
					Team: GetIdentifierParam(c),
					User: GetUserParam(c),
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
