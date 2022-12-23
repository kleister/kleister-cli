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
// 	"gopkg.in/guregu/null.v3"
// 	"gopkg.in/urfave/cli.v2"
// )

// // tmplBuildList represents a row within build listing.
// var tmplBuildList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// `

// // tmplBuildShow represents a build within details view.
// var tmplBuildShow = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}{{ with .Pack }}
// Pack: {{ .Name }}{{ end }}{{ with .Minecraft }}
// Minecraft: {{ . }}{{ end }}{{ with .Forge }}
// Forge: {{ . }}{{ end }}{{ with .MinJava }}
// Java: {{ . }}{{ end }}{{ with .MinMemory }}
// Memory: {{ . }}{{ end }}
// Published: {{ .Published }}
// Private: {{ .Private }}{{ with .Versions }}
// Versions: {{ versionlist . }}{{ end }}
// Created: {{ .CreatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// Updated: {{ .UpdatedAt.Format "Mon Jan _2 15:04:05 MST 2006" }}
// `

// // tmplBuildVersionList represents a row within build version listing.
// var tmplBuildVersionList = "Slug: \x1b[33m{{ .Version.Slug }}\x1b[0m" + `
// ID: {{ .Version.ID }}
// Name: {{ .Version.Name }}
// Mod: {{ with .Version.Mod }}{{ . }}{{ else }}n/a{{ end }}
// `

// // Build provides the sub-command for the build API.
// func Build() *cli.Command {
// 	return &cli.Command{
// 		Name:  "build",
// 		Usage: "Build related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "list all builds",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "pack, p",
// 						Value: "",
// 						Usage: "id or slug of the related pack",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplBuildList,
// 						Usage: "custom output format",
// 					},
// 					&cli.StringFlag{
// 						Name:  "output",
// 						Value: "text",
// 						Usage: "output as format, json or xml",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, BuildList)
// 				},
// 			},
// 			{
// 				Name:      "show",
// 				Usage:     "display a build",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "pack, p",
// 						Value: "",
// 						Usage: "id or slug of the related pack",
// 					},
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "build id or slug to show",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplBuildShow,
// 						Usage: "custom output format",
// 					},
// 					&cli.StringFlag{
// 						Name:  "output",
// 						Value: "text",
// 						Usage: "output as format, json or xml",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, BuildShow)
// 				},
// 			},
// 			{
// 				Name:      "delete",
// 				Aliases:   []string{"rm"},
// 				Usage:     "delete a build",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "pack, p",
// 						Value: "",
// 						Usage: "id or slug of the related pack",
// 					},
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "build id or slug to delete",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, BuildDelete)
// 				},
// 			},
// 			{
// 				Name:      "update",
// 				Usage:     "update a build",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "pack, p",
// 						Value: "",
// 						Usage: "id or slug of the related pack",
// 					},
// 					&cli.StringFlag{
// 						Name:  "id, i",
// 						Value: "",
// 						Usage: "build id or slug to update",
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
// 						Name:  "min-java",
// 						Value: "",
// 						Usage: "minimal Java version",
// 					},
// 					&cli.StringFlag{
// 						Name:  "min-memory",
// 						Value: "",
// 						Usage: "minimal memory alloc",
// 					},
// 					&cli.StringFlag{
// 						Name:  "minecraft",
// 						Value: "",
// 						Usage: "provide a minecraft id or slug",
// 					},
// 					&cli.StringFlag{
// 						Name:  "forge",
// 						Value: "",
// 						Usage: "provide a forge id or slug",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "published",
// 						Value: false,
// 						Usage: "mark build published",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "hidden",
// 						Value: false,
// 						Usage: "mark pack hidden",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "private",
// 						Value: false,
// 						Usage: "mark build private",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "public",
// 						Value: false,
// 						Usage: "mark pack public",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, BuildUpdate)
// 				},
// 			},
// 			{
// 				Name:      "create",
// 				Usage:     "create a build",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "pack, p",
// 						Value: "",
// 						Usage: "id or slug of the related pack",
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
// 						Name:  "min-java",
// 						Value: "",
// 						Usage: "minimal Java version",
// 					},
// 					&cli.StringFlag{
// 						Name:  "min-memory",
// 						Value: "",
// 						Usage: "minimal memory alloc",
// 					},
// 					&cli.StringFlag{
// 						Name:  "minecraft",
// 						Value: "",
// 						Usage: "provide a minecraft id or slug",
// 					},
// 					&cli.StringFlag{
// 						Name:  "forge",
// 						Value: "",
// 						Usage: "provide a forge id or slug",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "published",
// 						Value: false,
// 						Usage: "mark build published",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "hidden",
// 						Value: false,
// 						Usage: "mark pack hidden",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "private",
// 						Value: false,
// 						Usage: "mark build private",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "public",
// 						Value: false,
// 						Usage: "mark pack public",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, BuildCreate)
// 				},
// 			},
// 			{
// 				Name:  "version",
// 				Usage: "version assignments",
// 				Subcommands: []*cli.Command{
// 					{
// 						Name:      "list",
// 						Aliases:   []string{"ls"},
// 						Usage:     "list assigned versions",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "id or slug of the related pack",
// 							},
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "build id or slug to list versions",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplBuildVersionList,
// 								Usage: "custom output format",
// 							},
// 							&cli.StringFlag{
// 								Name:  "output",
// 								Value: "text",
// 								Usage: "output as format, json or xml",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, BuildVersionList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "append a version to build",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "id or slug of the related pack",
// 							},
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "build id or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "mod, m",
// 								Value: "",
// 								Usage: "mod id or slug to append",
// 							},
// 							&cli.StringFlag{
// 								Name:  "version, V",
// 								Value: "",
// 								Usage: "version id or slug to append",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, BuildVersionAppend)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "remove a version from build",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "id or slug of the related pack",
// 							},
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "build id or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "mod, m",
// 								Value: "",
// 								Usage: "mod id or slug to remove",
// 							},
// 							&cli.StringFlag{
// 								Name:  "version, V",
// 								Value: "",
// 								Usage: "version id or slug to remove",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, BuildVersionRemove)
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // BuildList provides the sub-command to list all builds.
// func BuildList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.BuildList(
// 		GetPackParam(c),
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

