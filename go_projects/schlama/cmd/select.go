/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
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
				fmt.Println("Model not found locally. Pulling model...")
				err := ollama.PullModel(args[0])
				if err != nil {
					fmt.Println(err)
					fmt.Println("Here is a list of available models:")
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
			fmt.Println("Model set to:", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
