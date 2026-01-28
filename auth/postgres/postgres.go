package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jekiti/citydrive/auth/internal/config"
)

type Postgres struct {
	master     *pgxpool.Pool
	stopsignal chan struct{}
}

func NewPostgres(cfg config.DatabaseConfig) (*Postgres, error) {
	masterCfg := &PostgresConfig{
		Host:     cfg.Host,
		Database: cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
		MaxConn:  cfg.MaxConn,
	}
	master, err := initPool(masterCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating masterCfg: %w", err)
	}

	return &Postgres{
		master:     master,
		stopsignal: make(chan struct{}),
	}, nil
}

func initPool(cfg *PostgresConfig) (*pgxpool.Pool, error) {

	if cfg.MaxConn < 1 {
		cfg.MaxConn = defaultConnect
	}

	if cfg.Host == "" {
		return nil, errors.New("host can't be empty")
	}

	if cfg.Database == "" {
		return nil, errors.New("database can't be empty")
	}

	if cfg.User == "" {
		return nil, errors.New("user can't be empty")
	}

	if cfg.Password == "" {
		return nil, errors.New("password can't be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?pool_max_conns=%d&pool_max_conn_idle_time=3s",
		cfg.User, cfg.Password, cfg.Host, cfg.Database, cfg.MaxConn)

	pgxpoolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing config for DB: %w", err)
	}

	return createPool(ctx, pgxpoolConfig)

}

func createPool(ctx context.Context, cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("cannot ping DB: %w", err)
	}

	return pool, nil
}

func (postgres *Postgres) Master() *pgxpool.Pool {
	return postgres.master
}

func (postgres *Postgres) Close() {
	postgres.master.Close()
}
