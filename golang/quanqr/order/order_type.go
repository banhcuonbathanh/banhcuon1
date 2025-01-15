package order_grpc

import (
	"time"
)

// OrderDish now includes modification tracking fields
type OrderDish struct {
    ID                int64     `json:"id"`                    // Added primary key
    DishID           int64     `json:"dish_id"`
    Quantity         int64     `json:"quantity"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    OrderName        string    `json:"order_name"`
    ModificationType string    `json:"modification_type"`      // Track if item was initial or added later
    ModificationNumber int32   `json:"modification_number"`    // Track which modification this was part of
}

// OrderSet now includes modification tracking fields
type OrderSet struct {
    ID                int64     `json:"id"`                    // Added primary key
    SetID            int64     `json:"set_id"`
    Quantity         int64     `json:"quantity"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    OrderName        string    `json:"order_name"`
    ModificationType string    `json:"modification_type"`      // Track if item was initial or added later
    ModificationNumber int32   `json:"modification_number"`    // Track which modification this was part of
}

// New type for tracking order modifications
type OrderModification struct {
    ID                int64     `json:"id"`
    OrderID          int64     `json:"order_id"`
    ModificationNumber int32    `json:"modification_number"`
    ModificationType  string    `json:"modification_type"`
    ModifiedAt       time.Time `json:"modified_at"`
    ModifiedByUserID int64     `json:"modified_by_user_id"`
    OrderName        string    `json:"order_name"`
}

// New type for tracking dish deliveries
type DishDelivery struct {
    ID                int64     `json:"id"`
    OrderID          int64     `json:"order_id"`
    OrderName        string    `json:"order_name"`
    QuantityDelivered int32    `json:"quantity_delivered"`
    DeliveryStatus   string    `json:"delivery_status"`
    DeliveredAt      time.Time `json:"delivered_at"`
    DeliveredByUserID int64    `json:"delivered_by_user_id"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    ModificationNumber int32   `json:"modification_number"`
    DishOrderItemID  int64     `json:"dish_order_item_id"`
}

// CreateOrderRequestType updated with new fields
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
    DishItems      []OrderDish     `json:"dish_items"`
    SetItems       []OrderSet      `json:"set_items"`
    Topping        string          `json:"topping"`
    TrackingOrder  string          `json:"tracking_order"`
    TakeAway       bool            `json:"take_away"`
    ChiliNumber    int64           `json:"chili_number"`
    TableToken     string          `json:"table_token"`
    OrderName      string          `json:"order_name"`
    Version        int32           `json:"version"`
    ParentOrderID  int64           `json:"parent_order_id"`
}

// UpdateOrderRequestType updated with new fields
type UpdateOrderRequestType struct {
    ID             int64           `json:"id"`
    GuestID        int64           `json:"guest_id"`
    UserID         int64           `json:"user_id"`
    TableNumber    int64           `json:"table_number"`
    OrderHandlerID int64           `json:"order_handler_id"`
    Status         string          `json:"status"`
    TotalPrice     int32           `json:"total_price"`
    DishItems      []OrderDish     `json:"dish_items"`
    SetItems       []OrderSet      `json:"set_items"`
    IsGuest        bool            `json:"is_guest"`
    Topping        string          `json:"topping"`
    TrackingOrder  string          `json:"tracking_order"`
    TakeAway       bool            `json:"take_away"`
    ChiliNumber    int64           `json:"chili_number"`
    TableToken     string          `json:"table_token"`
    OrderName      string          `json:"order_name"`
    Version        int32           `json:"version"`
    ParentOrderID  int64           `json:"parent_order_id"`
}

// OrderDetailedResponse updated with new fields


// Existing support types that remain unchanged
type GetOrdersRequestType struct {
    Page     int32 `json:"page"`
    PageSize int32 `json:"page_size"`
}

type PaginationInfo struct {
    CurrentPage int32 `json:"current_page"`
    TotalPages  int32 `json:"total_pages"`
    TotalItems  int64 `json:"total_items"`
    PageSize    int32 `json:"page_size"`
}

type OrderListResponse struct {
    Data       []OrderType    `json:"data"`
    Pagination PaginationInfo `json:"pagination"`
}

type PayOrdersRequestType struct {
    GuestID *int64 `json:"guest_id,omitempty"`
    UserID  *int64 `json:"user_id,omitempty"`
}

type OrderResponse struct {
    Data OrderType `json:"data"`
}

type OrderIDParam struct {
    ID int64 `json:"id"`
}

type OrderDetailIDParam struct {
    ID int64 `json:"id"`
}

type OrderDetailedDish struct {
    DishID      int64  `json:"dish_id"`
    Quantity    int64  `json:"quantity"`
    Name        string `json:"name"`
    Price       int32  `json:"price"`
    Description string `json:"description"`
    Image       string `json:"image"`
    Status      string `json:"status"`
}

type OrderSetDetailed struct {
    ID          int64             `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Dishes      []OrderDetailedDish `json:"dishes"`
    UserID      int32             `json:"userId"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    IsFavourite bool              `json:"is_favourite"`
    LikeBy      []int64           `json:"like_by"`
    IsPublic    bool              `json:"is_public"`
    Image       string            `json:"image"`
    Price       int32             `json:"price"`
    Quantity    int64             `json:"quantity"`
}

type OrderDetailedListResponse struct {
    Data       []OrderDetailedResponse `json:"data"`
    Pagination PaginationInfo         `json:"pagination"`
}

type FetchOrdersByCriteriaRequestType struct {
    OrderIds    []int64    `json:"order_ids,omitempty"`
    OrderName   string     `json:"order_name,omitempty"`
    StartDate   *time.Time `json:"start_date,omitempty"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    Page        int32      `json:"page"`
    PageSize    int32      `json:"page_size"`
}

