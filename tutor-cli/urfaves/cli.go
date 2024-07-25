package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"urfaves/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-netdisk"
	app.Usage = "A simple net-disk service"
	app.Version = "v1"
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Migrate,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
