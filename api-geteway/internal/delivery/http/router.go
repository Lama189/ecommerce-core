package http

import (
	"net/http"

	"github.com/Lama189/ecommerce-core/api-geteway/internal/infrastructure/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *UserHandler, tokenValidator jwt.TokenValidator) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.Register)
			r.Post("/login", handler.Login)
			r.Post("/refresh", handler.Refresh)
		})

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(tokenValidator))
			r.Get("/users/me", handler.GetMe)
		})
	})

	return r
}
