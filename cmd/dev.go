package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dev called")
	},
}

var createTemplateCmd = &cobra.Command{
	Use:   "create-template path",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copy-template called")
		if len(args) < 1 {
			cmd.Println("Error: at least one argument is required")
			cmd.Usage()
			return
		}
		cmd.Println("Creating template at path:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(createTemplateCmd)
}
