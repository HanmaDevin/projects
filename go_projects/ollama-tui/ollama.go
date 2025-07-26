package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Ollama struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func NewOllamaModel() *Ollama {
	return &Ollama{
		Stream: false,
	}
}

func getResponse(ollama *Ollama) string {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(ollama)

	req, err := http.Post("http://localhost:11434/api/generate", "application/json", body)
	if err != nil {
		log.Fatal("<-!--- Could not read Body --->\n")
	}

	defer req.Body.Close()

	return ""
}
