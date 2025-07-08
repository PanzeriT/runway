package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/panzerit/runway/internal/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "runway",
	Short: "A brief description of your application",
}

var once sync.Once

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringP("output", "o", "", "Specify the output directory for the new project (default is current directory)")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		once.Do(func() {
			path, err := rootCmd.PersistentFlags().GetString("output")
			generator.CheckError(err, generator.ErrInvalidFlag)

			if path == "" {
				path, err = os.Getwd()
				generator.CheckError(err, generator.ErrUnableToDetermineWorkingDirectory)
			}

			cfgFile := filepath.Join(path, "config.yaml")
			viper.SetConfigFile(cfgFile)

			err = viper.ReadInConfig()
			generator.CheckError(err, generator.ErrMissingConfiguration,
				fmt.Sprintf("No configuration file found (%s). Please run 'runway init <module_name>' to start a new project.", cfgFile))
		})
	}
}
