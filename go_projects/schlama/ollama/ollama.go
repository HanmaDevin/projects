package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	return &OllamaModel{}
}

func GetResponse(ollama *OllamaModel) string {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(ollama)

	c := http.Client{Timeout: time.Minute * 3}

	resp, err := c.Post(ollama_api, "application/json", body)
	if err != nil {
		log.Fatal("<-!-- Post request to ollama api failed --->\n")
	}

	defer resp.Body.Close()
	var ai Response
	json.NewDecoder(resp.Body).Decode(&ai)

	return clean(ai.Resp)
}

func PullModel(model string) error {
	models := ListModels()
	for _, m := range models {
		if m.Name == model {
			cmd := exec.Command("ollama", "pull", model)
			fmt.Println("Pulling manifest...")
			fmt.Println("Could take a while depending on the model size.")
			if err := cmd.Run(); err != nil {
				log.Fatalf("<-!-- Could not pull model: %s --->\n", model)
			}
			fmt.Println("Finished!")
			return nil
		}
	}
	return fmt.Errorf("model %s not found in the list of available models", model)
}

func ListLocalModels() {
	cmd := exec.Command("ollama", "list")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal("<-!-- Could not run 'ollama list' --->\n")
	}
}

func IsModelPresent(model string) bool {
	cmd := exec.Command("ollama", "list")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("<-!-- Could not run 'ollama list' --->\n")
	}

	table := strings.Split(string(out), "\n")[1:] // Skip the header line

	for _, line := range table {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == model {
			return true
		}
	}
	return false
}

type ModelInfo struct {
	Name  string
	Sizes []string
}

func extractModels(n *html.Node) []ModelInfo {
	var models []ModelInfo

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			hasModelAttr := false
			for _, attr := range n.Attr {
				if attr.Key == "x-test-model" {
					hasModelAttr = true
					break
				}
			}
			if hasModelAttr {
				var modelName string
				var sizes []string

				var findModelName func(*html.Node)
				findModelName = func(node *html.Node) {
					if node.Type == html.ElementNode && node.Data == "div" {
						for _, attr := range node.Attr {
							if attr.Key == "x-test-model-title" {
								for _, a := range node.Attr {
									if a.Key == "title" {
										modelName = a.Val
										break
									}
								}
							}
						}
					}
					for c := node.FirstChild; c != nil; c = c.NextSibling {
						findModelName(c)
					}
				}
				findModelName(n)

				var findSizes func(*html.Node)
				findSizes = func(node *html.Node) {
					if node.Type == html.ElementNode && node.Data == "span" {
						for _, attr := range node.Attr {
							if attr.Key == "x-test-size" {
								if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
									sizes = append(sizes, strings.TrimSpace(node.FirstChild.Data))
								}
							}
						}
					}
					for c := node.FirstChild; c != nil; c = c.NextSibling {
						findSizes(c)
					}
				}
				findSizes(n)

				if modelName != "" && len(sizes) > 0 {
					models = append(models, ModelInfo{
						Name:  modelName,
						Sizes: sizes,
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return models
}

func ListModels() []ModelInfo {
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

	doc, err := html.Parse(bytes.NewReader(b))

	return extractModels(doc)
}

func clean(s string) string {
	return s
}
