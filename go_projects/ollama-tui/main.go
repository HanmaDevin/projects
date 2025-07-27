package main

func main() {
	ollama := NewOllamaModel()
	ollama.setModel("deepseek-r1:14b")
	ollama.setPrompt("hello")

	listLocalModels()
}
