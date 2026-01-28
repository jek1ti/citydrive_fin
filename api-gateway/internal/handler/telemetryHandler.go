package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jekiti/citydrive/api-gateway/internal/common"
	"github.com/jekiti/citydrive/api-gateway/internal/model"
	"github.com/jekiti/citydrive/api-gateway/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TelemetryHandler struct {
	telemetryClient *service.TelemetryClient
}

func NewTelemetryHandler(telemetryClient *service.TelemetryClient) *TelemetryHandler {
	return &TelemetryHandler{telemetryClient: telemetryClient}
}

func (h *TelemetryHandler) PutCarInfo(c *gin.Context) {
	var req model.CarInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, 400, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}
	carID, exists := c.Get("car_id")
	if !exists {
		common.Response(c, 401, "CAR_ID_MISSING", "car_id not found in context", "")
		return
	}

	traceID := common.GetTraceID(c)
	carData := &model.CarData{
		Brand:             req.Brand,
		Model:             req.Model,
		YearOfManufacture: int32(req.YearOfManufacture),
		Odo:               req.Odo,
		Lat:               req.Lat,
		Lon:               req.Lon,
		Fuel:              req.Fuel,
		FuelType:          req.FuelType,
		Speed:             int32(req.Speed),
		EngineOn:          req.EngineOn,
		Locked:            req.Locked,
		Activated:         req.Activated,
		Rpm:               int32(req.Rpm),
		Handbrake:         req.Handbrake,
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.telemetryClient.PutTelemetry(ctx, traceID, carID.(string), carData)
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Telemetry service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid telemetry data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}

	c.JSON(202, gin.H{
		"status":   "accepted",
		"message":  resp.Message,
		"car_id":   carID,
		"trace_id": traceID,
	})

}
