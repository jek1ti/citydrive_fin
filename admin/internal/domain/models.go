package domain

type CarShort struct {
	ID    string  `json:"id" db:"id" redis:"id"`
	Brand string  `json:"brand" db:"brand" redis:"brand"`
	Model string  `json:"model" db:"model" redis:"model"`
	Lat   float64 `json:"lat" db:"lat" redis:"lat"`
	Lon   float64 `json:"lon" db:"lon" redis:"lon"`
	Speed int32   `json:"speed" db:"speed" redis:"speed"`
}

type CarDetails struct {
	Brand             string  `json:"brand" db:"brand" redis:"brand"`
	Model             string  `json:"model" db:"model" redis:"model"`
	YearOfManufacture int32   `json:"year_of_manufacture" db:"year_of_manufacture" redis:"year_of_manufacture"`
	Odo               int64   `json:"odo" db:"odo" redis:"odo"`
	Lat               float64 `json:"lat" db:"lat" redis:"lat"`
	Lon               float64 `json:"lon" db:"lon" redis:"lon"`
	Fuel              float64 `json:"fuel" db:"fuel" redis:"fuel"`
	FuelType          string  `json:"fuel_type" db:"fuel_type" redis:"fuel_type"`
	Speed             int32   `json:"speed" db:"speed" redis:"speed"`
	EngineOn          bool    `json:"engine_on" db:"engine_on" redis:"engine_on"`
	Locked            bool    `json:"locked" db:"locked" redis:"locked"`
	Activated         bool    `json:"activated" db:"activated" redis:"activated"`
	RPM               int32   `json:"rpm" db:"rpm" redis:"rpm"`
	Handbrake         bool    `json:"handbrake" db:"handbrake" redis:"handbrake"`
}

type CarHistoryPoint struct {
	Brand string  `json:"brand" db:"brand"`
	Model string  `json:"model" db:"model"`
	Lat   float64 `json:"lat" db:"lat"`
	Lon   float64 `json:"lon" db:"lon"`
	Speed int32   `json:"speed" db:"speed"`
	Time  int64   `json:"time" db:"timestamp"`
}

type CarState struct {
	Lat       float64 `json:"lat" db:"lat"`
	Lon       float64 `json:"lon" db:"lon"`
	Fuel      float64 `json:"fuel" db:"fuel"`
	Speed     int32   `json:"speed" db:"speed"`
	EngineOn  bool    `json:"engine_on" db:"engine_on"`
	Locked    bool    `json:"locked" db:"locked"`
	Activated bool    `json:"activated" db:"activated"`
	RPM       int32   `json:"rpm" db:"rpm"`
	Handbrake bool    `json:"handbrake" db:"handbrake"`
	Time      int64   `json:"time" db:"timestamp"`
}

type HistoryFilter struct {
	From      int64
	To        int64
	Activated *bool
	CarID     string
}
