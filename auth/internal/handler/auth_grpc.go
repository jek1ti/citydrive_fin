package handler

import (
	"context"
	"log/slog"

	"github.com/jekiti/citydrive/auth/internal/models"
	authservice "github.com/jekiti/citydrive/auth/internal/service"
	auth "github.com/jekiti/citydrive/gen/proto/auth"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	service authservice.UserService
	log     *slog.Logger
}

func NewAuthHandler(service authservice.UserService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{service: service, log: logger}
}

func (h *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	op := "auth.handler.Register"
	log := h.log.With("op", op)
	log.Info("Register request received", slog.String("email", req.Email))
	modelReq := &models.RegisterRequest{
		Email:      req.Email,
		Name:       req.Name,
		Surname:    req.Surname,
		Department: req.Department,
	}
	res, err := h.service.Register(ctx, modelReq)
	if err != nil {
		log.Error("error in Register handler:", slog.Any("error", err))
		return nil, err
	}
	return &auth.RegisterResponse{
		Password: res.Password,
		Id:       res.ID,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	op := "auth.handler.Login"
	log := h.log.With("op", op)
	log.Info("Login request received", slog.String("email", req.Email))
	modelReq := &models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	res, err := h.service.Login(ctx, modelReq)
	if err != nil {
		log.Error("error in Login handler:", slog.Any("error", err))
		return nil, err
	}
	return &auth.LoginResponse{
		AccessToken: res.AccessToken,
	}, nil
}
