package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

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

	csvWriter := csv.NewWriter(file)

	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")
	done := false

	file, err = os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	content, _ := csv.NewReader(file).ReadAll()
	id := len(content)

	data := []string{strconv.Itoa(id), desc, createdAt, strconv.FormatBool(done)}

	csvWriter.Write(data)
	csvWriter.Flush()
	file.Close()
}

func DeleteTask(id int) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)
	header, err := csvReader.Read() // ignore the header first
	if err != nil {
		panic(err)
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Filter out the row with the matching ID
	var updatedRecords [][]string
	for _, record := range records {
		if len(record) > 0 {
			recordID, err := strconv.Atoi(record[0])
			if err != nil {
				panic(err)
			}
			if recordID != id { // Keep rows that don't match the given ID
				updatedRecords = append(updatedRecords, record)
			}
		}
	}

	// Rewrite the file with the updated records
	file.Close()

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644) // Open in write mode with truncation
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	err = csvWriter.Write(header) // rewrite the header
	if err != nil {
		panic(err)
	}

	err = csvWriter.WriteAll(updatedRecords) // Write all remaining records
	if err != nil {
		panic(err)
	}
	csvWriter.Flush()
}

func PrintTasks() {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	// Read the header row and ignore it
	_, err = csvReader.Read()
	if err != nil {
		panic(err)
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%3s\t%20s\t%20s\t%s\n", "ID", "Description", "CreatedAt", "IsComplete")
	fmt.Println("------------------------------------------------------------------")
	for _, record := range records {
		isComplete, err := strconv.ParseBool(record[3])
		if err != nil {
			panic(err)
		}

		// Parse the CreatedAt timestamp
		createdAt, err := time.Parse("2006-01-02 15:04:05", record[2])
		if err != nil {
			panic(err)
		}

		// Convert CreatedAt to a human-readable relative time
		humanReadableTime := timediff.TimeDiff(createdAt)

		if !isComplete {
			fmt.Printf("%3s\t%20s\t%20s\t%s\n", record[0], record[1], humanReadableTime, record[3])
		}
	}
}

func PrintAllTasks() {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)
	_, err = csvReader.Read() // skip header
	if err != nil {
		panic(err)
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%3s\t%20s\t%20s\t%s\n", "ID", "Description", "CreatedAt", "IsComplete")
	fmt.Println("------------------------------------------------------------------")

	for _, record := range records {
		// Parse the CreatedAt timestamp
		createdAt, err := time.Parse("2006-01-02 15:04:05", record[2])
		if err != nil {
			panic(err)
		}

		// Convert CreatedAt to a human-readable relative time
		humanReadableTime := timediff.TimeDiff(createdAt)

		fmt.Printf("%3s\t%20s\t%20s\t%s\n", record[0], record[1], humanReadableTime, record[3])
	}

}

func MarkTaskDone(id int) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)
	header, err := csvReader.Read() // ignore the header first
	if err != nil {
		panic(err)
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	var updatedRecords [][]string
	for _, record := range records {
		recordID, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		if recordID == id {
			record[3] = "true" // set the isComplete value to true
		}
		updatedRecords = append(updatedRecords, record)
	}

	// Rewrite the file with the updated records
	file.Close()

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644) // Open in write mode with truncation
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	err = csvWriter.Write(header) // rewrite the header
	if err != nil {
		panic(err)
	}

	err = csvWriter.WriteAll(updatedRecords) // Write all remaining records
	if err != nil {
		panic(err)
	}
	csvWriter.Flush()
}
