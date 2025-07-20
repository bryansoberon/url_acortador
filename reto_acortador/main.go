package main

import (
	"log"
	"net/http"

	"reto_acordator/handlers"
	"reto_acordator/shortener"

	"github.com/go-chi/chi/v5")

func main() {
	// Inicializar almacenamiento y servicios
	store := shortener.NewStore()
	service := shortener.NewShortener(store, 6)
	handler := &handlers.Handler{
		Shortener: service,
		Store:     store,
	}

	r := chi.NewRouter()

	// Rutas
	r.Post("/shorten", handler.ShortenURL)
	r.Get("/{code}", handler.RedirectURL)

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
