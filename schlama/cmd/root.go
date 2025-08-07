/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "schlama",
	Short: "A better ollama user interface.",
	Long:  `Schlama is a CLI and a web-chat app, depending on what you perfer, which allows for easy communication with local LLMs. It allows file/directory input and images are also supported (Only works with multimodal models). Basically an easier way to chat with local LLMs and install new ones. For more control over the models please use the ollama CLI. This is just a simpler way to interact with the ollama api and having a bit of control over the models.`,
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
	if ollama.IsOllamaRunning() {
		var home, _ = os.UserHomeDir()
		var config_Path string = filepath.Dir(home + "/.config/schlama/")
		if _, err := os.Stat(config_Path); os.IsNotExist(err) {
			err := os.MkdirAll(config_Path, 0755)
			if err != nil {
				fmt.Println(styles.ErrorStyle("Error creating config directory: ~/.config/schlama/config.yaml"))
				os.Exit(-1)
			}
		}

		if _, err := os.Stat(config_Path + "/config.yaml"); os.IsNotExist(err) {
			config.WriteConfig(config.Config{
				Model: "",
			})
		}
	} else {
		fmt.Println(styles.ErrorStyle("Ollama is not running."))
		fmt.Println(styles.OutputStyle("Please start ollama first."))
		fmt.Println(styles.HintStyle("You can start ollama with the command: 'ollama serve'"))
		fmt.Println(styles.HintStyle("Or you can install ollama with the command: 'curl -fsSL https://ollama.com/install.sh | sh'"))
		os.Exit(1)
	}
}
