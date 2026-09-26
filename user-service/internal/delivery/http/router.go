package http

import (
	"github.com/Lama189/ecommerce-core/user-service/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *UserHandler, jwtManager user.TokenManager) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.Register)
			r.Post("/login", handler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(jwtManager))
			r.Get("/users/me", handler.GetMe)
		})
	})

	return r
}
