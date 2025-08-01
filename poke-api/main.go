package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// important comment to ensure the embed directive works
//
//go:embed views/*.html
var views embed.FS

const base_url = "https://pokeapi.co/api/v2/pokemon/"

var t, _ = template.New("").ParseFS(views, "views/*.html")

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("POST /fetch", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		}

		resp, err := http.Get(base_url + strings.TrimSpace(strings.ToLower(r.FormValue("pokemon"))))
		if err != nil {
			http.Error(w, "Couldn't fetch shit!!!", http.StatusInternalServerError)
		}

		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse data", http.StatusInternalServerError)
		}

		if err := t.ExecuteTemplate(w, "response.html", data); err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		}
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Server listening on localhost:3000")
	server.ListenAndServe()
}
