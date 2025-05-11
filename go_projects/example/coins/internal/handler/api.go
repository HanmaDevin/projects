package handler

import (
	"github.com/HanmaDevin/example_coins/internal/middleware"
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(router *chi.Mux) {
	router.Use(chimiddle.StripSlashes)

	router.Route("/account", func(route chi.Router) {
		route.Use(middleware.Authorization)
		route.Get("/coins", GetCoinsBalance)
	})
}
