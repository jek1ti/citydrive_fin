package repository

import (
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jekiti/citydrive/admin/internal/config"
	"github.com/jekiti/citydrive/admin/internal/domain"
)

type DBRepository interface {
	GetCarHistory(ctx context.Context, carID string, from, to int64) ([]domain.CarState, error)
	GetCarsHistory(ctx context.Context, from, to int64, activated *bool) (map[string][]domain.CarHistoryPoint, error)
	Close() error
}

type PostgresRepository struct {
	db     *sql.DB
	log    *slog.Logger
	config *config.PostgresConfig
}

func NewPostgresRepository(cfg *config.PostgresConfig, log *slog.Logger) (DBRepository, error) {
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

func (r *PostgresRepository) GetCarHistory(ctx context.Context, carID string, from, to int64) ([]domain.CarState, error) {
	log := r.log.With("module", "repository", "function", "GetCarHistory", "car_id", carID)

	exists, err := r.carExists(ctx, carID)
	if err != nil {
		log.Error("error checking if car exists", "error", err)
		return nil, err
	}
	if !exists {
		log.Info("car does not exist", "car_id", carID)
		return nil, domain.ErrCarNotFound
	}

	query := `
			SELECT
				t.lat, t.lon, t.fuel, t.speed, t.engine_on, t.locked,
				t.activated, t.rpm, t.handbrake, t."timestamp"
			FROM citydrive.car_telemetry_history AS t
			WHERE t.car_id = $1
				AND t."timestamp" >= $2
				AND t."timestamp" <= $3
			ORDER BY t."timestamp" ASC
			`
	log.Info("querying car history from postgres", "car_id", carID, "from", from, "to", to)
	rows, err := r.db.QueryContext(ctx, query, carID, from, to)
	if err != nil {
		log.Error("error querying car history", "error", err)
		return nil, err
	}

	defer rows.Close()
	var history []domain.CarState
	for rows.Next() {
		var state domain.CarState
		err := rows.Scan(
			&state.Lat,
			&state.Lon,
			&state.Fuel,
			&state.Speed,
			&state.EngineOn,
			&state.Locked,
			&state.Activated,
			&state.RPM,
			&state.Handbrake,
			&state.Time,
		)
		if err != nil {
			log.Error("error scanning car history row", "error", err)
			return nil, err
		}
		history = append(history, state)
	}
	return history, nil
}

func (r *PostgresRepository) GetCarsHistory(ctx context.Context, from, to int64, activated *bool) (map[string][]domain.CarHistoryPoint, error) {
	log := r.log.With("module", "repository", "function", "GetCarsHistory")
	query := `
		SELECT
			t.car_id, c.brand, c.model,
			t.lat, t.lon, t.speed, t."timestamp"
		FROM citydrive.car_telemetry_history AS t
		JOIN citydrive.cars AS c ON c.id = t.car_id
		WHERE t."timestamp" >= $1
			AND t."timestamp" <= $2
			AND t.activated = $3
		ORDER BY t.car_id, t."timestamp" ASC
		`		

	rows, err := r.db.QueryContext(ctx, query, from, to, *activated)
	if err != nil {
		log.Error("error querying cars history", "error", err)
		return nil, err
	}
	defer rows.Close()
	history := make(map[string][]domain.CarHistoryPoint)
	for rows.Next() {
		var point domain.CarHistoryPoint
		var carID string
		err := rows.Scan(
			&carID,
			&point.Brand,
			&point.Model,
			&point.Lat,
			&point.Lon,
			&point.Speed,
			&point.Time,
		)
		if err != nil {
			log.Error("error scanning cars history row", "error", err)
			return nil, err
		}
		history[carID] = append(history[carID], point)
	}
	return history, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) carExists(ctx context.Context, carID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM cars WHERE id = $1)`
	err := r.db.QueryRowContext(ctx, query, carID).Scan(&exists)
	return exists, err
}
