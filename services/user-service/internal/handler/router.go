package handler

import (
	"net/http"

	"github.com/SovetkanB/payflow/user-service/internal/config"
	mw "github.com/SovetkanB/payflow/user-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	handler *Handler
	cfg     *config.Config
}

func NewRouter(h *Handler, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", h.Health)

	r.Route("/api/v1", func(r chi.Router) {
		//auth endpoint
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.CreateUser)
			r.Post("/login", h.Login)
		})

		//protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware(cfg.JWTSecret))
			//users endpoints
			r.Route("/users", func(r chi.Router) {
				r.Get("/me", h.Me)
				r.Get("/{id}", h.GetUser)
			})
		})

	})

	return r
}
