/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/spf13/cobra"
)

var limit int64

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models for download.",
	Long:  `List gets all the available models from ollama.com and displays them.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	listCmd.Flags().Int64VarP(&limit, "limit", "l", 50, "Limit the output.")
	rootCmd.AddCommand(listCmd)
}
