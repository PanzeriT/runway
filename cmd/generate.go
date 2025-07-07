package cmd

import (
	"os"

	"github.com/panzerit/runway/internal/generator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir := Must(cmd.Flags().GetString("output"))
		if outputDir == "" {
			outputDir = Must(os.Getwd())
		}

		generator.New(generator.WithOutputDir(outputDir)).
			Run()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