// // BuildShow provides the sub-command to show build details.
// func BuildShow(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.BuildGet(
// 		GetPackParam(c),
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

// // BuildDelete provides the sub-command to delete a build.
// func BuildDelete(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.BuildDelete(
// 		GetPackParam(c),
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully deleted\n")
// 	return nil
// }

// // BuildUpdate provides the sub-command to update a build.
// func BuildUpdate(c *cli.Context, client kleister.ClientAPI) error {
// 	record, err := client.BuildGet(
// 		GetPackParam(c),
// 		GetIdentifierParam(c),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	changed := false

// 	if c.IsSet("minecraft") {
// 		if match, _ := regexp.MatchString("^([0-9]+)$", c.String("minecraft")); match {
// 			if val, err := strconv.ParseInt(c.String("minecraft"), 10, 64); err == nil && val != record.MinecraftID.Int64 {
// 				record.MinecraftID = null.NewInt(val, val > 0)
// 				changed = true
// 			}
// 		} else {
// 			if c.String("minecraft") != "" {
// 				related, err := client.MinecraftGet(
// 					c.String("minecraft"),
// 				)

// 				if err != nil {
// 					return err
// 				}

// 				if related.ID != record.MinecraftID.Int64 {
// 					record.MinecraftID = null.NewInt(related.ID, related.ID > 0)
// 					changed = true
// 				}
// 			}
// 		}
// 	}

// 	if c.IsSet("forge") {
// 		if match, _ := regexp.MatchString("^([0-9]+)$", c.String("forge")); match {
// 			if val, err := strconv.ParseInt(c.String("forge"), 10, 64); err == nil && val != record.ForgeID.Int64 {
// 				record.ForgeID = null.NewInt(val, val > 0)
// 				changed = true
// 			}
// 		} else {
// 			if c.String("forge") != "" {
// 				related, err := client.ForgeGet(
// 					c.String("forge"),
// 				)

// 				if err != nil {
// 					return err
// 				}

// 				if related.ID != record.ForgeID.Int64 {
// 					record.ForgeID = null.NewInt(related.ID, related.ID > 0)
// 					changed = true
// 				}
// 			}
// 		}
// 	}

// 	if val := c.String("name"); c.IsSet("name") && val != record.Name {
// 		record.Name = val
// 		changed = true
// 	}

// 	if val := c.String("slug"); c.IsSet("slug") && val != record.Slug {
// 		record.Slug = val
// 		changed = true
// 	}

// 	if val := c.String("min-java"); c.IsSet("min-java") && val != record.MinJava {
// 		record.MinJava = val
// 		changed = true
// 	}

// 	if val := c.String("min-memory"); c.IsSet("min-memory") && val != record.MinMemory {
// 		record.MinMemory = val
// 		changed = true
// 	}

// 	if c.IsSet("published") && c.IsSet("hidden") {
// 		return fmt.Errorf("conflict, you can mark it only published or hidden")
// 	}

// 	if c.IsSet("published") {
// 		record.Published = true
// 		changed = true
// 	}

// 	if c.IsSet("hidden") {
// 		record.Published = false
// 		changed = true
// 	}

