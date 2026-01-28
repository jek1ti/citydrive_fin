package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jekiti/citydrive/api-gateway/internal/common"
	"github.com/jekiti/citydrive/api-gateway/internal/model"
	"github.com/jekiti/citydrive/api-gateway/internal/service"
	adminpb "github.com/jekiti/citydrive/gen/proto/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminHandler struct {
	adminClient *service.AdminClient
}

func NewAdminHandler(adminClient *service.AdminClient) *AdminHandler {
	return &AdminHandler{adminClient: adminClient}
}

func (h *AdminHandler) GetCar(c *gin.Context) {
	carID := c.Param("id")
	if carID == "" {
		common.Response(c, 400, "INVALID_DATA", "CarID is required", "")
		return
	}
	req := &adminpb.GetCarRequest{Id: carID}

	traceID := common.GetTraceID(c)

	ctx := c.Request.Context()
	respGrpc, err := h.adminClient.GetCar(ctx, traceID, req)
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Admin service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid Admin data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		case codes.NotFound:
			common.Response(c, 404, "CAR_NOT_FOUND", "Car not found", err.Error())
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}
	resp := model.GetCarResponse{
		Car: model.CarDetails{
			Brand:             respGrpc.Car.Brand,
			Model:             respGrpc.Car.Model,
			YearOfManufacture: respGrpc.Car.YearOfManufacture,
			Odo:               respGrpc.Car.Odo,
			Lat:               respGrpc.Car.Lat,
			Lon:               respGrpc.Car.Lon,
			Fuel:              respGrpc.Car.Fuel,
			FuelType:          fuelTypeToString(respGrpc.Car.FuelType),
			Speed:             respGrpc.Car.Speed,
			EngineOn:          respGrpc.Car.EngineOn,
			Locked:            respGrpc.Car.Locked,
			Activated:         respGrpc.Car.Activated,
			RPM:               respGrpc.Car.Rpm,
			Handbrake:         respGrpc.Car.Handbrake,
		},
	}

	c.JSON(200, resp)
}

func (h *AdminHandler) GetCarsNow(c *gin.Context) {
	req := &adminpb.GetCarsNowRequest{}

	traceID := common.GetTraceID(c)

	ctx := c.Request.Context()
	respGrpc, err := h.adminClient.GetCarsNow(ctx, traceID, req)
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Admin service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid Admin data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		case codes.NotFound:
			c.JSON(200, model.GetCarsNowResponse{Cars: []model.CarShort{}})
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}
	var cars []model.CarShort

	if respGrpc.Cars == nil {
		respGrpc.Cars = []*adminpb.CarShort{}
	}

	for _, carPb := range respGrpc.Cars {
		car := model.CarShort{
			ID:    carPb.Id,
			Brand: carPb.Brand,
			Model: carPb.Model,
			Lat:   carPb.Lat,
			Lon:   carPb.Lon,
			Speed: carPb.Speed,
		}
		cars = append(cars, car)
	}

	resp := model.GetCarsNowResponse{
		Cars: cars,
	}

	c.JSON(200, resp)
}

