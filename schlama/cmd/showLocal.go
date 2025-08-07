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
	Short: "Lists all locally downloaded models.",
	Long:  `Lists all locally downloaded models on the current machine. If you want to downlaod more models check out the 'schlama list' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		ollama.ListLocalModels()
	},
}

func init() {
	rootCmd.AddCommand(showLocalCmd)
}
