package command

import (
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/urfave/cli/v2"
)

// sprigFuncMap provides template helpers provided by sprig.
var sprigFuncMap = sprig.TxtFuncMap()

// globalFuncMap provides global template helper functions.
var globalFuncMap = template.FuncMap{
	// "buildlist": func(s []*kleister.Build) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		if row.Pack != nil {
	// 			res = append(res, fmt.Sprintf("%s@%s", row.Pack.Slug, row.String()))
	// 		} else {
	// 			res = append(res, row.String())
	// 		}
	// 	}

	// 	return strings.Join(res, ", ")
	// },
	// "clientlist": func(s []*kleister.Client) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		res = append(res, row.String())
	// 	}

	// 	return strings.Join(res, ", ")
	// },
	// "modlist": func(s []*kleister.Mod) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		res = append(res, row.String())
	// 	}

	// 	return strings.Join(res, ", ")
	// },
	// "packlist": func(s []*kleister.Pack) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		res = append(res, row.String())
	// 	}

	// 	return strings.Join(res, ", ")
	// },
	// "teamlist": func(s []*kleister.Team) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		res = append(res, row.String())
	// 	}

	// 	return strings.Join(res, ", ")
	// },
	// "userlist": func(s []*kleister.User) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		res = append(res, row.String())
	// 	}

	// 	return strings.Join(res, ", ")
	// },
	// "versionlist": func(s []*kleister.Version) string {
	// 	res := []string{}

	// 	for _, row := range s {
	// 		if row.Mod != nil {
	// 			res = append(res, fmt.Sprintf("%s@%s", row.Mod.Slug, row.String()))
	// 		} else {
	// 			res = append(res, row.String())
	// 		}
	// 	}

	// 	return strings.Join(res, ", ")
	// },
}

// GetIdentifierParam checks and returns the record id/slug parameter.
func GetIdentifierParam(c *cli.Context) string {
	val := c.String("id")

	if val == "" {
		fmt.Fprintf(os.Stderr, "Error: you must provide an ID or a slug.\n")
		os.Exit(1)
	}

	return val
}

// GetUserParam checks and returns the user id/slug parameter.
func GetUserParam(c *cli.Context) string {
	val := c.String("user")

	if val == "" {
		fmt.Fprintf(os.Stderr, "Error: you must provide a user ID or slug.\n")
		os.Exit(1)
	}

	return val
}

// GetTeamParam checks and returns the team id/slug parameter.
func GetTeamParam(c *cli.Context) string {
	val := c.String("team")

	if val == "" {
		fmt.Fprintf(os.Stderr, "Error: you must provide a team ID or slug.\n")
		os.Exit(1)
	}

	return val
}

// GetModParam checks and returns the mod id/slug parameter.
func GetModParam(c *cli.Context) string {
	val := c.String("mod")

	if val == "" {
		fmt.Println("Error: you must provide a mod id or slug.")
		os.Exit(1)
	}

	return val
}

// GetVersionParam checks and returns the version id/slug parameter.
func GetVersionParam(c *cli.Context) string {
	val := c.String("version")

	if val == "" {
		fmt.Println("Error: you must provide a version id or slug.")
		os.Exit(1)
	}

	return val
}

// GetPackParam checks and returns the pack id/slug parameter.
func GetPackParam(c *cli.Context) string {
	val := c.String("pack")

	if val == "" {
		fmt.Println("Error: you must provide a pack id or slug.")
		os.Exit(1)
	}

	return val
}

// GetBuildParam checks and returns the build id/slug parameter.
func GetBuildParam(c *cli.Context) string {
	val := c.String("build")

	if val == "" {
		fmt.Println("Error: you must provide a build id or slug.")
		os.Exit(1)
	}

	return val
}

// GetClientParam checks and returns the client id/slug parameter.
func GetClientParam(c *cli.Context) string {
	val := c.String("client")

	if val == "" {
		fmt.Println("Error: you must provide a client id or slug.")
		os.Exit(1)
	}

	return val
}