func (h *AdminHandler) GetCarsHistory(c *gin.Context) {
	fromString := c.Query("from")
	if fromString == "" {
		common.Response(c, 400, "INVALID_DATA", "Query parameter FROM is required", "")
		return
	}
	fromInt64, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter From is invalid or nil", err.Error())
		return
	}

	toString := c.Query("to")
	if toString == "" {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter TO is required", "")
		return
	}

	activatedString := c.Query("activated")
	if activatedString == "" {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter ACTIVATED is required", "")
		return
	}
	activatedBool, err := strconv.ParseBool(activatedString)
	if err != nil {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter ACTIVATED is invalid or nil", err.Error())
		return
	}

	toInt64, err := strconv.ParseInt(toString, 10, 64)
	if err != nil {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter To is invalid or nil", err.Error())
		return
	}

	if fromInt64 >= toInt64 {
		common.Response(c, 400, "INVALID_DATA", "Query parameter FROM >= TO", "")
		return
	}

	req := &adminpb.GetCarsHistoryRequest{
		From:      fromInt64,
		To:        toInt64,
		Activated: &activatedBool,
	}

	traceID := common.GetTraceID(c)

	ctx := c.Request.Context()
	respGrpc, err := h.adminClient.GetCarsHistory(ctx, traceID, req)
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Admin service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid Admin data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		case codes.NotFound:
			c.JSON(200, model.GetCarsHistoryResponse{
				HistoryByCar: map[string][]model.CarHistoryPoint{},
			})
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}
	historyByCar := make(map[string][]model.CarHistoryPoint)

	for carID, carHistoryList := range respGrpc.HistoryByCar {
		carHistory := make([]model.CarHistoryPoint, len(carHistoryList.Items))
		for i, historyPoint := range carHistoryList.Items {
			carHistory[i] = model.CarHistoryPoint{
				Brand: historyPoint.Brand,
				Model: historyPoint.Model,
				Lat:   historyPoint.Lat,
				Lon:   historyPoint.Lon,
				Speed: historyPoint.Speed,
				Time:  historyPoint.Time,
			}
		}
		historyByCar[carID] = carHistory
	}

	resp := model.GetCarsHistoryResponse{
		HistoryByCar: historyByCar,
	}

	c.JSON(200, resp)
}

func (h *AdminHandler) GetCarHistory(c *gin.Context) {
	fromString := c.Query("from")
	if fromString == "" {
		common.Response(c, 400, "INVALID_DATA", "Query parameter FROM is required", "")
		return
	}
	fromInt64, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter From is invalid or nil", err.Error())
		return
	}
	toString := c.Query("to")
	if toString == "" {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter TO is required", "")
		return
	}
	toInt64, err := strconv.ParseInt(toString, 10, 64)
	if err != nil {
		common.Response(c, 400, "INVALID_DATA", "Query Parameter To is invalid or nil", err.Error())
		return
	}

	if fromInt64 >= toInt64 {
		common.Response(c, 400, "INVALID_DATA", "Query parameter FROM >= TO", "")
		return
	}

	id := c.Param("id")
	if id == "" {
		common.Response(c, 400, "INVALID_DATA", "Query parameter id is nil", "")
		return
	}

	traceID := common.GetTraceID(c)

	req := &adminpb.GetCarHistoryRequest{
		From: fromInt64,
		To:   toInt64,
		Id:   id,
	}

	ctx := c.Request.Context()
	respGrpc, err := h.adminClient.GetCarHistory(ctx, traceID, req)
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Admin service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid Admin data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		case codes.NotFound:
			common.Response(c, 404, "CAR_NOT_FOUND", "Car not found", err.Error())
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}
	states := make([]model.CarState, len(respGrpc.States))

	for i, stateGrpc := range respGrpc.States {
		states[i] = model.CarState{
			Lat:       stateGrpc.Lat,
			Lon:       stateGrpc.Lon,
			Fuel:      stateGrpc.Fuel,
			Speed:     stateGrpc.Speed,
			EngineOn:  stateGrpc.EngineOn,
			Locked:    stateGrpc.Locked,
			Activated: stateGrpc.Activated,
			RPM:       stateGrpc.Rpm,
			Handbrake: stateGrpc.Handbrake,
			Time:      stateGrpc.Time,
		}
	}

	resp := model.GetCarHistoryResponse{
		States: states,
	}

	c.JSON(200, resp)
}

func fuelTypeToString(fuelType adminpb.FuelType) string {
	switch fuelType {
	case adminpb.FuelType_DIESEL:
		return "diesel"
	case adminpb.FuelType_GASOLINE_92:
		return "92"
	case adminpb.FuelType_GASOLINE_95:
		return "95"
	case adminpb.FuelType_GASOLINE_98:
		return "98"
	default:
		return "unknown"
	}
}
