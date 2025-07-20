package handler

import (
	"net/http"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/mapping"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/middleware"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/response"
)

type AuthHandler struct {
	register usecase.RegisterUserUseCase
	login    usecase.LoginUserUseCase
}

func NewAuthHandler(reg usecase.RegisterUserUseCase, login usecase.LoginUserUseCase) *AuthHandler {
	return &AuthHandler{register: reg, login: login}
}

var badJSONError = apperror.New("http.decode", "bad_json", apperror.KindValidation, nil)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req mapping.RegisterRequest
	if err := mapping.DecodeJSONStrict(r, &req); err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), badJSONError)
		return
	}
	out, err := h.register.Execute(r.Context(), usecase.RegisterUserInput{
		Login: req.Login, Password: req.Password,
	})
	if err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), err)
		return
	}
	response.WriteJSON(w, http.StatusCreated, map[string]any{"user": out})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req mapping.LoginRequest
	if err := mapping.DecodeJSONStrict(r, &req); err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), badJSONError)
		return
	}
	out, err := h.login.Execute(r.Context(), usecase.LoginUserInput{
		Login: req.Login, Password: req.Password,
	})
	if err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), err)
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token": out.AccessToken,
		"user":         out.User,
	})
}
