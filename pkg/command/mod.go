package command

// import (
// 	"encoding/json"
// 	"encoding/xml"
// 	"fmt"
// 	"os"
// 	"text/template"

// 	"github.com/kleister/kleister-go/kleister"
// 	"gopkg.in/urfave/cli.v2"
// )

// // tmplModList represents a row within forge listing.
// var tmplModList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// `

// // tmplModShow represents a mod within details view.
// var tmplModShow = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}{{with .Side}}
// Side: {{ . }}{{end}}{{with .Description}}
// Description: {{ . }}{{end}}{{with .Author}}
// Author: {{ . }}{{end}}{{with .Website}}
// Website: {{ . }}{{end}}{{with .Donate}}
// Donate: {{ . }}{{end}}{{with .Users}}
// Users: {{ userlist . }}{{end}}{{with .Teams}}
// Teams: {{ teamlist . }}{{end}}
// Created: {{ .CreatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// Updated: {{ .UpdatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// `

// // tmplModUserList represents a row within mod user listing.
// var tmplModUserList = "Slug: \x1b[33m{{ .User.Slug }}\x1b[0m" + `
// ID: {{ .User.ID }}
// Username: {{ .User.Username }}
// Permission: {{ .Perm }}
// `

// // tmplModTeamList represents a row within mod team listing.
// var tmplModTeamList = "Slug: \x1b[33m{{ .Team.Slug }}\x1b[0m" + `
// ID: {{ .Team.ID }}
// Name: {{ .Team.Name }}
// Permission: {{ .Perm }}
// `

