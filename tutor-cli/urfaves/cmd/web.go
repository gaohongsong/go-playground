package cmd

import (
	"github.com/urfave/cli"
	"log"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: `go-netdisk web server provide http service`,
	Action:      runWeb,
	Flags: []cli.Flag{
		intFlag("port", 5000, "Temporary port number to prevent conflict", []string{"p"}),
		stringFlag("config", "", "Custom configuration file path", []string{"c"}),
	},
}

func runWeb(c *cli.Context) error {
	log.Println("init database start")
	return nil
}
