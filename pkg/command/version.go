package command

// import (
// 	"encoding/json"
// 	"encoding/xml"
// 	"fmt"
// 	"os"
// 	"regexp"
// 	"strconv"
// 	"text/template"

// 	"github.com/kleister/kleister-go/kleister"
// 	"gopkg.in/urfave/cli.v2"
// )

// // tmplVersionList represents a row within forge listing.
// var tmplVersionList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// `

// // tmplVersionShow represents a version within details view.
// var tmplVersionShow = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}{{with .Mod}}
// Mod: {{ . }}{{end}}{{with .File}}
// File: {{ . }}{{end}}{{with .Builds}}
// Builds: {{ buildlist . }}{{end}}
// Created: {{ .CreatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// Updated: {{ .UpdatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// `

// // tmplVersionBuildList represents a row within version build listing.
// var tmplVersionBuildList = "Slug: \x1b[33m{{ .Build.Slug }}\x1b[0m" + `
// ID: {{ .Build.ID }}
// Name: {{ .Build.Name }}
// Pack: {{with .Build.Pack}}{{ . }}{{else}}n/a{{end}}
// `

// // Version provides the sub-command for the version API.
// func Version() *cli.Command {
// 	return &cli.Command{
// 		Name:  "version",
// 		Usage: "Version related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "List all versions",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "mod, m",
// 						Value: "",
// 						Usage: "ID or slug of the related mod",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplVersionList,
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
// 					return Handle(c, VersionList)
// 				},
// 			},
// 			{
// 				Name:      "show",
// 				Usage:     "Display a version",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "mod, m",
// 						Value: "",
// 						Usage: "ID or slug of the related mod",
// 					},
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Version ID or slug to show",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplVersionShow,
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
// 					return Handle(c, VersionShow)
// 				},
// 			},
// 			{
// 				Name:      "delete",
// 				Aliases:   []string{"rm"},
// 				Usage:     "Delete a version",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "mod, m",
// 						Value: "",
// 						Usage: "ID or slug of the related mod",
// 					},
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Version ID or slug to show",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, VersionDelete)
// 				},
// 			},
// 			{
// 				Name:      "update",
// 				Usage:     "Update a version",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "mod, m",
// 						Value: "",
// 						Usage: "ID or slug of the related mod",
// 					},
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "Version ID or slug to show",
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
// 						Name:  "file-url",
// 						Value: "",
// 						Usage: "Provide a file URL",
// 					},
// 					&cli.StringFlag{
// 						Name:  "file-path",
// 						Value: "",
// 						Usage: "Provide a file path",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, VersionUpdate)
// 				},
// 			},
// 			{
// 				Name:      "create",
// 				Usage:     "Create a version",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "mod, m",
// 						Value: "",
// 						Usage: "ID or slug of the related mod",
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
// 						Name:  "file-url",
// 						Value: "",
// 						Usage: "Provide a file URL",
// 					},
// 					&cli.StringFlag{
// 						Name:  "file-path",
// 						Value: "",
// 						Usage: "Provide a file path",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, VersionCreate)
// 				},
// 			},
// 			{
// 				Name:  "build",
// 				Usage: "Build assignments",
// 				Subcommands: []*cli.Command{
// 					{
// 						Name:      "list",
// 						Aliases:   []string{"ls"},
// 						Usage:     "List assigned builds",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "mod, m",
// 								Value: "",
// 								Usage: "ID or slug of the related mod",
// 							},
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Version ID or slug to list builds",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplVersionBuildList,
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
// 							return Handle(c, VersionBuildList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "Append a build to version",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "mod, m",
// 								Value: "",
// 								Usage: "ID or slug of the related mod",
// 							},
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Version ID or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "Pack ID or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "build, b",
// 								Value: "",
// 								Usage: "Build ID or slug to append to",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, VersionBuildAppend)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "Remove a build from version",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "mod, m",
// 								Value: "",
// 								Usage: "ID or slug of the related mod",
// 							},
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Version ID or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "Pack ID or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "build, b",
// 								Value: "",
// 								Usage: "Build ID or slug to append to",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, VersionBuildRemove)
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // VersionList provides the sub-command to list all versions.
// func VersionList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.VersionList(
// 		GetModParam(c),
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

