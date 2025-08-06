/*
Copyright © 2025 Devin Brunk Cardosa
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
	"github.com/spf13/cobra"
)

var file string
var directory string

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prompt the model with a message and/or file",
	Long:  `Makes an api call to localhost:11343/api/generate and outputs the response in a more readable fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
		} else {
			body := config.ReadConfig()
			if body.Model == "" {
				fmt.Println(styles.HintStyle("No model specified in config. Please set a model using 'schlama select <model_name>'."))
				return
			}

			body.Msg.Content = args[0]

			var f []byte
			var err error
			if cmd.Flags().Changed("file") {
				fmt.Println(styles.HintStyle("Reading file: " + file))
				f, err = os.ReadFile(file)
				if err != nil {
					fmt.Println(styles.ErrorStyle("Not able to read the specified file!"))
					os.Exit(1)
				}
				body.Msg.Content += "\n" + string(f)
			}

			if cmd.Flags().Changed("directory") {
				fmt.Println(styles.HintStyle("Reading directory: " + directory))
				data, err := getDirContent(directory)
				if err != nil {
					fmt.Println(styles.ErrorStyle("Not able to read the specified directory!"))
					os.Exit(1)
				}
				body.Msg.Content += "\n" + data
			}

			resp, err := ollama.GetResponse(body)
			if err != nil {
				fmt.Println(styles.ErrorStyle(err.Error()))
				return
			}
			ollama.PrintMarkdown(resp)
		}
	},
}

func getDirContent(root string) (string, error) {
	var sb strings.Builder
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Println(styles.HintStyle("Reading file: " + path))
			sb.WriteString("File: " + filepath.Base(path) + "\n")
			sb.Write(content)
			sb.WriteString("\n\n")
		}
		return nil
	})
	return sb.String(), err
}

func init() {
	promptCmd.Flags().StringVarP(&file, "file", "f", "", "Prompt with file content")
	promptCmd.Flags().StringVarP(&directory, "directory", "d", "", "Prompt with directory content")
	rootCmd.AddCommand(promptCmd)
}
