package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/HanmaDevin/calculator/types"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my calculator API :)")
}

func add(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	if header != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(header, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not 'application/json'"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	var obj types.Object

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&obj)

	if obj.Number1 == nil || obj.Number2 == nil {
		msg := "Request body must contain 'number1' and 'number2' fields"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := "Request body contains badly-formed JSON"
			http.Error(w, msg, http.StatusBadRequest)
			return
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
			return
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
	}

	var res types.Result
	res.Res = obj.Add()
	res.Desc = "Successfully added two numbers"
	json.NewEncoder(w).Encode(res)
}

func subtract(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	if header != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(header, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not 'application/json'"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	var obj types.Object

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&obj)

	if obj.Number1 == nil || obj.Number2 == nil {
		msg := "Request body must contain 'number1' and 'number2' fields"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := "Request body contains badly-formed JSON"
			http.Error(w, msg, http.StatusBadRequest)
			return
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
			return
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
	}

	var res types.Result
	res.Res = obj.Subtract()
	res.Desc = "Successfully subtracted two numbers"
	json.NewEncoder(w).Encode(res)
}

func Handle(router *http.ServeMux) {
	router.HandleFunc("GET /", hello)
	router.HandleFunc("POST /add", add)
	router.HandleFunc("POST /subtract", subtract)
}