// // VersionShow provides the sub-command to show version details.
// func VersionShow(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.VersionGet(
// 		GetModParam(c),
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

// // VersionDelete provides the sub-command to delete a version.
// func VersionDelete(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.VersionDelete(
// 		GetModParam(c),
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully delete\n")
// 	return nil
// }

// // VersionUpdate provides the sub-command to update a version.
// func VersionUpdate(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.VersionGet(
// 		GetModParam(c),
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

// 	if val := c.String("file-url"); c.IsSet("file-url") && val != "" {
// 		err := record.DownloadFile(
// 			val,
// 		)

// 		if err != nil {
// 			return fmt.Errorf("failed to download and encode file")
// 		}

// 		changed = true
// 	}

// 	if val := c.String("file-path"); c.IsSet("file-path") && val != "" {
// 		err := record.EncodeFile(
// 			val,
// 		)

// 		if err != nil {
// 			return fmt.Errorf("failed to encode file")
// 		}

// 		changed = true
// 	}

// 	if changed {
// 		_, patch := client.VersionPatch(
// 			GetModParam(c),
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

// // VersionCreate provides the sub-command to create a version.
// func VersionCreate(c *cli.Context, client kleister.ClientAPI) error {
// 	record := &kleister.Version{}

// 	if c.String("mod") == "" {
// 		return fmt.Errorf("you must provide a mod id or slug")
// 	}

// 	if c.IsSet("mod") {
// 		if match, _ := regexp.MatchString("^([0-9]+)$", c.String("mod")); match {
// 			if val, err := strconv.ParseInt(c.String("mod"), 10, 64); err == nil && val != 0 {
// 				record.ModID = val
// 			}
// 		} else {
// 			if c.String("mod") != "" {
// 				related, err := client.ModGet(
// 					c.String("mod"),
// 				)

// 				if err != nil {
// 					return err
// 				}

// 				if related.ID != record.ModID {
// 					record.ModID = related.ID
// 				}
// 			}
// 		}
// 	}

// 	if val := c.String("name"); c.IsSet("name") && val != "" {
// 		record.Name = val
// 	} else {
// 		return fmt.Errorf("you must provide a name")
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != "" {
// 		record.Slug = val
// 	}

// 	if val := c.String("file-url"); c.IsSet("file-url") && val != "" {
// 		err := record.DownloadFile(
// 			val,
// 		)

// 		if err != nil {
// 			return fmt.Errorf("failed to download and encode file")
// 		}
// 	}

// 	if val := c.String("file-path"); c.IsSet("file-path") && val != "" {
// 		err := record.EncodeFile(
// 			val,
// 		)

// 		if err != nil {
// 			return fmt.Errorf("failed to encode file")
// 		}
// 	}

// 	_, err := client.VersionPost(
// 		GetModParam(c),
// 		record,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully created\n")
// 	return nil
// }

// // VersionBuildList provides the sub-command to list builds of the version.
// func VersionBuildList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.VersionBuildList(
// 		kleister.VersionBuildParams{
// 			Mod:     GetModParam(c),
// 			Version: GetIdentifierParam(c),
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

// // VersionBuildAppend provides the sub-command to append a build to the version.
// func VersionBuildAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.VersionBuildAppend(
// 		kleister.VersionBuildParams{
// 			Mod:     GetModParam(c),
// 			Version: GetIdentifierParam(c),
// 			Pack:    GetPackParam(c),
// 			Build:   GetBuildParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully appended to version\n")
// 	return nil
// }

// // VersionBuildRemove provides the sub-command to remove a build from the version.
// func VersionBuildRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.VersionBuildDelete(
// 		kleister.VersionBuildParams{
// 			Mod:     GetModParam(c),
// 			Version: GetIdentifierParam(c),
// 			Pack:    GetPackParam(c),
// 			Build:   GetBuildParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully removed from version\n")
// 	return nil
// }
