package types

import "time"
type ProductReq struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Rating       int64   `json:"rating"`
	NumReviews   int64   `json:"num_reviews"`
	Price        float32 `json:"price"`
	CountInStock int64   `json:"count_in_stock"`
}

type ProductRes struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Image        string     `json:"image"`
	Category     string     `json:"category"`
	Description  string     `json:"description"`
	Rating       int64      `json:"rating"`
	NumReviews   int64      `json:"num_reviews"`
	Price        float32    `json:"price"`
	CountInStock int64      `json:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type OrderReq struct {
	Items         []*OrderItem `json:"items"`
	PaymentMethod string       `json:"payment_method"`
	TaxPrice      float32      `json:"tax_price"`
	ShippingPrice float32      `json:"shipping_price"`
	TotalPrice    float32      `json:"total_price"`
}

type OrderItem struct {
	Name      string  `json:"name"`
	Quantity  int64   `json:"quantity"`
	Image     string  `json:"image"`
	Price     float32 `json:"price"`
	ProductID int64   `json:"product_id"`
}

type OrderRes struct {
	ID            int64        `json:"id"`
	Items         []*OrderItem `json:"items"`
	PaymentMethod string       `json:"payment_method"`
	TaxPrice      float32      `json:"tax_price"`
	ShippingPrice float32      `json:"shipping_price"`
	TotalPrice    float32      `json:"total_price"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
}

type UserReqModel struct {
    Password  string     `json:"password"`
    ID        int64      `json:"id"`
    Name      string     `json:"name"`
    Email     string     `json:"email"`
    Role      string     `json:"role"`  // Changed from IsAdmin bool to Role string
    Phone     *int64     `json:"phone"`
    Image     string     `json:"image"`
    Address   string     `json:"address"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

type UserResModel struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`  // Changed from IsAdmin bool to Role string
	Phone     int64     `json:"phone"`
	Image     string    `json:"image"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUserRes struct {
	Users []UserReqModel `json:"users"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserReq struct {
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`  // Changed from IsAdmin bool to Role string
	Phone     int64     `json:"phone"`
	Image     string    `json:"image"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserRes struct {
	SessionID             string       `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResModel `json:"user"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type Session struct {
    ID           string
    UserEmail    string
    RefreshToken string
    IsRevoked    bool
    ExpiresAt    time.Time
}