package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/HanmaDevin/calculator/api"
	"github.com/HanmaDevin/calculator/database"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	router := http.NewServeMux()
	api.Handle(router)

	var port *int = flag.Int("port", 8080, "portnumber to use for the server")
	flag.Parse()

	database.StartDbConnection()

	fmt.Printf("Starting Server at localhost:%v\n", *port)

	err := http.ListenAndServe("localhost:"+strconv.Itoa(*port), router)
	if err != nil {
		log.Error(err.Error())
	}
}
