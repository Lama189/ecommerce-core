package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Lama189/ecommerce-core/user-service/internal/domain"
	"github.com/Lama189/ecommerce-core/user-service/internal/service/user"
	"github.com/google/uuid"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid reqiest body")
		return
	}

	u, err := h.service.Create(r.Context(), user.CreateUserDTO{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConflict):
			sendError(w, http.StatusConflict, "user with this phone already exists")

		case errors.Is(err, domain.ErrInvalidInput):
			sendError(w, http.StatusBadRequest, err.Error())

		default:
			sendError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	sendJSON(w, http.StatusCreated, userResponse{
		ID:        u.ID,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid reqiest body")
		return
	}

	res, err := h.service.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			sendError(w, http.StatusUnauthorized, "invalid phone or password")
			return
		}
		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	sendJSON(w, http.StatusOK, authResponse{
		User: userResponse{
			ID:        res.User.ID,
			Phone:     res.User.Phone,
			CreatedAt: res.User.CreatedAt,
		},
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	})
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		sendError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "user not found")
			return
		}
		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	sendJSON(w, http.StatusOK, userResponse{
		ID:        u.ID,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
	})
}
