package handler

import (
	"errors"
	"net/http"

	mw "github.com/SovetkanB/payflow/user-service/internal/middleware"
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

	res, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	claims := mw.GetUserIDFromContext(r.Context())
	if claims == nil {
		writeJSON(w, http.StatusUnauthorized, "authorization is required")
		return
	}
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			writeJSON(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, model.ErrInvalidPassword):
			writeJSON(w, http.StatusUnauthorized, err.Error())
			return
		default:
			writeJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := mw.GetUserIDFromContext(r.Context())
	if claims == nil {
		writeJSON(w, http.StatusUnauthorized, "authorization is required")
		return
	}

	res, err := h.service.GetUser(r.Context(), claims.UserID)
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
