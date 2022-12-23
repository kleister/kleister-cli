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

// // tmplKeyList represents a row within key listing.
// var tmplKeyList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// `

// // tmplKeyShow represents a key within details view.
// var tmplKeyShow = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// Key: {{ .Value }}
// Created: {{ .CreatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// Updated: {{ .UpdatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// `

// // Key provides the sub-command for the key API.
// func Key() *cli.Command {
// 	return &cli.Command{
// 		Name:  "key",
// 		Usage: "Key related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "List all keys",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplKeyList,
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
// 					return Handle(c, KeyList)
// 				},
// 			},
// 			{
// 				Name:      "show",
// 				Usage:     "Display a key",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Key ID or slug to show",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplKeyShow,
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
// 					return Handle(c, KeyShow)
// 				},
// 			},
// 			{
// 				Name:      "delete",
// 				Aliases:   []string{"rm"},
// 				Usage:     "Delete a key",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Key ID or slug to delete",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, KeyDelete)
// 				},
// 			},
// 			{
// 				Name:      "update",
// 				Usage:     "Update a key",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Key ID or slug to update",
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
// 						Name:  "key",
// 						Value: "",
// 						Usage: "Provide a key",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, KeyUpdate)
// 				},
// 			},
// 			{
// 				Name:      "create",
// 				Usage:     "Create a key",
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
// 						Name:  "key",
// 						Value: "",
// 						Usage: "Provide a key",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, KeyCreate)
// 				},
// 			},
// 		},
// 	}
// }

// // KeyList provides the sub-command to list all keys.
// func KeyList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.KeyList()

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

// // KeyShow provides the sub-command to show key details.
// func KeyShow(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.KeyGet(
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

// // KeyDelete provides the sub-command to delete a key.
// func KeyDelete(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.KeyDelete(
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully delete\n")
// 	return nil
// }

// // KeyUpdate provides the sub-command to update a key.
// func KeyUpdate(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.KeyGet(
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

// 	if val := c.String("key"); c.IsSet("key") && val != record.Value {
// 		record.Value = val
// 		changed = true
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != record.Slug {
// 		record.Slug = val
// 		changed = true
// 	}

// 	if changed {
// 		_, patch := client.KeyPatch(
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

// // KeyCreate provides the sub-command to create a key.
// func KeyCreate(c *cli.Context, client kleister.ClientAPI) error {
// 	record := &kleister.Key{}

// 	if val := c.String("name"); c.IsSet("name") && val != "" {
// 		record.Name = val
// 	} else {
// 		return fmt.Errorf("you must provide a name")
// 	}

// 	if val := c.String("key"); c.IsSet("key") && val != "" {
// 		record.Value = val
// 	} else {
// 		return fmt.Errorf("you must provide a key")
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != "" {
// 		record.Slug = val
// 	}

// 	_, err := client.KeyPost(
// 		record,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully created\n")
// 	return nil
// }
