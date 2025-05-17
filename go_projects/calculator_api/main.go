package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/HanmaDevin/calculator/api"
)

func main() {
	router := http.NewServeMux()
	api.Handle(router)

	var port *int = flag.Int("port", 8080, "portnumber to use for the server")
	flag.Parse()

	fmt.Printf("Starting Server at localhost:%v\n", *port)

	err := http.ListenAndServe("localhost:"+strconv.Itoa(*port), router)
	if err != nil {
		log.Fatal(err)
	}
}
