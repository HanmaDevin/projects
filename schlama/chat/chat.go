package chat

import (
	"embed"
	"fmt"
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

func Start() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(styles.HintStyle("Starting chat..."))
		cfg := config.ReadConfig()
		data := data{
			Model:  cfg.Model,
			Prompt: "",
			Resp:   "",
		}

		models, err := getLocalModels()
		if err != nil {
			http.Error(w, "Failed to get local models: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.Options = models

		if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("POST /set-model", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(styles.HintStyle("Setting model..."))
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		cfg := config.Config{
			Model: r.FormValue("model"),
		}
		if err := config.WriteConfig(cfg); err != nil {
			http.Error(w, "Failed to write config: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(styles.HintStyle("Processing chat..."))
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		cfg := config.ReadConfig()

		resp, err := ollama.GetResponse(cfg)
		if err != nil {
			http.Error(w, "Failed to get response from Ollama: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.ExecuteTemplate(w, "response.html", resp); err != nil {
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
