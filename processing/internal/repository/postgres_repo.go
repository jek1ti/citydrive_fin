package repository

import (
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jekiti/citydrive/processing/internal/config"
	"github.com/jekiti/citydrive/processing/internal/domain"
)

type DBRepository interface {
	SaveTelemetry(telemetry domain.CarTelemetry) error
	Close() error
}

type PostgresRepository struct {
	db     *sql.DB
	log    *slog.Logger
	config *config.DBConfig
}

func NewPostgresRepository(cfg *config.DBConfig, log *slog.Logger) (DBRepository, error) {
	log = log.With("module", "repository", "function", "NewPostgresRepository")
	dsn := cfg.URL
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error("failed to ping postgres", "error", err)
		return nil, err
	}
	return &PostgresRepository{
		db:     db,
		log:    log,
		config: cfg,
	}, nil
}

func (r *PostgresRepository) SaveTelemetry(tel domain.CarTelemetry) error {
	log := r.log.With("module", "repository", "function", "SaveCarState", "car_id", tel.CarID)
	query := `
		INSERT INTO citydrive.car_telemetry_history
		(car_id, lat, lon, fuel, speed, engine_on, locked, activated, rpm, handbrake, odo, "timestamp")
		VALUES
		($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`

	_, err := r.db.Exec(query,
		tel.CarID,      
		tel.Lat,
		tel.Lon,
		tel.Fuel,
		tel.Speed,
		tel.EngineOn,
		tel.Locked,
		tel.Activated,
		tel.Rpm,
		tel.Handbrake,
		tel.Odo,
		tel.ReceivedAt,
	)


	if err != nil {
		log.Error("error saving car state to postgres", "error", err)
		return err
	}
	log.Info("car state saved to postgres", "car_id", tel.CarID)
	return nil
}

func (r *PostgresRepository) Close() error {
	log := r.log.With("module", "repository", "function", "Close")
	log.Info("closing postgres connection")
	err := r.db.Close()
	if err != nil {
		log.Error("error closing postgres connection", "error", err)
		return err
	}
	return nil
}
