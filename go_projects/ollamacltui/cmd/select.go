/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Selects which model to chat with.",
	Long:  `This command sets the model to chat with, short description explains it completely pretty much`,
	Run: func(cmd *cobra.Command, args []string) {
		Ollama.SetModel(args[0])
		fmt.Println("select called")
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
