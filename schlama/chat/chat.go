package chat

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func Start() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Println(styles.OutputStyle("Starting chat..."))
		cfg := config.ReadConfig()
		data := data{
			Model:  cfg.Model,
			Prompt: "",
			Resp:   "",
		}

		models, err := getLocalModels()
		if err != nil {
			log.Println(styles.OutputStyle("Failed to get local models: " + err.Error()))
			http.Error(w, "Failed to get local models: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.Options = models

		if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
			log.Println(styles.OutputStyle("Failed to render index template: " + err.Error()))
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("POST /set-model", func(w http.ResponseWriter, r *http.Request) {
		log.Println(styles.OutputStyle("Setting model..."))
		if err := r.ParseForm(); err != nil {
			log.Println(styles.OutputStyle("Failed to parse form: " + err.Error()))
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		cfg := config.Config{
			Model: r.FormValue("model"),
		}
		if err := config.WriteConfig(cfg); err != nil {
			log.Println(styles.OutputStyle("Failed to write config: " + err.Error()))
			http.Error(w, "Failed to write config: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
		log.Println(styles.OutputStyle("Handling chat request..."))
		if err := r.ParseMultipartForm(200 << 20); err != nil {
			log.Println(styles.OutputStyle("Failed to parse form: " + err.Error()))
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

		file, fileHeader, err := r.FormFile("file")
		if err != nil && err != http.ErrMissingFile {
			return
		}
		if err == nil {
			defer file.Close()
			log.Println(styles.OutputStyle("Received file: " + fileHeader.Filename))
			content, err := io.ReadAll(file)
			if err != nil {
				log.Println(styles.OutputStyle("Failed to read file: " + err.Error()))
				http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			cfg.Msg[0].Content += "\n" + string(content)
		} else if err == http.ErrMissingFile {
			log.Println(styles.OutputStyle("No file uploaded"))
		}

		dir := r.FormValue("dir")
		log.Println(styles.OutputStyle("Received directory: " + dir))
		if dir != "" {
			content, err := getDirContent(dir)
			if err != nil {
				log.Println(styles.OutputStyle("Failed to read directory: " + err.Error()))
				http.Error(w, "Failed to read directory: "+err.Error(), http.StatusInternalServerError)
				return
			}
			cfg.Msg[0].Content += "\n" + content
		}

		images := r.FormValue("images")
		log.Println(styles.OutputStyle("Received image: " + images))
		if images != "" {
			encoded, err := encodeImageToBase64(images)
			if err != nil {
				log.Println(styles.OutputStyle("Failed to read image: " + err.Error()))
				http.Error(w, "Failed to read image: "+err.Error(), http.StatusInternalServerError)
				return
			}
			cfg.Msg[0].Images = append(cfg.Msg[0].Images, encoded)
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
	})

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
			sb.WriteString("File: " + filepath.Base(path) + "\n")
			sb.Write(content)
			sb.WriteString("\n")
		}
		return nil
	})
	return sb.String(), err
}

func encodeImageToBase64(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}
