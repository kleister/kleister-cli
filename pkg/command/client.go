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

// // tmplClientList represents a row within client listing.
// var tmplClientList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// `

// // tmplClientShow represents a client within details view.
// var tmplClientShow = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// UUID: {{ .Value }}{{ with .Packs }}
// Packs: {{ packlist . }}{{ end }}
// Created: {{ .CreatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// Updated: {{ .UpdatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// `

// // tmplClientPackList represents a row within client pack listing.
// var tmplClientPackList = "Slug: \x1b[33m{{ .Pack.Slug }}\x1b[0m" + `
// ID: {{ .Pack.ID }}
// Name: {{ .Pack.Name }}
// `

// // Client provides the sub-command for the client API.
// func Client() *cli.Command {
// 	return &cli.Command{
// 		Name:  "client",
// 		Usage: "client related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "list all clients",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplClientList,
// 						Usage: "custom output format",
// 					},
// 					&cli.StringFlag{
// 						Name:  "output",
// 						Value: "text",
// 						Usage: "output as format, json or xml",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ClientList)
// 				},
// 			},
// 			{
// 				Name:      "show",
// 				Usage:     "display a client",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "client id or slug to show",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplClientShow,
// 						Usage: "custom output format",
// 					},
// 					&cli.StringFlag{
// 						Name:  "output",
// 						Value: "text",
// 						Usage: "output as format, json or xml",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ClientShow)
// 				},
// 			},
// 			{
// 				Name:      "delete",
// 				Aliases:   []string{"rm"},
// 				Usage:     "delete a client",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "client id or slug to delete",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ClientDelete)
// 				},
// 			},
// 			{
// 				Name:      "update",
// 				Usage:     "update a client",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "client id or slug to update",
// 					},
// 					&cli.StringFlag{
// 						Name:  "slug",
// 						Value: "",
// 						Usage: "provide a slug",
// 					},
// 					&cli.StringFlag{
// 						Name:  "name",
// 						Value: "",
// 						Usage: "provide a name",
// 					},
// 					&cli.StringFlag{
// 						Name:  "uuid",
// 						Value: "",
// 						Usage: "provide a uuid",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ClientUpdate)
// 				},
// 			},
// 			{
// 				Name:      "create",
// 				Usage:     "create a client",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "slug",
// 						Value: "",
// 						Usage: "provide a slug",
// 					},
// 					&cli.StringFlag{
// 						Name:  "name",
// 						Value: "",
// 						Usage: "provide a name",
// 					},
// 					&cli.StringFlag{
// 						Name:  "uuid",
// 						Value: "",
// 						Usage: "provide a uuid",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ClientCreate)
// 				},
// 			},
// 			{
// 				Name:  "pack",
// 				Usage: "pack assignments",
// 				Subcommands: []*cli.Command{
// 					{
// 						Name:      "list",
// 						Aliases:   []string{"ls"},
// 						Usage:     "list assigned packs",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "client id or slug to list packs",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplClientPackList,
// 								Usage: "custom output format",
// 							},
// 							&cli.StringFlag{
// 								Name:  "output",
// 								Value: "text",
// 								Usage: "output as format, json or xml",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ClientPackList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "append a pack to client",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "client id or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "pack id or slug to append",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ClientPackAppend)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "remove a pack from client",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "client id or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "pack id or slug to remove",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ClientPackRemove)
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // ClientList provides the sub-command to list all clients.
// func ClientList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.ClientList()

// 	if err != nil {
// 		return err
// 	}

// 	switch c.String("output") {
// 	case "json":
// 		res, err := json.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 	case "xml":
// 		res, err := xml.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 	case "text":
// 		if len(records) == 0 {
// 			fmt.Fprintf(os.Stderr, "empty result\n")
// 			return nil
// 		}

// 		tmpl, err := template.New(
// 			"_",
// 		).Funcs(
// 			globalFuncMap,
// 		).Funcs(
// 			sprigFuncMap,
// 		).Parse(
// 			fmt.Sprintf("%s\n", c.String("format")),
// 		)

// 		if err != nil {
// 			return err
// 		}

// 		for _, record := range records {
// 			err := tmpl.Execute(os.Stdout, record)

