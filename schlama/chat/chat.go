package chat

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/HanmaDevin/schlama/styles"
)

//go:embed views/*.html
var views embed.FS
var t, _ = template.New("").ParseFS(views, "views/*.html")

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := openURL("http://localhost:8080")
	if err != nil {
		fmt.Println(styles.TableBorder(styles.ErrorStyle("Failed to open browser: " + err.Error())))
		os.Exit(1)
	}

	fmt.Println(styles.TableBorder(styles.HintStyle("Chat started in your browser at http://localhost:8080")))
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
