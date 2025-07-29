/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/HanmaDevin/schlama/ollama"
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

		fmt.Printf("%-25s %-20s\n", "MODEL NAME", "SIZES")
		fmt.Printf("%-25s %-20s\n", strings.Repeat("-", 25), strings.Repeat("-", 20))
		for i, model := range models {
			if i >= int(limit) {
				break
			}
			fmt.Printf("%-25s %-20s\n", model.Name, strings.Join(model.Sizes, ", "))
		}
	},
}

func init() {
	listCmd.Flags().Int64VarP(&limit, "limit", "l", 50, "Limit the output.")
	rootCmd.AddCommand(listCmd)
}
