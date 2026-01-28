package model

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
    Email      string `json:"email" binding:"required,email"`
    Name       string `json:"name" binding:"required"`
    Surname    string `json:"surname" binding:"required"` 
    Department string `json:"department" binding:"required"`
}