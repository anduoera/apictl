package cmd

import (
	"apictl/core"
	"github.com/spf13/cobra"
	"os"
)

var ctx *core.Context
var rootCmd = &cobra.Command{
	Use:   "apictl",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.SetLoggerLevelStr(args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	ctx = core.GetInstance()
}
