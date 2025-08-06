/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"
	"regexp"
	"strings"

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
			nameReg := regexp.MustCompile(`[\w\-]+\d?\.?\d?`)
			name := nameReg.FindString(args[0])

			labelReg := regexp.MustCompile(`:\w+`)
			label := labelReg.FindString(args[0])
			if label == "" {
				label = ":latest"
			}

			model := name + label

			if !ollama.IsModelPresent(model) {
				fmt.Println(styles.HintStyle("Model not found locally. Pulling model..."))
				err := ollama.PullModel(model)
				if err != nil {
					fmt.Println(styles.ErrorStyle(err.Error()))
					fmt.Println(styles.HintStyle("Here is a list of available models:"))
					models := ollama.ListModels()

					var rows []string
					header := fmt.Sprintf("%-25s %-40s", "MODEL NAME", "SIZES")
					rows = append(rows, styles.HeaderStyle(header))
					divider := fmt.Sprintf("%-25s %-40s", strings.Repeat("-", 25), strings.Repeat("-", 40))
					rows = append(rows, styles.RowStyle(divider))
					for i, model := range models {
						if i >= int(limit) {
							break
						}
						line := fmt.Sprintf("%-25s %-40s", model.Name, strings.Join(model.Sizes, ", "))
						rows = append(rows, styles.RowStyle(line))
					}
					table := strings.Join(rows, "\n")
					fmt.Println(table)
					return
				}
			}

			cfg := config.Config{
				Model: model,
				Msg: ollama.Message{
					Role:    "user",
					Content: "",
					Image:   nil,
				},
				Stream: false,
			}
			config.WriteConfig(cfg)
			out := fmt.Sprintf("Current Model: %s", cfg.Model)
			fmt.Println(styles.OutputStyle(out))
		}
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
