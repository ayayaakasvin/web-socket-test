package models

import "time"

type User struct {
	ID           uint   `json:"user_id" example:"123"`
	Username     string `json:"username" example:"alice"`
	PasswordHash string
	Role         Role
	CreatedAt    time.Time `json:"created_at" example:"2025-11-14 10:23:45"`
}

type Role string

// Roles
const (
	AdminRole  Role = "admin"
	ClientRole Role = "client"
)