// // Mod provides the sub-command for the mod API.
// func Mod() *cli.Command {
// 	return &cli.Command{
// 		Name:  "mod",
// 		Usage: "Mod related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "List all mods",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplModList,
// 						Usage: "Custom output format",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "json",
// 						Value: false,
// 						Usage: "Print in JSON format",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "xml",
// 						Value: false,
// 						Usage: "Print in XML format",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ModList)
// 				},
// 			},
// 			{
// 				Name:      "show",
// 				Usage:     "Display a mod",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Mod ID or slug to show",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplModShow,
// 						Usage: "Custom output format",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "json",
// 						Value: false,
// 						Usage: "Print in JSON format",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "xml",
// 						Value: false,
// 						Usage: "Print in XML format",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ModShow)
// 				},
// 			},
// 			{
// 				Name:      "delete",
// 				Aliases:   []string{"rm"},
// 				Usage:     "Delete a mod",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Mod ID or slug to delete",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ModDelete)
// 				},
// 			},
// 			{
// 				Name:      "update",
// 				Usage:     "Update a mod",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Mod ID or slug to update",
// 					},
// 					&cli.StringFlag{
// 						Name:  "slug",
// 						Value: "",
// 						Usage: "Provide a slug",
// 					},
// 					&cli.StringFlag{
// 						Name:  "name",
// 						Value: "",
// 						Usage: "Provide a name",
// 					},
// 					&cli.StringFlag{
// 						Name:  "side",
// 						Value: "both",
// 						Usage: "Provide a side",
// 					},
// 					&cli.StringFlag{
// 						Name:  "description",
// 						Value: "",
// 						Usage: "Provide a description",
// 					},
// 					&cli.StringFlag{
// 						Name:  "author",
// 						Value: "",
// 						Usage: "Provide an author",
// 					},
// 					&cli.StringFlag{
// 						Name:  "website-link",
// 						Value: "",
// 						Usage: "Provide a website link",
// 					},
// 					&cli.StringFlag{
// 						Name:  "donate-link",
// 						Value: "",
// 						Usage: "Provide a donate link",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ModUpdate)
// 				},
// 			},
// 			{
// 				Name:      "create",
// 				Usage:     "Create a mod",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "slug",
// 						Value: "",
// 						Usage: "Provide a slug",
// 					},
// 					&cli.StringFlag{
// 						Name:  "name",
// 						Value: "",
// 						Usage: "Provide a name",
// 					},
// 					&cli.StringFlag{
// 						Name:  "side",
// 						Value: "both",
// 						Usage: "Provide a side",
// 					},
// 					&cli.StringFlag{
// 						Name:  "description",
// 						Value: "",
// 						Usage: "Provide a description",
// 					},
// 					&cli.StringFlag{
// 						Name:  "author",
// 						Value: "",
// 						Usage: "Provide an author",
// 					},
// 					&cli.StringFlag{
// 						Name:  "website-link",
// 						Value: "",
// 						Usage: "Provide a website link",
// 					},
// 					&cli.StringFlag{
// 						Name:  "donate-link",
// 						Value: "",
// 						Usage: "Provide a donate link",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ModCreate)
// 				},
// 			},
// 			{
// 				Name:  "user",
// 				Usage: "User assignments",
// 				Subcommands: []*cli.Command{
// 					{
// 						Name:      "list",
// 						Aliases:   []string{"ls"},
// 						Usage:     "List assigned users",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to list users",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplModUserList,
// 								Usage: "Custom output format",
// 							},
// 							&cli.BoolFlag{
// 								Name:  "json",
// 								Value: false,
// 								Usage: "Print in JSON format",
// 							},
// 							&cli.BoolFlag{
// 								Name:  "xml",
// 								Value: false,
// 								Usage: "Print in XML format",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModUserList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "Append a user to mod",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "user, u",
// 								Value: "",
// 								Usage: "User ID or slug to append",
// 							},
// 							&cli.StringFlag{
// 								Name:  "perm",
// 								Value: "user",
// 								Usage: "Permission for the user, can be user, admin or owner",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModUserAppend)
// 						},
// 					},
// 					{
// 						Name:      "perm",
// 						Usage:     "Update mod user permissions",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to update",
// 							},
// 							&cli.StringFlag{
// 								Name:  "user, u",
// 								Value: "",
// 								Usage: "User ID or slug to update",
// 							},
// 							&cli.StringFlag{
// 								Name:  "perm",
// 								Value: "user",
// 								Usage: "Permission for the user, can be user, admin or owner",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModUserPerm)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "Remove a user from mod",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "user, u",
// 								Value: "",
// 								Usage: "User ID or slug to remove",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModUserRemove)
// 						},
// 					},
// 				},
// 			},
// 			{
// 				Name:  "team",
// 				Usage: "Team assignments",
// 				Subcommands: []*cli.Command{
// 					{
// 						Name:      "list",
// 						Aliases:   []string{"ls"},
// 						Usage:     "List assigned teams",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to list teams",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplModTeamList,
// 								Usage: "Custom output format",
// 							},
// 							&cli.BoolFlag{
// 								Name:  "json",
// 								Value: false,
// 								Usage: "Print in JSON format",
// 							},
// 							&cli.BoolFlag{
// 								Name:  "xml",
// 								Value: false,
// 								Usage: "Print in XML format",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModTeamList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "Append a team to mod",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "team, t",
// 								Value: "",
// 								Usage: "Team ID or slug to append",
// 							},
// 							&cli.StringFlag{
// 								Name:  "perm",
// 								Value: "user",
// 								Usage: "Permission for the team, can be user, admin or owner",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModTeamAppend)
// 						},
// 					},
// 					{
// 						Name:      "perm",
// 						Usage:     "Update mod team permissions",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to update",
// 							},
// 							&cli.StringFlag{
// 								Name:  "team, t",
// 								Value: "",
// 								Usage: "Team ID or slug to update",
// 							},
// 							&cli.StringFlag{
// 								Name:  "perm",
// 								Value: "user",
// 								Usage: "Permission for the team, can be user, admin or owner",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModTeamPerm)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "Remove a team from mod",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Mod ID or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "team, t",
// 								Value: "",
// 								Usage: "Team ID or slug to remove",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ModTeamRemove)
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // ModList provides the sub-command to list all mods.
// func ModList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.ModList()

