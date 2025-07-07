package cmd

import (
	"fmt"
	"os"

	"github.com/panzerit/runway/internal/generator"
	"github.com/spf13/cobra"
)

func Must[T any](res T, err error) T {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return any(res).(T)
}

var initCmd = &cobra.Command{
	Use:     "init <module_name>",
	Example: "ruwway init mymodule  # inialize a new Go project",
	Short:   "Start a new runway project from scratch; create the go.mod file",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force := Must(cmd.Flags().GetBool("force"))

		outputDir := Must(cmd.Flags().GetString("output"))
		if outputDir == "" {
			outputDir = Must(os.Getwd())
		}

		fileCount := Must(os.ReadDir(outputDir))
		if len(fileCount) > 0 && !force {
			cmd.Println("Warning: The current directory is not empty. Use --force to override this warning.")
			os.Exit(2)
		}

		generator.New(generator.WithOutputDir(outputDir)).
			Init(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("force", "f", false, "Force initialization even if the directory is not empty")
}
