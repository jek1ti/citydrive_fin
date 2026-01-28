package handler

import (
	"context"
	"log/slog"

	telemetrypb "github.com/jekiti/citydrive/gen/proto/telemetry"
	"github.com/jekiti/citydrive/telemetry/internal/models"
	"github.com/jekiti/citydrive/telemetry/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TelemetryHandler struct {
	telemetrypb.UnimplementedTelemetryServiceServer
	telemetryService *service.TelemetryService
	log              *slog.Logger
}

func NewTelemetryHandler(telemetryService *service.TelemetryService, log *slog.Logger) *TelemetryHandler {
	return &TelemetryHandler{telemetryService: telemetryService, log: log}
}

func (h *TelemetryHandler) PutTelemetry(ctx context.Context, req *telemetrypb.PutRequest) (*telemetrypb.PutResponse, error) {
	h.log.Info("received PutTelemetry request in")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "car_id required in metadata")
	}
	h.log.Info("incoming metadata", "metadata", md)
	carIDs := md.Get("car_id")
	traceIDs := md.Get("trace_id")
	var traceID string
	if len(traceIDs) > 0 {
		traceID = traceIDs[0]
	}
	if len(carIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "car_id required in metadata")
	}
	carID := carIDs[0]

	ctxNew := context.WithValue(ctx, "trace_id", traceID)
	log := h.log.With(
		"module", "handler",
		"function", "PutTelemetry",
		"car_id", carID,
		"trace_id", traceID,
	)
	log.Info("processing telemetry request")
	data := &models.TelemetryData{
		Brand:             req.Brand,
		Model:             req.Model,
		YearOfManufacture: req.YearOfManufacture,
		Odo:               req.Odo,
		Lat:               req.Lat,
		Lon:               req.Lon,
		Fuel:              req.Fuel,
		FuelType:          req.FuelType,
		Speed:             req.Speed,
		EngineOn:          req.EngineOn,
		Locked:            req.Locked,
		Activated:         req.Activated,
		RPM:               req.Rpm,
		Handbrake:         req.Handbrake,
	}

	err := h.telemetryService.ProcessTelemetry(ctxNew, carID, data)
	if err != nil {
		log.Error("telemetry processing failed", "error", err)
		return nil, status.Error(codes.Internal, "can't put telemetry")
	}
	log.Info("telemetry processed successfully")
	return &telemetrypb.PutResponse{Message: "telemetry processed successfully"}, nil
}
