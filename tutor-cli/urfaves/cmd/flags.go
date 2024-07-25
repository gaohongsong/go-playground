package cmd

import (
	"github.com/urfave/cli"
	"time"
)

func stringFlag(name, value, usage string, aliases []string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string, aliases []string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

//nolint:deadcode,unused
func intFlag(name string, value int, usage string, aliases []string) *cli.IntFlag {
	return &cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

//nolint:deadcode,unused
func durationFlag(name string, value time.Duration, usage string, aliases []string) *cli.DurationFlag {
	return &cli.DurationFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
