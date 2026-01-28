package authrepository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jekiti/citydrive/auth/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	op := "auth.user_repository"
	log := r.log.With("op", op)

	query := `INSERT INTO users (
	email, name, surname, department, password_hash
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Name, user.Surname, user.Department, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		log.Error("error saving user:", slog.Any("error", err))
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	op := "auth.user_repository"
	log := r.log.With("op", op)

	var user models.User

	query := `SELECT id, email, name, surname, department, password_hash, created_at, updated_at
	FROM users
	WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.Department,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("user not found by email", "email", email)
			return nil, nil
		}
		log.Error("error getting user by email", slog.Any("error", err))
		return nil, err
	}

	return &user, nil
}
