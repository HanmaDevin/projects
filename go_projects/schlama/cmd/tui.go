/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start a TUI application",
	Long:  `The tui command starts a terminal user interface application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tui called")
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
