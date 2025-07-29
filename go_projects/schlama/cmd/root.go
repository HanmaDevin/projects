/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/HanmaDevin/schlama/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "schlama",
	Short: "A better ollama user interface.",
	Long:  `Schlama is a cli or tui user interface, depending on what you perfer, which allows for easy communication with the local ollama api. Basically an easier way to chat with local model or install new ones.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	var home, _ = os.UserHomeDir()
	var config_Path string = filepath.Dir(home + "/.config/schlama/")
	if _, err := os.Stat(config_Path); os.IsNotExist(err) {
		err := os.MkdirAll(config_Path, 0755)
		if err != nil {
			fmt.Println("Error creating config directory:", err.Error())
			os.Exit(-1)
		}
	}

	if _, err := os.Stat(config_Path + "/config.yaml"); os.IsNotExist(err) {
		config.WriteConfig(config.Config{
			Prompt: "What is the meaning of life?",
			Model:  "",
			Stream: false,
		})
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
