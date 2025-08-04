/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/spf13/cobra"
)

var file string

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
				fmt.Println(styles.TableBorder(styles.HintStyle("No model specified in config. Please set a model using 'schlama select <model_name>'.")))
				return
			}

			var f []byte
			var err error
			if file != "" {
				f, err = os.ReadFile(file)
				if err != nil {
					fmt.Println(styles.TableBorder(styles.ErrorStyle("Not able to read the specified file!")))
					os.Exit(1)
				}
			}

			body.Prompt = args[0] + "\n" + string(f)
			resp, err := ollama.GetResponse(body)
			if err != nil {
				fmt.Println(styles.TableBorder(styles.ErrorStyle(err.Error())))
				return
			}
			ollama.PrintMarkdown(resp)
		}
	},
}

func init() {
	promptCmd.Flags().StringVarP(&file, "file", "f", "", "Prompt with file content")
	rootCmd.AddCommand(promptCmd)
}
