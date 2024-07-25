package main

import (
	"cli/cmd"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "short info",
	Long:  `long info`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Do Stuff Here")
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	cmd.CmdTimes.Flags().IntVarP(&cmd.EchoTimes, "times", "t", 1, "times to echo the input")

	// 两个顶层的命令，和一个cmdEcho命令下的子命令cmdTimes
	RootCmd.AddCommand(cmd.CmdPrint, cmd.CmdEcho)
	cmd.CmdEcho.AddCommand(cmd.CmdTimes)
}

func initConfig() {
	// 勿忘读取config文件，无论是从cfgFile还是从home文件
	if cfgFile != "" {
		viper.SetConfigName(cfgFile)
	} else {
		// 找到home文件
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 在home文件夹中搜索以“.cobra”为名称的config
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	// 读取符合的环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can not read config:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
