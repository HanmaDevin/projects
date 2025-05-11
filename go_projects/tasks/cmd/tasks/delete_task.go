package tasks

import (
	"Go_Projects/todo_list/pkg/tasks"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <ID>",
	Short:   "Delete task with given ID",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"del"},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		tasks.DeleteTask(id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
