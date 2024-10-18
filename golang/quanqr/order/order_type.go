package order_grpc

import (
	"time"
)

// Order struct aligned with the proto definition
type OrderType struct {
	ID             int64          `json:"id"`
	GuestID        int64         `json:"guest_id"`
	UserID         int64         `json:"user_id"`
	IsGuest			bool			`json:"is_guest"`
	TableNumber    int64         `json:"table_number"`
	OrderHandlerID int64         `json:"order_handler_id"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	TotalPrice     int32          `json:"total_price"`
	DishItems      []DishOrderItem `json:"dish_items"`
	SetItems       []SetOrderItemType  `json:"set_items"`
}

// CreateOrderRequest struct
type CreateOrderRequestType struct {
	GuestID        int64         `json:"guest_id"`
	UserID         int64         `json:"user_id"`
	IsGuest			bool			`json:"is_guest"`
	TableNumber    int64         `json:"table_number"`
	OrderHandlerID int64         `json:"order_handler_id"`
	Status         string         `json:"status"`
	TotalPrice     int32          `json:"total_price"`
	DishItems      []CreateOrderItemType `json:"dish_items"`
	SetItems       []CreateOrderItemType `json:"set_items"`
}

// UpdateOrderRequest struct
type UpdateOrderRequestType struct {
	ID             int64          `json:"id"`
	GuestID        int64         `json:"guest_id"`
	UserID         int64         `json:"user_id"`
		IsGuest			bool			`json:"is_guest"`
	TableNumber    int64         `json:"table_number"`
	OrderHandlerID int64         `json:"order_handler_id"`
	Status         string         `json:"status"`
	TotalPrice     int32          `json:"total_price"`
	DishItems      []DishOrderItem `json:"dish_items"`
	SetItems       []SetOrderItemType  `json:"set_items"`
}


type DishOrderType struct {
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
	ID              int64       `json:"id"`


	Quantity        int32       `json:"quantity"`

	Dish 		DishOrderType `json:"dish"`
	// DishSnapshot    DishSnapshot `json:"dish_snapshot,omitempty"` // For response data
}

// SetOrderItem struct
type SetOrderItemType struct {
	ID             int64          `json:"id"`

	Quantity       int32          `json:"quantity"`
	Set            SetProtoType       `json:"set,omitempty"` // For response data
	ModifiedDishes []SetProtoDishType `json:"modified_dishes"`
}

// CreateOrderItem struct
type CreateOrderItemType struct {
	DishSnapshotID *int64         `json:"dish_snapshot_id,omitempty"`
	SetSnapshotID  *int64         `json:"set_snapshot_id,omitempty"`
	Quantity       int32          `json:"quantity"`
	ModifiedDishes []SetProtoDishType `json:"modified_dishes,omitempty"`
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

// SetProto struct
type SetProtoType struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`

	UserID      int64         `json:"user_id"`

	IsFavourite bool           `json:"is_favourite"`
	LikeBy      []int64        `json:"like_by"`
			CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	IsPublic    bool           `json:"is_public"`
	Image       string         `json:"image"`

		Dishes      []SetProtoDishType `json:"dishes"`
	
}

// SetProtoDish struct
type SetProtoDishType struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

// Guest struct
type GuestType struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	TableNumber int32     `json:"table_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Account struct
type AccountType struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
}

// Table struct
type TableType struct {
	Number    int32     `json:"number"`
	Capacity  int32     `json:"capacity"`
	Status    string    `json:"status"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderDetailIDParam struct
type OrderDetailIDParamType struct {
	ID int64 `json:"id"`
}

// Order struct aligned with the proto definition
// type Order struct {
// 	ID              int64                  `json:"id"`
// 	GuestID         int64                  `json:"guest_id"`
// 	TableNumber     int64                  `json:"table_number"`
// 	DishSnapshotID  int64                  `json:"dish_snapshot_id"`
// 	OrderHandlerID  int64                  `json:"order_handler_id"`
// 	Status          string                 `json:"status"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	TotalPrice      int32                  `json:"total_price"`
// 	DishItems       []DishOrderItem        `json:"dish_items"`
// 	SetItems        []SetOrderItem         `json:"set_items"`
// }

