/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/HanmaDevin/ollamacltui/service"
	"github.com/spf13/cobra"
)

var limit int64

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models for download.",
	Long:  `List gets all the available models from ollama.com and displays them.`,
	Run: func(cmd *cobra.Command, args []string) {
		models := service.ListModels()
		fmt.Println("All available Models:")
		for i := range limit {
			fmt.Println(models[i])
		}
	},
}

func init() {
	listCmd.Flags().Int64VarP(&limit, "limit", "l", 50, "Limit the size of output.")
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
