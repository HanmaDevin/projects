package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/mergestat/timediff"
)

// global variables
var (
	err      error
	homeDir  string
	filename string
)

// init function to set global veriables
func init() {
	homeDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filename = homeDir + "\\.go_data\\tasks.csv"
}

func AddTask(desc string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Fscanf(os.Stderr, "File: %s does not exist!!!", filename)
	} else if err != nil {
		panic(err)
	}

	CreatedAt := fileInfo.ModTime()
	diff := timediff.TimeDiff(CreatedAt)
	done := false

	content, _ := csv.NewReader(file).ReadAll()
	id := len(content)

	data := []string{strconv.Itoa(id), desc, diff, strconv.FormatBool(done)}

	csvWriter.Write(data)
	csvWriter.Flush()
}

func DeleteTask(id int) {}

func PrintTasks() {}

func PrintAllTasks() {}

func MarkTaskDone(id int) {}
