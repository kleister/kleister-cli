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

// // tmplMinecraftList represents a row within forge listing.
// var tmplMinecraftList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Version: {{ .Version }}
// Type: {{ .Type }}
// `

// // tmplMinecraftBuildList represents a row within minecraft build listing.
// var tmplMinecraftBuildList = "Slug: \x1b[33m{{ .Slug }}\x1b[0m" + `
// ID: {{ .ID }}
// Name: {{ .Name }}
// Pack: {{with .Pack}}{{ . }}{{else}}n/a{{end}}
// `

// // Minecraft provides the sub-command for the Minecraft API.
// func Minecraft() *cli.Command {
// 	return &cli.Command{
// 		Name:  "minecraft",
// 		Usage: "Minecraft related sub-commands",
// 		Subcommands: []*cli.Command{
// 			{
// 				Name:      "list",
// 				Aliases:   []string{"ls"},
// 				Usage:     "List all Minecraft versions",
// 				ArgsUsage: " ",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:  "sort",
// 						Value: "slug",
// 						Usage: "Sort by this field",
// 					},
// 					&cli.StringFlag{
// 						Name:  "format",
// 						Value: tmplMinecraftList,
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
// 					&cli.StringFlag{
// 						Name:  "filter",
// 						Value: "",
// 						Usage: "Filter by values",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "first",
// 						Value: false,
// 						Usage: "Return only first record",
// 					},
// 					&cli.BoolFlag{
// 						Name:  "last",
// 						Value: false,
// 						Usage: "Return only last record",
// 					},
// 				},
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, MinecraftList)
// 				},
// 			},
// 			{
// 				Name:      "refresh",
// 				Aliases:   []string{"ref"},
// 				Usage:     "Refresh Minecraft versions",
// 				ArgsUsage: " ",
// 				Action: func(c *cli.Context) error {
// 					return Handle(c, MinecraftRefresh)
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
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Minecraft ID or slug to list builds",
// 							},
// 							&cli.StringFlag{
// 								Name:  "format",
// 								Value: tmplMinecraftBuildList,
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
// 							return Handle(c, MinecraftBuildList)
// 						},
// 					},
// 					{
// 						Name:      "append",
// 						Usage:     "Append a build to Minecraft",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Minecraft ID or slug to append to",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "Pack ID or slug to append",
// 							},
// 							&cli.StringFlag{
// 								Name:  "build, b",
// 								Value: "",
// 								Usage: "Build ID or slug to append",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, MinecraftBuildAppend)
// 						},
// 					},
// 					{
// 						Name:      "remove",
// 						Aliases:   []string{"rm"},
// 						Usage:     "Remove a build from Minecraft",
// 						ArgsUsage: " ",
// 						Flags: []cli.Flag{
// 							&cli.StringFlag{
// 								Name:  "id, i",
// 								Value: "",
// 								Usage: "Minecraft ID or slug to remove from",
// 							},
// 							&cli.StringFlag{
// 								Name:  "pack, p",
// 								Value: "",
// 								Usage: "Pack ID or slug to remove",
// 							},
// 							&cli.StringFlag{
// 								Name:  "build, b",
// 								Value: "",
// 								Usage: "Build ID or slug to remove",
// 							},
// 						},
// 						Action: func(c *cli.Context) error {
// 							return Handle(c, MinecraftBuildRemove)
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // MinecraftList provides the sub-command to list all Minecraft versions.
// func MinecraftList(c *cli.Context, client kleister.ClientAPI) error {
// 	var (
// 		result []*kleister.Minecraft
// 	)

// 	records, err := client.MinecraftList()

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
// 			params["Type"] = record.Type

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
// 			kleister.MinecraftBySlug(
// 				result,
// 			),
// 		)
// 	case "version":
// 		sort.Sort(
// 			kleister.MinecraftByVersion(
// 				result,
// 			),
// 		)
// 	case "type":
// 		sort.Sort(
// 			kleister.MinecraftByType(
// 				result,
// 			),
// 		)
// 	default:
// 		return fmt.Errorf("the sort value %s is invalid, can be slug, version or type", c.String("sort"))
// 	}

// 	if c.Bool("first") {
// 		result = []*kleister.Minecraft{
// 			result[0],
// 		}
// 	}

// 	if c.Bool("last") {
// 		result = []*kleister.Minecraft{
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

// // MinecraftRefresh provides the sub-command to refresh the Minecraft versions.
// func MinecraftRefresh(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.MinecraftRefresh()

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully refreshed\n")
// 	return nil
// }

// // MinecraftBuildList provides the sub-command to list builds of the Minecraft.
// func MinecraftBuildList(c *cli.Context, client kleister.ClientAPI) error {
// 	records, err := client.MinecraftBuildList(
// 		kleister.MinecraftBuildParams{
// 			Minecraft: GetIdentifierParam(c),
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

// // MinecraftBuildAppend provides the sub-command to append a build to the Minecraft.
// func MinecraftBuildAppend(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.MinecraftBuildAppend(
// 		kleister.MinecraftBuildParams{
// 			Minecraft: GetIdentifierParam(c),
// 			Pack:      GetPackParam(c),
// 			Build:     GetBuildParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully appended to Minecraft\n")
// 	return nil
// }

// // MinecraftBuildRemove provides the sub-command to remove a build from the Minecraft.
// func MinecraftBuildRemove(c *cli.Context, client kleister.ClientAPI) error {
// 	err := client.MinecraftBuildDelete(
// 		kleister.MinecraftBuildParams{
// 			Minecraft: GetIdentifierParam(c),
// 			Pack:      GetPackParam(c),
// 			Build:     GetBuildParam(c),
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(os.Stderr, "Successfully removed from Minecraft\n")
// 	return nil
// }
