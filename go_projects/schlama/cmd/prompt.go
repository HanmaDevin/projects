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

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Get a response with specified prompt.",
	Long:  `Makes an api call to localhost:11343/api/generate and outputs the response in a more readable fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
		} else {
			body := config.ReadConfig()
			if body.Model == "" {
				fmt.Println(styles.HintStyle("No model specified in config. Please set a model using 'schlama select <model_name>'."))
				return
			}
			body.Prompt = args[0]
			resp, err := ollama.GetResponse(body)
			if err != nil {
				fmt.Println(err)
				return
			}
			ollama.PrintMarkdown(resp)
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
