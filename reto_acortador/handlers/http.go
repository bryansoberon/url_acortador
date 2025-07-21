package handlers

import (
	"encoding/json"
	"net/http"
	"reto_acordator/shortener"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Shortener *shortener.Shortener
	Store     *shortener.Store
}

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !isValidURL(req.LongURL) {
		http.Error(w, "URL larga inválida o malformada", http.StatusBadRequest)
		return
	}

	code, err := h.Shortener.GenerateShortCode(req.LongURL)
	if err != nil {
		http.Error(w, "Falla al generar el código corto", http.StatusInternalServerError)
		return
	}

	if err := h.Store.Save(code, req.LongURL); err != nil {
		http.Error(w, "Short code collision", http.StatusInternalServerError)
		return
	}

	shortURL := "http://localhost:8080/" + code
	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, exists := h.Store.Get(code)
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

// Validación simple para URLs no vacías, que empiecen con http:// o https:// y que tengan dominio
func isValidURL(url string) bool {
	if url == "" {
		return false
	}
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return false
	}
	if len(url) <= len("https://") {
		return false
	}
	return true
}
