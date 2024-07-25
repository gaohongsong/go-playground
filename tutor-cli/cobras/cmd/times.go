package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var EchoTimes int

var CmdTimes = &cobra.Command{
	Use:   "times [# times] [string to echo]",
	Short: "Echo anything to the screen more times",
	Long: `echo things multiple times back to the user by providing
a count and a string.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < EchoTimes; i++ {
			fmt.Println("Echo: " + strings.Join(args, " "))
		}
	},
}
