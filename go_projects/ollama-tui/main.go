package main

import "fmt"

func main() {
	ollama := NewOllamaModel()
	ollama.setModel("deepseek-r1:14b")
	ollama.setPrompt("hello")

	resp := getResponse(ollama)

	fmt.Println(resp)
}