// 	if c.IsSet("private") && c.IsSet("public") {
// 		return fmt.Errorf("conflict, you can mark it only private or public")
// 	}

// 	if c.IsSet("private") {
// 		record.Private = true
// 		changed = true
// 	}

// 	if c.IsSet("public") {
// 		record.Private = false
// 		changed = true
// 	}

// 	if changed {
// 		_, patch := client.BuildPatch(
// 			GetPackParam(c),
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

// // BuildCreate provides the sub-command to create a build.
// func BuildCreate(c *cli.Context, client kleister.ClientAPI) error {
// 	record := &kleister.Build{}

// 	if c.String("pack") == "" {
// 		return fmt.Errorf("you must provide a pack id or slug")
// 	}

// 	if c.IsSet("pack") {
// 		if match, _ := regexp.MatchString("^([0-9]+)$", c.String("pack")); match {
// 			if val, err := strconv.ParseInt(c.String("pack"), 10, 64); err == nil && val != 0 {
// 				record.PackID = val
// 			}
// 		} else {
// 			if c.String("pack") != "" {
// 				related, err := client.PackGet(
// 					c.String("pack"),
// 				)

// 				if err != nil {
// 					return err
// 				}

// 				if related.ID != record.PackID {
// 					record.PackID = related.ID
// 				}
// 			}
// 		}
// 	}

// 	if c.IsSet("minecraft") {
// 		if match, _ := regexp.MatchString("^([0-9]+)$", c.String("minecraft")); match {
// 			if val, err := strconv.ParseInt(c.String("minecraft"), 10, 64); err == nil && val != 0 {
// 				record.MinecraftID = null.NewInt(val, val > 0)
// 			}
// 		} else {
// 			if c.String("minecraft") != "" {
// 				related, err := client.MinecraftGet(
// 					c.String("minecraft"),
// 				)

// 				if err != nil {
// 					return err
// 				}

// 				if related.ID != record.MinecraftID.Int64 {
// 					record.MinecraftID = null.NewInt(related.ID, related.ID > 0)
// 				}
// 			}
// 		}
// 	}

// 	if c.IsSet("forge") {
// 		if match, _ := regexp.MatchString("^([0-9]+)$", c.String("forge")); match {
// 			if val, err := strconv.ParseInt(c.String("forge"), 10, 64); err == nil && val != 0 {
// 				record.ForgeID = null.NewInt(val, val > 0)
// 			}
// 		} else {
// 			if c.String("forge") != "" {
// 				related, err := client.ForgeGet(
// 					c.String("forge"),
// 				)

// 				if err != nil {
// 					return err
// 				}

// 				if related.ID != record.ForgeID.Int64 {
// 					record.ForgeID = null.NewInt(related.ID, related.ID > 0)
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

// 	if val := c.String("min-java"); c.IsSet("min-java") && val != "" {
// 		record.MinJava = val
// 	}

// 	if val := c.String("min-memory"); c.IsSet("min-memory") && val != "" {
// 		record.MinMemory = val
// 	}

// 	if c.IsSet("published") && c.IsSet("hidden") {
// 		return fmt.Errorf("conflict, you can mark it only published or hidden")
// 	}

// 	if c.IsSet("published") {
// 		record.Published = true
// 	}

// 	if c.IsSet("hidden") {
// 		record.Published = false
// 	}

// 	if c.IsSet("private") && c.IsSet("public") {
// 		return fmt.Errorf("conflict, you can mark it only private or public")
// 	}

// 	if c.IsSet("private") {
// 		record.Private = true
// 	}

// 	if c.IsSet("public") {
// 		record.Private = false
// 	}

// 	_, err := client.BuildPost(
// 		GetPackParam(c),
// 		record,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully created\n")
// 	return nil
// }

// // BuildVersionList provides the sub-command to list versions of the build.
// func BuildVersionList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.BuildVersionList(
// 		kleister.BuildVersionParams{
// 			Pack:  GetPackParam(c),
// 			Build: GetIdentifierParam(c),
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

// // BuildVersionAppend provides the sub-command to append a version to the build.
// func BuildVersionAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.BuildVersionAppend(
// 		kleister.BuildVersionParams{
// 			Pack:    GetPackParam(c),
// 			Build:   GetIdentifierParam(c),
// 			Mod:     GetModParam(c),
// 			Version: GetVersionParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully appended to build\n")
// 	return nil
// }

// // BuildVersionRemove provides the sub-command to remove a version from the build.
// func BuildVersionRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.BuildVersionDelete(
// 		kleister.BuildVersionParams{
// 			Pack:    GetPackParam(c),
// 			Build:   GetIdentifierParam(c),
// 			Mod:     GetModParam(c),
// 			Version: GetVersionParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "successfully removed from build\n")
// 	return nil
// }
