package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/solderapp/solder-cli/solder"
)

// Forge provides the sub-command for the forge API.
func Forge() cli.Command {
	return cli.Command{
		Name:  "forge",
		Usage: "Forge related sub-commands",
		Subcommands: []cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all Forge versions",
				Action: func(c *cli.Context) {
					Handle(c, ForgeList)
				},
			},
			{
				Name:    "refresh",
				Aliases: []string{"ref"},
				Usage:   "Refresh the Forge versions",
				Action: func(c *cli.Context) {
					Handle(c, ForgeRefresh)
				},
			},
		},
	}
}

// ForgeList provides the sub-command to list all Forge versions.
func ForgeList(c *cli.Context, client solder.API) error {
	records, err := client.ForgeList()

	if err != nil || len(records) == 0 {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"ID", "Version", "Minecraft"})

	for _, record := range records {
		table.Append(
			[]string{
				strconv.FormatInt(record.ID, 10),
				record.Version,
				record.Minecraft,
			},
		)
	}

	table.Render()
	return nil
}

// ForgeRefresh provides the sub-command to refresh the Forge versions.
func ForgeRefresh(c *cli.Context, client solder.API) error {
	msg, err := client.ForgeRefresh()

	if err != nil {
		return err
	}

	fmt.Println(msg.Message)
	return nil
}
