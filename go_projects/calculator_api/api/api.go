package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my calculator API :)")
}

func add(w http.ResponseWriter, r *http.Request) {
	param1 := r.URL.Query().Get("a")
	param2 := r.URL.Query().Get("b")

	if param1 == "" || param2 == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing required query parameters 'a' and/or 'b'")
		return
	}

	a, err := strconv.ParseInt(param1, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "'a' must be a valid integer")
		return
	}

	b, err := strconv.ParseInt(param2, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "'b' must be a valid integer")
		return
	}

	sum := a + b
	fmt.Fprintf(w, "The sum of %v + %v = %v", a, b, sum)
}

func subtract(w http.ResponseWriter, r *http.Request) {
	param1 := r.URL.Query().Get("a")
	param2 := r.URL.Query().Get("b")

	if param1 == "" || param2 == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing required query parameters 'a' and/or 'b'")
		return
	}

	a, err := strconv.ParseInt(param1, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "'a' must be a valid integer")
		return
	}

	b, err := strconv.ParseInt(param2, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "'b' must be a valid integer")
		return
	}

	sum := a - b
	fmt.Fprintf(w, "The result of %v - %v = %v", a, b, sum)
}

func Handle(router *http.ServeMux) {
	router.HandleFunc("GET /", hello)
	router.HandleFunc("GET /add", add)
	router.HandleFunc("GET /subtract", subtract)
}
