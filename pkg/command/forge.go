package command

// import (
// 	"encoding/json"
// 	"encoding/xml"
// 	"fmt"
// 	"os"
// 	"sort"
// 	"strings"
// 	"text/template"

// 	"github.com/Knetic/govaluate"
// 	"github.com/kleister/kleister-go/kleister"
// 	"gopkg.in/urfave/cli.v2"
// )

// // tmplForgeList represents a row within forge listing.
// var tmplForgeList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Version: {{ .Version }}
// Minecraft: {{ .Minecraft }}
// `

// // tmplForgeBuildList represents a row within forge build listing.
// var tmplForgeBuildList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// Pack: {{ with .Pack }}{{ . }}{{ else }}n/a{{ end }}
// `

// // Forge provides the sub-command for the Forge API.
// func Forge() *cli.Command {
// 	return &cli.Command{
// 		Name:  "forge",
// 		Usage: "Forge related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "list all forge versions",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "sort",
// 						Value: "Slug",
// 						Usage: "sort by this field",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplForgeList,
// 						Usage: "custom output format",
// 					},
// 					&cli.StringFlag{
// 						Name:  "output",
// 						Value: "text",
// 						Usage: "output as format, json or xml",
// 					},
// 					&cli.StringFlag{
// 						Name:  "filter",
// 						Value: "",
// 						Usage: "filter by values",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "first",
// 						Value: false,
// 						Usage: "return only first record",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "last",
// 						Value: false,
// 						Usage: "return only last record",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ForgeList)
// 				},
// 			},
// 			{
// 				Name:      "refresh",
// 				Aliases:   []string{"ref"},
// 				Usage:     "refresh forge versions",
// 				ArgsUsage: " ",
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, ForgeRefresh)
// 				},
// 			},
// 			{
// 				Name:  "build",
// 				Usage: "build assignments",
// 				Subcommands: []*cli.Command{
// 					{
// 						Name:      "list",
// 						Aliases:   []string{"ls"},
// 						Usage:     "list assigned builds",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "forge id or slug to list builds",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplForgeBuildList,
// 								Usage: "custom output format",
// 							},
// 							&cli.StringFlag{
// 								Name:  "output",
// 								Value: "text",
// 								Usage: "output as format, json or xml",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ForgeBuildList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "append a build to forge",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "forge id or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "pack id or slug to append",
// 							},
// 							&cli.StringFlag{
// 								Name:  "build, b",
// 								Value: "",
// 								Usage: "build id or slug to append",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ForgeBuildAppend)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "remove a build from forge",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "forge id or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "pack id or slug to remove",
// 							},
// 							&cli.StringFlag{
// 								Name:  "build, b",
// 								Value: "",
// 								Usage: "build id or slug to remove",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, ForgeBuildRemove)
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // ForgeList provides the sub-command to list all Forge versions.
// func ForgeList(c *cli.Context, client kleister.ClientAPI) error {
// 	var (
// 		result []*kleister.Forge
// 	)

// 	records, err := client.ForgeList()

// 	if err != nil {
// 		return err
// 	}

// 	if c.IsSet("json") && c.IsSet("xml") {
// 		return fmt.Errorf("conflict, you can only use json or xml at once")
// 	}

// 	if c.IsSet("first") && c.IsSet("last") {
// 		return fmt.Errorf("conflict, you can only use first or last at once")
// 	}

// 	if c.IsSet("filter") {
// 		expression, err := govaluate.NewEvaluableExpression(
// 			c.String("filter"),
// 		)

// 		if err != nil {
// 			return fmt.Errorf("failed to parse filter. %s", err)
// 		}

// 		for _, record := range records {
// 			params := make(map[string]interface{}, 3)
// 			params["Slug"] = record.Slug
// 			params["Version"] = record.Version
// 			params["Minecraft"] = record.Minecraft

// 			match, err := expression.Evaluate(
// 				params,
// 			)

// 			if err != nil {
// 				return fmt.Errorf("failed to evaluate filter. %s", err)
// 			}

// 			if match.(bool) {
// 				result = append(result, record)
// 			}
// 		}
// 	} else {
// 		result = records
// 	}

// 	switch strings.ToLower(c.String("sort")) {
// 	case "slug":
// 		sort.Sort(
// 			kleister.ForgeBySlug(
// 				result,
// 			),
// 		)
// 	case "version":
// 		sort.Sort(
// 			kleister.ForgeByVersion(
// 				result,
// 			),
// 		)
// 	case "minecraft":
// 		sort.Sort(
// 			kleister.ForgeByMinecraft(
// 				result,
// 			),
// 		)
// 	default:
// 		return fmt.Errorf("the sort value %s is invalid, can be slug, version or minecraft", c.String("sort"))
// 	}

// 	if c.Bool("first") {
// 		result = []*kleister.Forge{
// 			result[0],
// 		}
// 	}

// 	if c.Bool("last") {
// 		result = []*kleister.Forge{
// 			result[len(result)-1],
// 		}
// 	}

// 	if c.Bool("xml") {
// 		res, err := xml.MarshalIndent(result, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if c.Bool("json") {
// 		res, err := json.MarshalIndent(result, "", "  ")

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Fprintf(os.Stdout, "%s\n", res)
// 		return nil
// 	}

// 	if len(result) == 0 {
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

// 	for _, record := range result {
// 		err := tmpl.Execute(os.Stdout, record)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // ForgeRefresh provides the sub-command to refresh the Forge versions.
// func ForgeRefresh(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ForgeRefresh()

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully refreshed\n")
// 	return nil
// }

// // ForgeBuildList provides the sub-command to list builds of the Forge.
// func ForgeBuildList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.ForgeBuildList(
// 		kleister.ForgeBuildParams{
// 			Forge: GetIdentifierParam(c),
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

// // ForgeBuildAppend provides the sub-command to append a build to the Forge.
// func ForgeBuildAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ForgeBuildAppend(
// 		kleister.ForgeBuildParams{
// 			Forge: GetIdentifierParam(c),
// 			Pack:  GetPackParam(c),
// 			Build: GetBuildParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully appended to Forge\n")
// 	return nil
// }

// // ForgeBuildRemove provides the sub-command to remove a build from the Forge.
// func ForgeBuildRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.ForgeBuildDelete(
// 		kleister.ForgeBuildParams{
// 			Forge: GetIdentifierParam(c),
// 			Pack:  GetPackParam(c),
// 			Build: GetBuildParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully removed from Forge\n")
// 	return nil
// }
