package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jekiti/citydrive/admin/internal/domain"
	"github.com/jekiti/citydrive/admin/internal/service"
	adminpb "github.com/jekiti/citydrive/gen/proto/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	adminpb.UnimplementedAdminServiceServer
	service service.Service
	log     *slog.Logger
}

func NewHandler(service service.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) GetCarsNow(ctx context.Context, req *adminpb.GetCarsNowRequest) (*adminpb.GetCarsNowResponse, error) {
	log := h.log.With("module", "Handler", "function", "GetActiveCars")
	log.Info("received GetActiveCars request")
	cars, err := h.service.GetCarsNow(ctx)
	if err != nil {
		log.Error("error fetching active cars", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	var resp adminpb.GetCarsNowResponse
	for _, car := range cars {
		resp.Cars = append(resp.Cars, &adminpb.CarShort{
			Id:    car.ID,
			Brand: car.Brand,
			Model: car.Model,
			Lat:   car.Lat,
			Lon:   car.Lon,
			Speed: car.Speed,
		})
	}
	log.Info("succesfully fetched active cars", "count", len(resp.Cars))
	return &resp, nil
}

func (h *Handler) GetCar(ctx context.Context, req *adminpb.GetCarRequest) (*adminpb.GetCarResponse, error) {
	log := h.log.With("module", "Handler", "function", "GetCar", "car_id", req.Id)
	log.Info("received GetCar request")
	car, err := h.service.GetCar(ctx, req.Id)
	log.Info("id", "id", req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrCarNotFound) {
			log.Info("car not found", "car_id", req.Id)
			return nil, status.Error(codes.NotFound, "car not found")
		} else if errors.Is(err, domain.ErrInvalidCarID) {
			log.Info("invalid car id", "car_id", req.Id)
			return nil, status.Error(codes.InvalidArgument, "invalid car id")
		} else {
			log.Error("error fetching car details", "error", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	resp := &adminpb.GetCarResponse{
		Car: &adminpb.CarDetails{
			Brand:             car.Brand,
			Model:             car.Model,
			YearOfManufacture: car.YearOfManufacture,
			Lat:               car.Lat,
			Odo:               car.Odo,
			Lon:               car.Lon,
			Fuel:              car.Fuel,
			FuelType:          fuelTypeToProto(car.FuelType),
			Speed:             car.Speed,
			EngineOn:          car.EngineOn,
			Locked:            car.Locked,
			Activated:         car.Activated,
			Rpm:               car.RPM,
			Handbrake:         car.Handbrake,
		},
	}
	log.Info("succesfully fetched car details", "car_id", req.Id)
	return resp, nil
}

func (h *Handler) GetCarHistory(ctx context.Context, req *adminpb.GetCarHistoryRequest) (*adminpb.GetCarHistoryResponse, error) {
	log := h.log.With("module", "handler", "function", "GetCarHistory", "car_id", req.Id)
	log.Info("received GetCarHistory request", "car_id", req.Id, "from", req.From, "to", req.To)
	history, err := h.service.GetCarHistory(ctx, req.Id, req.From, req.To)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTimeRange) {
			return nil, status.Error(codes.InvalidArgument, "invalid time range: 'from' timestamp is greater than or equal to 'to' timestamp")
		} else if errors.Is(err, domain.ErrCarNotFound) {
			return nil, status.Error(codes.NotFound, "car not found")
		} else if errors.Is(err, domain.ErrInvalidCarID) {
			return nil, status.Error(codes.InvalidArgument, "invalid car id")
		} else {
			log.Error("error fetching car history", "error", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	var resp adminpb.GetCarHistoryResponse
	for _, state := range history {
		resp.States = append(resp.States, &adminpb.CarState{
			Lat:       state.Lat,
			Lon:       state.Lon,
			Fuel:      state.Fuel,
			Speed:     state.Speed,
			EngineOn:  state.EngineOn,
			Locked:    state.Locked,
			Activated: state.Activated,
			Rpm:       state.RPM,
			Handbrake: state.Handbrake,
			Time:      state.Time,
		})
	}
	log.Info("successfully fetched car history", "states_count", len(resp.States))
	return &resp, nil
}

func (h *Handler) GetCarsHistory(ctx context.Context, req *adminpb.GetCarsHistoryRequest) (*adminpb.GetCarsHistoryResponse, error) {
	log := h.log.With("module", "handler", "function", "GetCarsHistory")
	log.Info("received GetCarsHistory request", "from", req.From, "to", req.To, "activated", req.Activated)
	history, err := h.service.GetCarsHistory(ctx, req.From, req.To, req.Activated)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTimeRange) {
			log.Info("invalid time range: 'from' timestamp is greater than or equal to 'to' timestamp", "from", req.From, "to", req.To)
			return nil, status.Error(codes.InvalidArgument, "invalid time range: 'from' timestamp is greater than or equal to 'to' timestamp")
		}
		log.Info("error fetching cars history in requested time limits", "from", req.From, "to", req.To)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	historyMap := make(map[string]*adminpb.CarHistoryList)
	var historyList []*adminpb.CarHistoryPoint
	for carID, points := range history {
		for _, point := range points {
			items := adminpb.CarHistoryPoint{
				Brand: point.Brand,
				Model: point.Model,
				Lat:   point.Lat,
				Lon:   point.Lon,
				Speed: point.Speed,
				Time:  point.Time,
			}
			historyList = append(historyList, &items)
		}
		historyMap[carID] = &adminpb.CarHistoryList{
			Items: historyList}
		historyList = []*adminpb.CarHistoryPoint{}
	}
	log.Info("successfully fetched cars history", "cars_count", len(historyMap))
	return &adminpb.GetCarsHistoryResponse{HistoryByCar: historyMap}, nil
}

func fuelTypeToProto(fuelType string) adminpb.FuelType {
	switch fuelType {
	case "diesel":
		return adminpb.FuelType_DIESEL
	case "92":
		return adminpb.FuelType_GASOLINE_92
	case "95":
		return adminpb.FuelType_GASOLINE_95
	case "98":
		return adminpb.FuelType_GASOLINE_98
	default:
		return adminpb.FuelType_FUEL_TYPE_UNSPECIFIED
	}
}