// new



// OrderType represents the main order structure
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
    DishItems      []OrderDish     `json:"dish_items"`
    SetItems       []OrderSet      `json:"set_items"`
    Topping        string          `json:"topping"`
    TrackingOrder  string          `json:"tracking_order"`
    TakeAway       bool            `json:"take_away"`
    ChiliNumber    int64           `json:"chili_number"`
    TableToken     string          `json:"table_token"`
    OrderName      string          `json:"order_name"`
    Version        int32           `json:"version"`
    ParentOrderID  int64           `json:"parent_order_id"`
}

// OrderVersionSummary represents a summary of items in a specific order version
type OrderVersionSummary struct {
    VersionNumber     int32     `json:"version_number"`
    TotalDishesCount  int32     `json:"total_dishes_count"`
    TotalSetsCount    int32     `json:"total_sets_count"`
    VersionTotalPrice int32     `json:"version_total_price"`
    ModificationType  string    `json:"modification_type"`
    ModifiedAt       time.Time  `json:"modified_at"`
    Changes          []OrderItemChange `json:"changes"`
}

// OrderItemChange represents changes made in a specific version
type OrderItemChange struct {
    ItemType       string `json:"item_type"`
    ItemID         int64  `json:"item_id"`
    ItemName       string `json:"item_name"`
    QuantityChanged int32 `json:"quantity_changed"`
    Price          int32  `json:"price"`
}

// OrderTotalSummary represents the total summary across all versions
type OrderTotalSummary struct {
    TotalVersions       int32           `json:"total_versions"`
    TotalDishesOrdered  int32           `json:"total_dishes_ordered"`
    TotalSetsOrdered    int32           `json:"total_sets_ordered"`
    CumulativeTotalPrice int32          `json:"cumulative_total_price"`
    MostOrderedItems    []OrderItemCount `json:"most_ordered_items"`
}

// OrderItemCount helps track most ordered items
type OrderItemCount struct {
    ItemType     string `json:"item_type"`
    ItemID       int64  `json:"item_id"`
    ItemName     string `json:"item_name"`
    TotalQuantity int32 `json:"total_quantity"`
}

// OrderDetailedResponse updated with version tracking fields
type OrderDetailedResponse struct {
    DataSet         []OrderSetDetailed    `json:"data_set"`
    DataDish        []OrderDetailedDish   `json:"data_dish"`
    ID              int64                 `json:"id"`
    GuestID         int64                 `json:"guest_id"`
    UserID          int64                 `json:"user_id"`
    TableNumber     int64                 `json:"table_number"`
    OrderHandlerID  int64                 `json:"order_handler_id"`
    Status          string                `json:"status"`
    TotalPrice      int32                 `json:"total_price"`
    IsGuest         bool                  `json:"is_guest"`
    Topping         string                `json:"topping"`
    TrackingOrder   string                `json:"tracking_order"`
    TakeAway        bool                  `json:"take_away"`
    ChiliNumber     int64                 `json:"chili_number"`
    TableToken      string                `json:"table_token"`
    OrderName       string                `json:"order_name"`
    CurrentVersion  int32                 `json:"current_version"`  // Renamed from Version
    ParentOrderID   int64                 `json:"parent_order_id"`
    VersionHistory  []OrderVersionSummary `json:"version_history"` // Added field
    TotalSummary    OrderTotalSummary     `json:"total_summary"`   // Added field
}

// Rest of the existing types remain unchanged