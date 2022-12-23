package command

import (
	"os"

	"github.com/kleister/kleister-cli/pkg/version"
	"github.com/urfave/cli/v2"
)

const (
	defaultServer = "http://localhost:8080"
)

// Run parses the command line arguments and executes the program.
func Run() error {
	app := &cli.App{
		Name:    "kleister-cli",
		Version: version.String,
		Usage:   "Manage mod packs for minecraft",
		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas@webhippie.de",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server, s",
				Value:   defaultServer,
				Usage:   "API server",
				EnvVars: []string{"KLEISTER_SERVER"},
			},
			&cli.StringFlag{
				Name:    "token, t",
				Value:   "",
				Usage:   "API token",
				EnvVars: []string{"KLEISTER_TOKEN"},
			},
		},
		Commands: []*cli.Command{
			Profile(),
			User(),
			Team(),
			// Build(),
			// Client(),
			// Forge(),
			// Key(),
			// Minecraft(),
			// Mod(),
			// Pack(),
			// Version(),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the version",
	}

	return app.Run(os.Args)
}
