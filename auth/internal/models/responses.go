package models

type RegisterResponse struct {
	Password string `json:"password"`
	ID       int64  `json:"id"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
