package main

import (
	"Go_Projects/todo_list/cmd/tasks"
	"os"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	check(err)

	dir := homeDir + "\\.go_data"
	filename := homeDir + "\\.go_data\\tasks.csv"
	header := "ID,Description,CreatedAt,IsComplete\n"

	if !exists(dir) {
		err := os.Mkdir(dir, 0755)
		check(err)
	}

	if !exists(filename) {
		_, err := os.Create(filename)
		check(err)
	}
	os.WriteFile(filename, []byte(header), 0644)

	tasks.Execute()
}
