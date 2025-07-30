/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Selects which model to chat with.",
	Long:  `This command sets the model to chat with. To list available model use 'local' command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
		} else {
			// Check if the model is present in the local models
			if !ollama.IsModelPresent(args[0]) {
				fmt.Println(styles.TableBorder(styles.HintStyle("Model not found locally. Pulling model...")))
				err := ollama.PullModel(args[0])
				if err != nil {
					fmt.Println(styles.TableBorder(styles.ErrorStyle(err.Error())))
					fmt.Println(styles.TableBorder(styles.HintStyle("Here is a list of available models:")))
					ollama.ListLocalModels()
					return
				}
			}
			body := config.ReadConfig()

			cfg := config.Config{
				Prompt: body.Prompt,
				Model:  args[0],
				Stream: body.Stream,
			}
			config.WriteConfig(cfg)
			out := fmt.Sprintf("Current Model: %s", cfg.Model)
			fmt.Println(styles.TableBorder(styles.OutputStyle(out)))
		}
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
