package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/SovetkanB/payflow/user-service/internal/model"
	"github.com/SovetkanB/payflow/user-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{service: srv}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok", "service": "user-service"}
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(req)

	res, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			writeJSON(w, http.StatusNotFound, err.Error())
			return
		default:
			writeJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, res)
}
