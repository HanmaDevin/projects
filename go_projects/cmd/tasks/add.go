package tasks

import (
	"Go_Projects/todo_list/pkg/tasks"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add <description>",
	Aliases: []string{"a"},
	Short:   "Add a task with a description",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tasks.AddTask(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
