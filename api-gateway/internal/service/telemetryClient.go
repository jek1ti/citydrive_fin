package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jekiti/citydrive/api-gateway/internal/config"
	"github.com/jekiti/citydrive/api-gateway/internal/model"
	telemetrypb "github.com/jekiti/citydrive/gen/proto/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type TelemetryClient struct {
	client telemetrypb.TelemetryServiceClient
	conn   *grpc.ClientConn
	log    *slog.Logger
}

func NewTelemetryClient(cfg *config.GRPCConfig, log *slog.Logger) (*TelemetryClient, error) {
	log = log.With(
		"module", "telemetry.client",
		"function", "NewTelemetryClient")
	log.Info("creating new telemetry gRPC client", "target", cfg.TelemetryAddr)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.TelemetryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), 
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to telemetry service at %s: %w", cfg.TelemetryAddr, err)
	}

	client := telemetrypb.NewTelemetryServiceClient(conn)
	log.Info("telemetry gRPC client created successfully on", "addr", cfg.TelemetryAddr)
	return &TelemetryClient{client: client, conn: conn, log: log}, nil
}

func (c *TelemetryClient) PutTelemetry(ctx context.Context, traceID string, carID string, data *model.CarData) (*telemetrypb.PutResponse, error) {
	log := c.log.With(
		"module", "telemetry.client",
		"function", "PutTelemetry",
		"car_id", carID,
		"trace_id", traceID,
	)
	log.Info("sending telemetry data to telemetry service from client")
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	md := metadata.Pairs("trace_id", traceID, "car_id", carID)
	ctx = metadata.NewOutgoingContext(ctx, md)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := ValidateCarData(data); err != nil {
		log.Error("invalid car data", "error", err)
		return nil, err
	}

	req := &telemetrypb.PutRequest{
		Brand:             data.Brand,
		Model:             data.Model,
		YearOfManufacture: data.YearOfManufacture,
		Odo:               data.Odo,
		Lat:               data.Lat,
		Lon:               data.Lon,
		Fuel:              data.Fuel,
		FuelType:          data.FuelType,
		Speed:             data.Speed,
		EngineOn:          data.EngineOn,
		Locked:            data.Locked,
		Activated:         data.Activated,
		Rpm:               data.Rpm,
		Handbrake:         data.Handbrake,
	}
	log.Info("request prepared, calling PutTelemetry on gRPC client")

	response, err := c.client.PutTelemetry(ctx, req)
	if err != nil {
		log.Error("error calling PutTelemetry on gRPC client", "error", err)
		return nil, fmt.Errorf("failed to send telemetry: %w", err)
	}
	log.Info("telemetry data sent successfully to telemetry service from client")
	return response, nil
}

func (c *TelemetryClient) Close() error {
	return c.conn.Close()
}

func ValidateCarData(data *model.CarData) error {
	if data == nil {
		return fmt.Errorf("car data cannot be nil")
	}
	if data.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if len(data.Brand) > 100 {
		return fmt.Errorf("brand too long")
	}

	if data.Model == "" {
		return fmt.Errorf("model is required")
	}
	if len(data.Model) > 100 {
		return fmt.Errorf("model too long")
	}

	currentYear := time.Now().Year()
	if data.YearOfManufacture < 1900 || data.YearOfManufacture > int32(currentYear+1) {
		return fmt.Errorf("year of manufacture must be between 1900 and %d", currentYear+1)
	}

	if data.Odo < 0 || data.Odo > 1000000 {
		return fmt.Errorf("odometer must be between 0 and 1,000,000 km")
	}

	if data.Lat < -90 || data.Lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if data.Lon < -180 || data.Lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	if data.Fuel < 0 || data.Fuel > 100 {
		return fmt.Errorf("fuel level must be between 0 and 100 percent")
	}

	validFuelTypes := map[string]bool{
		"diesel": true,
		"92":     true,
		"95":     true,
		"98":     true,
	}
	if !validFuelTypes[data.FuelType] {
		return fmt.Errorf("fuel type must be one of: diesel, 92, 95, 98")
	}

	if data.Speed < 0 || data.Speed > 300 {
		return fmt.Errorf("speed must be between 0 and 300 km/h")
	}

	if data.Rpm < 0 || data.Rpm > 10000 {
		return fmt.Errorf("RPM must be between 0 and 10,000")
	}
	return nil
}
