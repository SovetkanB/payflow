package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	data := `"service":"payflow"}`
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"message": "not implemented",
	}

	writeJSON(w, http.StatusCreated, data)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data := map[string]string{
		"message": "not implemented",
		"id":      id,
	}

	writeJSON(w, http.StatusOK, data)
}
