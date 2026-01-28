package domain

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type CarTelemetry struct {
	Brand             string  `json:"brand"`
	Model             string  `json:"model"`
	YearOfManufacture int32   `json:"year_of_manufacture"`
	Odo               int64   `json:"odo"`
	Lat               float64 `json:"lat"`
	Lon               float64 `json:"lon"`
	Fuel              float64 `json:"fuel"`
	FuelType          string  `json:"fuel_type"`
	Speed             int32   `json:"speed"`
	EngineOn          bool    `json:"engine_on"`
	Locked            bool    `json:"locked"`
	Activated         bool    `json:"activated"`
	Rpm               int32   `json:"rpm"`
	Handbrake         bool    `json:"handbrake"`
	CarID             string  `json:"car_id"`
	ReceivedAt        int64   `json:"received_at"`
}
