package chat

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
)

//go:embed views/*.html
var views embed.FS
var t, _ = template.New("").ParseFS(views, "views/*.html")

type data struct {
	Model   string   `json:"model"`
	Prompt  string   `json:"prompt"`
	Resp    string   `json:"response"`
	Options []string `json:"options"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(styles.OutputStyle("Starting chat..."))
	cfg := config.ReadConfig()
	data := data{
		Model:  cfg.Model,
		Prompt: "",
		Resp:   "",
	}

	models, err := getLocalModels()
	if err != nil {
		log.Println(styles.ErrorStyle("Failed to get local models: " + err.Error()))
		http.Error(w, "Failed to get local models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data.Options = models

	if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(styles.ErrorStyle("Failed to render index template: " + err.Error()))
		http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
	}
}

func setModelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(styles.OutputStyle("Setting model..."))
	if err := r.ParseForm(); err != nil {
		log.Println(styles.ErrorStyle("Failed to parse form: " + err.Error()))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	cfg := config.Config{
		Model: r.FormValue("model"),
	}
	if err := config.WriteConfig(cfg); err != nil {
		log.Println(styles.ErrorStyle("Failed to write config: " + err.Error()))
		http.Error(w, "Failed to write config: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(styles.OutputStyle("Handling chat request..."))
	if err := r.ParseMultipartForm(200 << 20); err != nil {
		log.Println(styles.ErrorStyle("Failed to parse form: " + err.Error()))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	cfg := config.ReadConfig()

	prompt := r.FormValue("prompt")
	if prompt == "" {
		log.Println(styles.OutputStyle("Prompt cannot be empty"))
		http.Error(w, "Prompt cannot be empty", http.StatusBadRequest)
		return
	}
	cfg.Msg[0].Content = prompt

	form := r.MultipartForm
	for _, fileHeader := range form.File["files"] {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println(styles.ErrorStyle("Failed to open file: " + err.Error()))
			http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			log.Println(styles.ErrorStyle("Failed to read file: " + err.Error()))
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		switch fileHeader.Header.Get("Content-Type") {
		case "text/plain":
			log.Println(styles.OutputStyle("Received text file: " + fileHeader.Filename))
			cfg.Msg[0].Content += "\n" + string(content)
		case "image/png", "image/jpeg", "image/gif":
			log.Println(styles.OutputStyle("Received image file: " + fileHeader.Filename))
			encoded := encodeImageToBase64(content)
			cfg.Msg[0].Images = append(cfg.Msg[0].Images, encoded)
		case "text/html":
			log.Println(styles.OutputStyle("Received HTML file: " + fileHeader.Filename))
			cfg.Msg[0].Content += "\n" + string(content)
		case "application/json":
			log.Println(styles.OutputStyle("Received JSON file: " + fileHeader.Filename))
			cfg.Msg[0].Content += "\n" + string(content)
		case "text/xml", "application/xml":
			cfg.Msg[0].Content += "\n" + string(content)
			log.Println(styles.OutputStyle("Received XML file: " + fileHeader.Filename))
		default:
			log.Println(styles.OutputStyle("Received file with unknown content type: " + fileHeader.Header.Get("Content-Type")))
			cfg.Msg[0].Content += "\n" + string(content)
		}
	}

	resp, err := ollama.GetResponse(cfg)
	if err != nil {
		log.Println(styles.OutputStyle("Failed to get response from Ollama: " + err.Error()))
		http.Error(w, "Failed to get response from Ollama: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := data{
		Model:  cfg.Model,
		Prompt: prompt,
		Resp:   resp,
	}

	if err := t.ExecuteTemplate(w, "response.html", data); err != nil {
		log.Println(styles.OutputStyle("Failed to render response template: " + err.Error()))
		http.Error(w, "Failed to render response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func Start() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", rootHandler)
	router.HandleFunc("POST /set-model", setModelHandler)
	router.HandleFunc("POST /chat", chatHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := openURL("http://localhost:8080")
	if err != nil {
		fmt.Println(styles.ErrorStyle("Failed to open browser: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(styles.HintStyle("Chat started in your browser at http://localhost:8080"))
	server.ListenAndServe()
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Use xdg-open on native Linux environments
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Run()
}

func getLocalModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	var models []string
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}
	return models, nil
}

func encodeImageToBase64(content []byte) string {
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}

func getDirContent(files []*multipart.FileHeader) (string, error) {
	var builder strings.Builder
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return "", err
		}
		content, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("File: %s\n", fh.Filename))
		builder.Write(content)
		builder.WriteString("\n")
	}
	return builder.String(), nil
}
