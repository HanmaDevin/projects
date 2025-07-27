package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/net/html"
)

const ollama_api = "http://localhost:11434/api/generate"

type OllamaModel struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Resp string `json:"response"`
}

func NewOllamaModel() *OllamaModel {
	return &OllamaModel{
		Stream: false,
	}
}

func getResponse(ollama *OllamaModel) string {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(ollama)

	c := http.Client{Timeout: time.Second * 15}

	resp, err := c.Post(ollama_api, "application/json", body)
	if err != nil {
		log.Fatal("<-!-- Post request to ollama api failed --->\n")
	}

	var ai Response
	json.NewDecoder(resp.Body).Decode(&ai)
	defer resp.Body.Close()

	return clean(ai.Resp)
}

func (o *OllamaModel) setPrompt(prompt string) {
	o.Prompt = prompt
}

func (o *OllamaModel) setModel(model string) {
	o.Model = model
}

func pullModel(model string) {
	cmd := exec.Command("ollama", "pull", model)
	fmt.Println("Pulling manifest...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("<-!-- Could not pull model: %s --->\n", model)
	}
	fmt.Println("Finished!")
}

func listLocalModels() {
	cmd := exec.Command("ollama", "list")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal("<-!-- Could not run 'ollama list' --->\n")
	}
}

func listModels() []string {
	c := http.Client{Timeout: time.Second * 15}
	resp, err := c.Get("https://ollama.com/library?sort=popular")
	if err != nil {
		log.Fatal("<-!-- Could not get a response from ollama.com --->\n")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("<-!-- Could not read response from ollama.com --->\n")
	}

	var models []string
	doc, err := html.Parse(bytes.NewReader(b))
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "title" {
					models = append(models, a.Val)
					break
				}
			}
		}
	}

	return models
}

func clean(s string) string {
	return s[19:]
}
