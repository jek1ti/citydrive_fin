package models

type Car struct {
	ID   int64  `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Plate string `json:"plate"`
	Year  int    `json:"year"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}