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
	"regexp"
	"strings"
	"time"

	"github.com/HanmaDevin/schlama/styles"
	"github.com/charmbracelet/glamour"
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

func GetResponse(ollama *OllamaModel) (string, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(ollama); err != nil {
		return "", fmt.Errorf(styles.ErrorStyle("failed to encode request: %w"), err)
	}

	c := http.Client{Timeout: time.Minute * 3}
	resp, err := c.Post(ollama_api, "application/json", body)
	if err != nil {
		return "", fmt.Errorf("post request to ollama api failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama api returned status %d: %s", resp.StatusCode, string(b))
	}

	var ai Response
	if err := json.NewDecoder(resp.Body).Decode(&ai); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	if ai.Resp == "" {
		return "", fmt.Errorf("ollama api returned empty response")
	}

	return clean(ai.Resp), nil
}

func PullModel(model string) error {
	models := ListModels()
	reg := regexp.MustCompile(`:\w+`)
	modelname := reg.ReplaceAllString(model, "")
	for _, m := range models {
		if m.Name == modelname {
			cmd := exec.Command("ollama", "pull", model)
			msg := fmt.Sprintf("Pulling %s...", model)
			fmt.Println(styles.OutputStyle(msg))
			fmt.Println(styles.HintStyle("Could take a while depending on the model size."))
			if err := cmd.Run(); err != nil {
				log.Fatalf("<-!-- Could not pull model: %s --->\n", model)
			}
			fmt.Println(styles.FinishedStyle("Finished!"))
			return nil
		}
	}
	return fmt.Errorf("model %s not found in the list of available models", model)
}

func IsOllamaRunning() bool {
	resp, err := http.Get("http://localhost:11434")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func ListLocalModels() []string {
	cmd := exec.Command("ollama", "list")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("<-!-- Could not run 'ollama list' --->\n")
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		fmt.Println(styles.OutputStyle("No models found."))
		return []string{}
	}

	var rows []string
	rows = append(rows, styles.HeaderStyle(lines[0]))
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		rows = append(rows, styles.RowStyle(line))
	}

	table := styles.TableBorder(strings.Join(rows, "\n"))
	fmt.Println(table)
	return rows
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

func PrintMarkdown(md string) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Fprintln(os.Stdout, md)
		return
	}
	out, err := r.Render(md)
	if err != nil {
		fmt.Fprintln(os.Stdout, md)
		return
	}
	fmt.Fprint(os.Stdout, out)
}

func clean(s string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleaned := re.ReplaceAllString(s, "")

	// Basic HTML to Markdown replacements
	replacements := []struct {
		old string
		new string
	}{
		{"<b>", "**"}, {"</b>", "**"},
		{"<strong>", "**"}, {"</strong>", "**"},
		{"<i>", "_"}, {"</i>", "_"},
		{"<em>", "_"}, {"</em>", "_"},
		{"<code>", "`"}, {"</code>", "`"},
		{"<pre>", "```\n"}, {"</pre>", "\n```"},
	}

	for _, r := range replacements {
		cleaned = strings.ReplaceAll(cleaned, r.old, r.new)
	}

	// Unescape HTML entities
	cleaned = html.UnescapeString(cleaned)

	return strings.TrimSpace(cleaned)
}
