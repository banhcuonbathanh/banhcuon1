package user

import (
	"english-ai-full/ecomm-grpc/models"

)
type UserResponse struct {
	ID             uint    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
	EmailVerified  string `json:"emailVerified"`
	Image          string `json:"image"`
	FavoriteIds    string `json:"favoriteIds"`
	PhoneNumber    string `json:"phoneNumber"`
	StreetAddress  string `json:"streetAddress"`

    Orders    []models.Order
    // Accounts  []model.Account
 
    // Posts     []model.BlogPost   
    // Comments  []model.BlogComment
}

type LoginUserResponse struct {
	TokenType string `json:"token_type"`
	Token     string `json:"token"`
}


