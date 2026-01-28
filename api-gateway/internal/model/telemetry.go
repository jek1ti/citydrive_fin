package model

type CarData struct {
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
}

type CarInfoRequest struct {
	Brand             string  `json:"brand" binding:"required"`
	Model             string  `json:"model" binding:"required"`
	YearOfManufacture int     `json:"year_of_manufacture"`
	Odo               int64   `json:"odo" binding:"min=0"`
	Lat               float64 `json:"lat" binding:"required"`
	Lon               float64 `json:"lon" binding:"required"`
	Fuel              float64 `json:"fuel" binding:"min=0,max=100"`
	FuelType          string  `json:"fuel_type" binding:"required,oneof=diesel 92 95 98"`
	Speed             int     `json:"speed" binding:"min=0"`
	EngineOn          bool    `json:"engine_on"`
	Locked            bool    `json:"locked"`
	Activated         bool    `json:"activated"`
	Rpm               int     `json:"rpm" binding:"min=0"`
	Handbrake         bool    `json:"handbrake"`
}
