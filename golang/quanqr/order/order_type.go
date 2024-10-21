package order_grpc

import (
	"time"
)
type OrderType struct {
    ID             int64           `json:"id"`
    GuestID        int64           `json:"guest_id"`
    UserID         int64           `json:"user_id"`
    IsGuest        bool            `json:"is_guest"`
    TableNumber    int64           `json:"table_number"`
    OrderHandlerID int64           `json:"order_handler_id"`
    Status         string          `json:"status"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
    TotalPrice     int32           `json:"total_price"`
    DishItems      []DishOrderItem `json:"dish_items"`
    SetItems       []SetOrderItem  `json:"set_items"`
    BowChili       int64           `json:"bow_chili"`
    BowNoChili     int64           `json:"bow_no_chili"`
}

// CreateOrderRequest struct
type CreateOrderRequestType struct {
    GuestID        int64           `json:"guest_id"`
    UserID         int64           `json:"user_id"`
    IsGuest        bool            `json:"is_guest"`
    TableNumber    int64           `json:"table_number"`
    OrderHandlerID int64           `json:"order_handler_id"`
    Status         string          `json:"status"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
    TotalPrice     int32           `json:"total_price"`
    DishItems      []DishOrderItem `json:"dish_items"`
    SetItems       []SetOrderItem  `json:"set_items"`
    BowChili       int64           `json:"bow_chili"`
    BowNoChili     int64           `json:"bow_no_chili"`
}

// UpdateOrderRequest struct
type UpdateOrderRequestType struct {
    ID             int64           `json:"id"`
    GuestID        int64           `json:"guest_id"`
    UserID         int64           `json:"user_id"`
    TableNumber    int64           `json:"table_number"`
    OrderHandlerID int64           `json:"order_handler_id"`
    Status         string          `json:"status"`
    TotalPrice     int32           `json:"total_price"`
    DishItems      []DishOrderItem `json:"dish_items"`
    SetItems       []SetOrderItem  `json:"set_items"`
    IsGuest        bool            `json:"is_guest"`
    BowChili       int64           `json:"bow_chili"`
    BowNoChili     int64           `json:"bow_no_chili"`
}

// DishOrder struct
type DishOrder struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       int32     `json:"price"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DishOrderItem struct
type DishOrderItem struct {
    ID       int64 `json:"id"`
    Quantity int32 `json:"quantity"`
}

// SetOrderItem struct
type SetOrderItem struct {
    ID       int64 `json:"id"`
    Quantity int32 `json:"quantity"`
}





// GetOrdersRequest struct
type GetOrdersRequestType struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	UserID   *int64    `json:"user_id,omitempty"`
	GuestID  *int64    `json:"guest_id,omitempty"`
}

// PayOrdersRequest struct
type PayOrdersRequestType struct {
	GuestID *int64 `json:"guest_id,omitempty"`
	UserID  *int64 `json:"user_id,omitempty"`
}

// OrderResponse struct
type OrderResponse struct {
	Data OrderType `json:"data"`
}

// OrderListResponse struct
type OrderListResponse struct {
	Data []OrderType `json:"data"`
}

// OrderIDParam struct
type OrderIDParam struct {
	ID int64 `json:"id"`
}

// OrderDetailIDParam struct
type OrderDetailIDParam struct {
	ID int64 `json:"id"`
}

// Guest struct
type Guest struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	TableNumber int32     `json:"table_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
