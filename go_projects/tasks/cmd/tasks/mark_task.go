package tasks

import (
	"Go_Projects/todo_list/pkg/tasks"
	"strconv"

	"github.com/spf13/cobra"
)

var markCmd = &cobra.Command{
	Use:   "mark <ID>",
	Short: "Mark a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		tasks.MarkTaskDone(id)
	},
}

func init() {
	rootCmd.AddCommand(markCmd)
}
