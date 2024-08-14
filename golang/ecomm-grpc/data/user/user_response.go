package user

import (

	"time"
)
type UserResponse struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	IsAdmin   bool       `json:"is_admin"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}


type LoginUserResponse struct {
	TokenType string `json:"token_type"`
	Token     string `json:"token"`
}


