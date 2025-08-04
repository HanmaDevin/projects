package chat

import (
	"embed"
	"net/http"
	"text/template"
)

//go:embed views/*.html
var views embed.FS
var t, _ = template.New("").ParseFS(views, "views/*.html")

func main() {
	router := http.NewServeMux()
}