// 			if err != nil {
// 				return err
// 			}
// 		}
// 	default:
// 		return fmt.Errorf("invalid output type")
// 	}

// 	return nil
// }

// // ClientShow provides the sub-command to show client details.
// func ClientShow(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.ClientGet(
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	switch c.String("output") {
// 	case "json":
// 		res, err := json.MarshalIndent(record, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 	case "xml":
// 		res, err := xml.MarshalIndent(record, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 	case "text":
// 		tmpl, err := template.New(
// 			"_",
// 		).Funcs(
// 			globalFuncMap,
// 		).Funcs(
// 			sprigFuncMap,
// 		).Parse(
// 			fmt.Sprintf("%s\n", c.String("format")),
// 		)

// 		if err != nil {
// 			return err
// 		}

// 		return tmpl.Execute(os.Stdout, record)
// 	default:
// 		return fmt.Errorf("invalid output type")
// 	}

// 	return nil
// }

// // ClientDelete provides the sub-command to delete a client.
// func ClientDelete(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ClientDelete(
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully deleted\n")
// 	return nil
// }

// // ClientUpdate provides the sub-command to update a client.
// func ClientUpdate(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.ClientGet(
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

// 	if val := c.String("uuid"); c.IsSet("uuid") && val != record.Value {
// 		record.Value = val
// 		changed = true
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != record.Slug {
// 		record.Slug = val
// 		changed = true
// 	}

// 	if changed {
// 		_, patch := client.ClientPatch(
// 			record,
// 		)

// 		if patch != nil {
// 			return patch
// 		}

// 		fmt.Fprintf(os.Stderr, "successfully updated\n")
// 	} else {
// 		fmt.Fprintf(os.Stderr, "nothing to update!\n")
// 	}

// 	return nil
// }

// // ClientCreate provides the sub-command to create a client.
// func ClientCreate(c *cli.Context, client kleister.ClientAPI) error {
// 	record := &kleister.Client{}

// 	if val := c.String("name"); c.IsSet("name") && val != "" {
// 		record.Name = val
// 	} else {
// 		return fmt.Errorf("you must provide a name")
// 	}

// 	if val := c.String("uuid"); c.IsSet("uuid") && val != "" {
// 		record.Value = val
// 	} else {
// 		return fmt.Errorf("you must provide a uuid")
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != "" {
// 		record.Slug = val
// 	}

// 	_, err := client.ClientPost(
// 		record,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully created\n")
// 	return nil
// }

// // ClientPackList provides the sub-command to list packs of the client.
// func ClientPackList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.ClientPackList(
// 		kleister.ClientPackParams{
// 			Client: GetIdentifierParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	switch c.String("output") {
// 	case "json":
// 		res, err := json.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 	case "xml":
// 		res, err := xml.MarshalIndent(records, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 	case "text":
// 		if len(records) == 0 {
// 			fmt.Fprintf(os.Stderr, "empty result\n")
// 			return nil
// 		}

// 		tmpl, err := template.New(
// 			"_",
// 		).Funcs(
// 			globalFuncMap,
// 		).Funcs(
// 			sprigFuncMap,
// 		).Parse(
// 			fmt.Sprintf("%s\n", c.String("format")),
// 		)

// 		if err != nil {
// 			return err
// 		}

// 		for _, record := range records {
// 			err := tmpl.Execute(os.Stdout, record)

// 			if err != nil {
// 				return err
// 			}
// 		}
// 	default:
// 		return fmt.Errorf("invalid output type")
// 	}

// 	return nil
// }

// // ClientPackAppend provides the sub-command to append a pack to the client.
// func ClientPackAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ClientPackAppend(
// 		kleister.ClientPackParams{
// 			Client: GetIdentifierParam(c),
// 			Pack:   GetPackParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully appended to client\n")
// 	return nil
// }

// // ClientPackRemove provides the sub-command to remove a pack from the client.
// func ClientPackRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ClientPackDelete(
// 		kleister.ClientPackParams{
// 			Client: GetIdentifierParam(c),
// 			Pack:   GetPackParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully removed from client\n")
// 	return nil
// }
