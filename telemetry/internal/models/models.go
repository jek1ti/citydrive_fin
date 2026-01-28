package models

type TelemetryData struct {
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
	RPM               int32   `json:"rpm"`
	Handbrake         bool    `json:"handbrake"`
}

type Violation struct {
	Type    string
	CarID   string
	Data    TelemetryData
	Details map[string]interface{} // {speed: 120, limit: 110}
}

const (
	ViolationTypeSpeedingLow    = "speeding_low"
	ViolationTypeSpeedingMedium = "speeding_medium"
	ViolationTypeSpeedingHigh   = "speeding_high"
	ViolationTypeLowFuel        = "low_fuel"
	ViolationTypeDrift          = "drift"
	ViolationTypeStealedAuto    = "stealed_auto"
)
