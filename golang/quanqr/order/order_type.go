package order_grpc

import (
	"time"
)

type Order struct {
	ID              int64     `json:"id"`
	GuestID         int64     `json:"guest_id"`
	TableNumber     int32     `json:"table_number"`
	DishSnapshotID  int64     `json:"dish_snapshot_id"`
	Quantity        int32     `json:"quantity"`
	OrderHandlerID  int64     `json:"order_handler_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Guest struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	TableNumber int32     `json:"table_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DishSnapshot struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       int32     `json:"price"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DishID      int64     `json:"dish_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Account struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
}

type Table struct {
	Number    int32     `json:"number"`
	Capacity  int32     `json:"capacity"`
	Status    string    `json:"status"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderDetail struct {
	Order        Order        `json:"order"`
	Guest        Guest        `json:"guest"`
	DishSnapshot DishSnapshot `json:"dish_snapshot"`
	OrderHandler Account      `json:"order_handler"`
	Table        Table        `json:"table"`
}

type CreateOrderItem struct {
	DishID   int64 `json:"dish_id"`
	Quantity int32 `json:"quantity"`
}

type CreateOrdersRequest struct {
	GuestID int64              `json:"guest_id"`
	Orders  []CreateOrderItem  `json:"orders"`
}

type UpdateOrderRequest struct {
	OrderID  int64  `json:"order_id"`
	Status   string `json:"status"`
	DishID   int64  `json:"dish_id"`
	Quantity int32  `json:"quantity"`
}

type GetOrdersRequest struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type PayGuestOrdersRequest struct {
	GuestID int64 `json:"guest_id"`
}

type OrderResponse struct {
	Message string      `json:"message"`
	Data    OrderDetail `json:"data"`
}

type OrderListResponse struct {
	Message string        `json:"message"`
	Data    []OrderDetail `json:"data"`
}

type OrderDetailIDParam struct {
	ID int64 `json:"id"`
}