// 	if err != nil {
// 		return err
// 	}

// 	if c.IsSet("json") && c.IsSet("xml") {
// 		return fmt.Errorf("conflict, you can only use json or xml at once")
// 	}

// 	if c.Bool("xml") {
// 		res, err := xml.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if c.Bool("json") {
// 		res, err := json.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if len(records) == 0 {
// 		fmt.Fprintf(os.Stderr, "Empty result\n")
// 		return nil
// 	}

// 	tmpl, err := template.New(
// 		"_",
// 	).Funcs(
// 		globalFuncMap,
// 	).Funcs(
// 		sprigFuncMap,
// 	).Parse(
// 		fmt.Sprintf("%s\n", c.String("format")),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	for _, record := range records {
// 		err := tmpl.Execute(os.Stdout, record)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // ModShow provides the sub-command to show mod details.
// func ModShow(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.ModGet(
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if c.IsSet("json") && c.IsSet("xml") {
// 		return fmt.Errorf("conflict, you can only use json or xml at once")
// 	}

// 	if c.Bool("xml") {
// 		res, err := xml.MarshalIndent(record, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if c.Bool("json") {
// 		res, err := json.MarshalIndent(record, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	tmpl, err := template.New(
// 		"_",
// 	).Funcs(
// 		globalFuncMap,
// 	).Funcs(
// 		sprigFuncMap,
// 	).Parse(
// 		fmt.Sprintf("%s\n", c.String("format")),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return tmpl.Execute(os.Stdout, record)
// }

// // ModDelete provides the sub-command to delete a mod.
// func ModDelete(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModDelete(
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully delete\n")
// 	return nil
// }

// // ModUpdate provides the sub-command to update a mod.
// func ModUpdate(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.ModGet(
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	changed := false

// 	if val := c.String("name"); c.IsSet("name") && val != record.Name {
// 		record.Name = val
// 		changed = true
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != record.Slug {
// 		record.Slug = val
// 		changed = true
// 	}

// 	if val := c.String("side"); c.IsSet("side") && val != record.Side {
// 		record.Side = val
// 		changed = true
// 	}

// 	if val := c.String("description"); c.IsSet("description") && val != record.Description {
// 		record.Description = val
// 		changed = true
// 	}

// 	if val := c.String("author"); c.IsSet("author") && val != record.Author {
// 		record.Author = val
// 		changed = true
// 	}

// 	if val := c.String("website"); c.IsSet("website") && val != record.Website {
// 		record.Website = val
// 		changed = true
// 	}

// 	if val := c.String("donate"); c.IsSet("donate") && val != record.Donate {
// 		record.Donate = val
// 		changed = true
// 	}

// 	if changed {
// 		_, patch := client.ModPatch(
// 			record,
// 		)

// 		if patch != nil {
// 			return patch
// 		}

// 		fmt.Fprintf(os.Stderr, "Successfully updated\n")
// 	} else {
// 		fmt.Fprintf(os.Stderr, "Nothing to update...\n")
// 	}

// 	return nil
// }

// // ModCreate provides the sub-command to create a mod.
// func ModCreate(c *cli.Context, client kleister.ClientAPI) error {
// 	record := &kleister.Mod{}

// 	if val := c.String("name"); c.IsSet("name") && val != "" {
// 		record.Name = val
// 	} else {
// 		return fmt.Errorf("you must provide a name")
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != "" {
// 		record.Slug = val
// 	}

// 	if val := c.String("side"); c.IsSet("side") && val != "" {
// 		record.Side = val
// 	}

// 	if val := c.String("description"); c.IsSet("description") && val != "" {
// 		record.Description = val
// 	}

// 	if val := c.String("author"); c.IsSet("author") && val != "" {
// 		record.Author = val
// 	}

// 	if val := c.String("website"); c.IsSet("website") && val != "" {
// 		record.Website = val
// 	}

// 	if val := c.String("donate"); c.IsSet("donate") && val != "" {
// 		record.Donate = val
// 	}

// 	_, err := client.ModPost(
// 		record,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully created\n")
// 	return nil
// }

// // ModUserList provides the sub-command to list users of the mod.
// func ModUserList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.ModUserList(
// 		kleister.ModUserParams{
// 			Mod: GetIdentifierParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if c.IsSet("json") && c.IsSet("xml") {
// 		return fmt.Errorf("conflict, you can only use json or xml at once")
// 	}

// 	if c.Bool("xml") {
// 		res, err := xml.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if c.Bool("json") {
// 		res, err := json.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if len(records) == 0 {
// 		fmt.Fprintf(os.Stderr, "Empty result\n")
// 		return nil
// 	}

// 	tmpl, err := template.New(
// 		"_",
// 	).Funcs(
// 		globalFuncMap,
// 	).Funcs(
// 		sprigFuncMap,
// 	).Parse(
// 		fmt.Sprintf("%s\n", c.String("format")),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	for _, record := range records {
// 		err := tmpl.Execute(os.Stdout, record)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // ModUserAppend provides the sub-command to append a user to the mod.
// func ModUserAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModUserAppend(
// 		kleister.ModUserParams{
// 			Mod:  GetIdentifierParam(c),
// 			User: GetUserParam(c),
// 			Perm: GetPermParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully appended to mod\n")
// 	return nil
// }

// // ModUserPerm provides the sub-command to update mod user permissions.
// func ModUserPerm(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModUserPerm(
// 		kleister.ModUserParams{
// 			Mod:  GetIdentifierParam(c),
// 			User: GetUserParam(c),
// 			Perm: GetPermParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully updated permissions\n")
// 	return nil
// }

// // ModUserRemove provides the sub-command to remove a user from the mod.
// func ModUserRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModUserDelete(
// 		kleister.ModUserParams{
// 			Mod:  GetIdentifierParam(c),
// 			User: GetUserParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully removed from mod\n")
// 	return nil
// }

// // ModTeamList provides the sub-command to list teams of the mod.
// func ModTeamList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.ModTeamList(
// 		kleister.ModTeamParams{
// 			Mod: GetIdentifierParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if c.IsSet("json") && c.IsSet("xml") {
// 		return fmt.Errorf("conflict, you can only use json or xml at once")
// 	}

// 	if c.Bool("xml") {
// 		res, err := xml.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if c.Bool("json") {
// 		res, err := json.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if len(records) == 0 {
// 		fmt.Fprintf(os.Stderr, "Empty result\n")
// 		return nil
// 	}

// 	tmpl, err := template.New(
// 		"_",
// 	).Funcs(
// 		globalFuncMap,
// 	).Funcs(
// 		sprigFuncMap,
// 	).Parse(
// 		fmt.Sprintf("%s\n", c.String("format")),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	for _, record := range records {
// 		err := tmpl.Execute(os.Stdout, record)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // ModTeamAppend provides the sub-command to append a team to the mod.
// func ModTeamAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModTeamAppend(
// 		kleister.ModTeamParams{
// 			Mod:  GetIdentifierParam(c),
// 			Team: GetTeamParam(c),
// 			Perm: GetPermParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully appended to team\n")
// 	return nil
// }

// // ModTeamPerm provides the sub-command to update mod team permissions.
// func ModTeamPerm(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModTeamPerm(
// 		kleister.ModTeamParams{
// 			Mod:  GetIdentifierParam(c),
// 			Team: GetTeamParam(c),
// 			Perm: GetPermParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully updated permissions\n")
// 	return nil
// }

// // ModTeamRemove provides the sub-command to remove a team from the mod.
// func ModTeamRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ModTeamDelete(
// 		kleister.ModTeamParams{
// 			Mod:  GetIdentifierParam(c),
// 			Team: GetTeamParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully removed from team\n")
// 	return nil
// }
