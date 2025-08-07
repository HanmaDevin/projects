package chat

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/HanmaDevin/schlama/config"
	"github.com/HanmaDevin/schlama/ollama"
	"github.com/HanmaDevin/schlama/styles"
)

//go:embed views/*.html
var views embed.FS
var t, _ = template.New("").ParseFS(views, "views/*.html")

type Response struct {
	Prompt string
	Resp   string
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
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

func Start() {
	main()
}
