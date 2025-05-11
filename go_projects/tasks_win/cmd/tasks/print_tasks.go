package tasks

import (
	"Go_Projects/todo_list/pkg/tasks"

	"github.com/spf13/cobra"
)

var all bool

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print tasks",
	Long:  "Prints tasks based on their completion status",
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			tasks.PrintAllTasks()
		} else {
			tasks.PrintTasks()
		}
	},
}

func init() {
	// Add the --all flag to the print command
	printCmd.Flags().BoolVarP(&all, "all", "a", false, "Print all tasks, including completed ones")
	rootCmd.AddCommand(printCmd)
}
