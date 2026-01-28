package authservice

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jekiti/citydrive/auth/internal/models"
	authrepository "github.com/jekiti/citydrive/auth/internal/repository"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)

}

type userService struct {
	repo         authrepository.UserRepository
	log          *slog.Logger
	jwtSecretKey string
}

func NewUserService(repo authrepository.UserRepository, log *slog.Logger, jwtSecretKey string) UserService {
	return &userService{repo: repo, log: log, jwtSecretKey: jwtSecretKey}
}

func (s *userService) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	op := "auth.user_service.Register"
	log := s.log.With("op", op)
	var user models.User
	password, err := password.Generate(15, 3, 2, false, false)
	if err != nil {
		log.Error("error generating password:", slog.Any("error", err))
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error hashing password:", slog.Any("error", err))
		return nil, err
	}
	passwordHashString := string(passwordHash)

	user = models.User{
		Email:        req.Email,
		Name:         req.Name,
		Surname:      req.Surname,
		Department:   req.Department,
		PasswordHash: passwordHashString, 
	}

	err = s.repo.Create(ctx, &user)
	if err != nil {
		log.Error("error creating user:", slog.Any("error", err))
		return nil, err
	}

	return &models.RegisterResponse{
		Password: password,
		ID:       user.ID,
	}, nil

}

func (s *userService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	op := "auth.user_service.Login"
	log := s.log.With("op", op)

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error("error fetching user by email:", slog.Any("error", err))
		return nil, err
	}

	if user == nil {
		log.Warn("user not found")
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Error("invalid password:", slog.Any("error", err))
		return nil, fmt.Errorf("invalid credentials")
	}
	claims := jwt.MapClaims{
		"iss":   "auth.citydrive",
		"sub":   fmt.Sprintf("user:%d", user.ID),
		"email": user.Email,
		"roles": []string{"admin"},
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		log.Error("error signing token:", slog.Any("error", err))
		return nil, err
	}
	return &models.LoginResponse{AccessToken: tokenString}, nil
}
