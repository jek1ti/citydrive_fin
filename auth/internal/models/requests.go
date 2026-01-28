package models

type RegisterRequest struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Department string `json:"department"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
