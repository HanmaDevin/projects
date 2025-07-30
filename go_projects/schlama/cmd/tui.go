/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"

	"github.com/HanmaDevin/schlama/styles"
	"github.com/HanmaDevin/schlama/tui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start a TUI application",
	Long:  `The tui command starts a terminal user interface application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.Start()
		if err != nil {
			out := styles.ErrorStyle(fmt.Sprintf("Error starting TUI: %v", err))
			fmt.Println(out)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
