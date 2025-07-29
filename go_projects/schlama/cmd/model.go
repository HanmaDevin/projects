/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"

	"github.com/HanmaDevin/schlama/config"
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Show the current selected model.",
	Long:  `Show the current selected model.`,
	Run: func(cmd *cobra.Command, args []string) {
		body := config.ReadConfig()
		fmt.Println("Current model:", body.Model)
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
