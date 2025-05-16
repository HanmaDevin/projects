package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/HanmaDevin/example_coins/internal/handler"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func print_banner(filename string) {
	banner, err := os.ReadFile(filename)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(banner))
}

func main() {
	log.SetReportCaller(true)
	var router *chi.Mux = chi.NewRouter()
	handler.Handler(router)

	fmt.Println("Starting GO API service...")
	fmt.Println("Listening on port 8080")
	print_banner("./banner.txt")

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Error(err)
	}

}
