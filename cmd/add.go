package cmd

import (
	"apictl/logic"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   logic.AddCommandUse,
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.SetLoggerLevelStr(args)
		command := logic.GetCommandTag(cmd.Use, ctx, args)
		command.Run()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