// // CreateOrderRequest struct
// type CreateOrderRequest struct {
// 	GuestID         int64                  `json:"guest_id"`
// 	TableNumber     int64                  `json:"table_number"`
// 	DishSnapshotID  int64                  `json:"dish_snapshot_id"`
// 	OrderHandlerID  int64                  `json:"order_handler_id"`
// 	Status          string                 `json:"status"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	TotalPrice      int32                  `json:"total_price"`
// 	DishItems       []DishOrderItem        `json:"dish_items"`
// 	SetItems        []SetOrderItem         `json:"set_items"`
// }

// // UpdateOrderRequest struct
// type UpdateOrderRequest struct {
// 	ID              int64                  `json:"id"`
// 	GuestID         int64                  `json:"guest_id"`
// 	TableNumber     int64                  `json:"table_number"`
// 	DishSnapshotID  int64                  `json:"dish_snapshot_id"`
// 	OrderHandlerID  int64                  `json:"order_handler_id"`
// 	Status          string                 `json:"status"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	TotalPrice      int32                  `json:"total_price"`
// 	DishItems       []DishOrderItem        `json:"dish_items"`
// 	SetItems        []SetOrderItem         `json:"set_items"`
// }

// // OrderResponse struct
// type OrderResponse struct {
// 	Data Order `json:"data"`
// }

// // OrderListResponse struct
// type OrderListResponse struct {
// 	Data []Order `json:"data"`
// }

// // OrderIDParam struct
// type OrderIDParam struct {
// 	ID int64 `json:"id"`
// }

// // DishOrderItem struct
// type DishOrderItem struct {
// 	ID       int64  `json:"id"`
// 	Quantity int32  `json:"quantity"`
// 	Dish     Dish   `json:"dish"`
// }

// // SetOrderItem struct
// type SetOrderItem struct {
// 	ID             int64          `json:"id"`
// 	Quantity       int32          `json:"quantity"`
// 	Set            SetProto       `json:"set"`
// 	ModifiedDishes []SetProtoDish `json:"modified_dishes"`
// }

// // Dish struct
// type Dish struct {
// 	ID          int64                  `json:"id"`
// 	Name        string                 `json:"name"`
// 	Price       int32                  `json:"price"`
// 	Description string                 `json:"description"`
// 	Image       string                 `json:"image"`
// 	Status      string                 `json:"status"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// // SetProto struct
// type SetProto struct {
// 	ID            int64           `json:"id"`
// 	Name          string          `json:"name"`
// 	Description   string          `json:"description"`
// 	Dishes        []SetProtoDish  `json:"dishes"`
// 	UserID        *int32          `json:"user_id,omitempty"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	IsFavourite   bool            `json:"is_favourite"`
// 	LikeBy        []int64         `json:"like_by"`
// 	IsPublic      bool            `json:"is_public"`
// 	Image         string          `json:"image"`
// }

// // SetProtoDish struct
// type SetProtoDish struct {
// 	ID    int64  `json:"id"`
// 	Name  string `json:"name"`
// 	Price int32  `json:"price"`
// }

// // CreateOrderItem struct
// type CreateOrderItem struct {
// 	ID             int64          `json:"id"`
// 	Quantity       int32          `json:"quantity"`
// 	ModifiedDishes []SetProtoDish `json:"modified_dishes"`
// }

// // GetOrdersRequest struct
// type GetOrdersRequest struct {
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// // Guest struct
// type Guest struct {
// 	ID          int64                  `json:"id"`
// 	Name        string                 `json:"name"`
// 	TableNumber int32                  `json:"table_number"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// // DishSnapshot struct
// type DishSnapshot struct {
// 	ID          int64                  `json:"id"`
// 	Name        string                 `json:"name"`
// 	Price       int32                  `json:"price"`
// 	Image       string                 `json:"image"`
// 	Description string                 `json:"description"`
// 	Status      string                 `json:"status"`
// 	DishID      int64                  `json:"dish_id"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// // Account struct
// type Account struct {
// 	ID     int64  `json:"id"`
// 	Name   string `json:"name"`
// 	Email  string `json:"email"`
// 	Role   string `json:"role"`
// 	Avatar string `json:"avatar"`
// }

// // Table struct
// type Table struct {
// 	Number    int32                  `json:"number"`
// 	Capacity  int32                  `json:"capacity"`
// 	Status    string                 `json:"status"`
// 	Token     string                 `json:"token"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// // PayGuestOrdersRequest struct
// type PayGuestOrdersRequest struct {
// 	GuestID int64 `json:"guest_id"`
// }

// // OrderDetailIDParam struct
// type OrderDetailIDParam struct {
// 	ID int64 `json:"id"`
// }
