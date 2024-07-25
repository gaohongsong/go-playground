package cmd

import (
	"github.com/urfave/cli"
	"log"
)

var Migrate = cli.Command{
	Name:        "migrate",
	Usage:       "Migrate init database",
	Description: `Backup create table and insert initial data.`,
	Action:      runMigrate,
	Flags: []cli.Flag{
		stringFlag("config", "", "Custom configuration file path", []string{"c"}),
		boolFlag("verbose, v", "Show process details", []string{"v"}),
	},
}

func runMigrate(c *cli.Context) error {
	log.Println("init database start")
	return nil
}
