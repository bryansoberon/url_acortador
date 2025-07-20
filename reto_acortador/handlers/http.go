package handlers

import (
	"encoding/json"
	"net/http"
	"reto_acordator/shortener"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Shortener *shortener.Shortener
	Store     *shortener.Store
}

// JSON de entrada para POST /shorten
type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

// JSON de salida
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// POST /shorten → genera una URL corta
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.LongURL == "" {
		http.Error(w, "Invalid JSON or empty long_url", http.StatusBadRequest)
		return
	}

	// Generar código corto
	code, err := h.Shortener.GenerateShortCode(req.LongURL)
	if err != nil {
		http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
		return
	}

	// Guardar en memoria
	if err := h.Store.Save(code, req.LongURL); err != nil {
		http.Error(w, "Short code collision", http.StatusInternalServerError)
		return
	}

	shortURL := "http://localhost:8080/" + code
	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
}

// GET /{short_code} → redirige a la URL original
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, exists := h.Store.Get(code)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Redirección 301 (permanente) justificada
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
