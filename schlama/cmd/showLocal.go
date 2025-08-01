/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/spf13/cobra"
)

// showLocalCmd represents the show command
var showLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Lists all downloaded models.",
	Long:  `Lists all downloaded models on current machine. If you want to downlaod more models check out the 'list' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		ollama.ListLocalModels()
	},
}

func init() {
	rootCmd.AddCommand(showLocalCmd)
}
