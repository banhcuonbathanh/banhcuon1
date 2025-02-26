package order_grpc

import (
	"time"
)

// OrderDish matches DishOrderItem from proto
type OrderDish struct {
	ID                int64     `json:"id"`
	DishID            int64     `json:"dish_id"`
	Quantity          int64     `json:"quantity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	OrderName         string    `json:"order_name"`
	ModificationType  string    `json:"modification_type"`
	ModificationNumber int32     `json:"modification_number"`
}

// CreateOrderDish matches CreateDishOrderItem from proto
type CreateOrderDish struct {
	DishID   int64 `json:"dish_id"`
	Quantity int64 `json:"quantity"`
}

// OrderSet matches SetOrderItem from proto
type OrderSet struct {
	ID                int64     `json:"id"`
	SetID             int64     `json:"set_id"`
	Quantity          int64     `json:"quantity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	OrderName         string    `json:"order_name"`
	ModificationType  string    `json:"modification_type"`
	ModificationNumber int32     `json:"modification_number"`
}

// OrderModification matches OrderModification from proto
type OrderModification struct {
	ID                 int64     `json:"id"`
	OrderID            int64     `json:"order_id"`
	ModificationNumber int32     `json:"modification_number"`
	ModificationType   string    `json:"modification_type"`
	ModifiedAt         time.Time `json:"modified_at"`
	ModifiedByUserID   int64     `json:"modified_by_user_id"`
	OrderName          string    `json:"order_name"`
}

// DeliveryStatus enum
type DeliveryStatus string

const (
	DeliveryStatusPending            DeliveryStatus = "PENDING"
	DeliveryStatusPartiallyDelivered DeliveryStatus = "PARTIALLY_DELIVERED"
	DeliveryStatusFullyDelivered     DeliveryStatus = "FULLY_DELIVERED"
	DeliveryStatusCancelled          DeliveryStatus = "CANCELLED"
)

// DishDelivery matches DishDelivery from proto
type DishDelivery struct {
	ID                 int64       `json:"id"`
	OrderID            int64       `json:"order_id"`
	OrderName          string      `json:"order_name"`
	GuestID            int64       `json:"guest_id,omitempty"`
	UserID             int64       `json:"user_id,omitempty"`
	TableNumber        int64       `json:"table_number,omitempty"`

	QuantityDelivered  int32       `json:"quantity_delivered"`
	DeliveryStatus     string      `json:"delivery_status"`
	DeliveredAt        time.Time   `json:"delivered_at,omitempty"`
	DeliveredByUserID  int64       `json:"delivered_by_user_id,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	DishID             int64       `json:"dish_id"`
	IsGuest            bool        `json:"is_guest"`
	ModificationNumber int32       `json:"modification_number"`
}

// CreateOrderRequestType matches CreateOrderRequest from proto
type CreateOrderRequestType struct {
	GuestID        int64       `json:"guest_id"`
	UserID         int64       `json:"user_id"`
	IsGuest        bool        `json:"is_guest"`
	TableNumber    int64       `json:"table_number"`
	OrderHandlerID int64       `json:"order_handler_id"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	TotalPrice     int32       `json:"total_price"`
	DishItems      []OrderDish `json:"dish_items"`
	SetItems       []OrderSet  `json:"set_items"`
	Topping        string      `json:"topping"`
	TrackingOrder  string      `json:"tracking_order"`
	TakeAway       bool        `json:"take_away"`
	ChiliNumber    int64       `json:"chili_number"`
	TableToken     string      `json:"table_token"`
	OrderName      string      `json:"order_name"`
	Version        int32       `json:"version"`

}

// UpdateOrderRequestType matches UpdateOrderRequest from proto
type UpdateOrderRequestType struct {
	ID             int64       `json:"id"`
	GuestID        int64       `json:"guest_id"`
	UserID         int64       `json:"user_id"`
	TableNumber    int64       `json:"table_number"`
	OrderHandlerID int64       `json:"order_handler_id"`
	Status         string      `json:"status"`
	TotalPrice     int32       `json:"total_price"`
	DishItems      []OrderDish `json:"dish_items"`
	SetItems       []OrderSet  `json:"set_items"`
	IsGuest        bool        `json:"is_guest"`
	Topping        string      `json:"topping"`
	TrackingOrder  string      `json:"tracking_order"`
	TakeAway       bool        `json:"take_away"`
	ChiliNumber    int64       `json:"chili_number"`
	TableToken     string      `json:"table_token"`
	OrderName      string      `json:"order_name"`
	Version        int32       `json:"version"`
	ParentOrderID  int64       `json:"parent_order_id"`
}

// OrderType matches Order from proto
type OrderType struct {
	ID             int64       `json:"id"`
	GuestID        int64       `json:"guest_id"`
	UserID         int64       `json:"user_id"`
	IsGuest        bool        `json:"is_guest"`
	TableNumber    int64       `json:"table_number"`
	OrderHandlerID int64       `json:"order_handler_id"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	TotalPrice     int32       `json:"total_price"`
	DishItems      []OrderDish `json:"dish_items"`
	SetItems       []OrderSet  `json:"set_items"`
	Topping        string      `json:"topping"`
	TrackingOrder  string      `json:"tracking_order"`
	TakeAway       bool        `json:"take_away"`
	ChiliNumber    int64       `json:"chili_number"`
	TableToken     string      `json:"table_token"`
	OrderName      string      `json:"order_name"`
	Version        int32       `json:"version"`
	ParentOrderID  int64       `json:"parent_order_id"`
}

// OrderDetailedDish matches OrderDetailedDish from proto
type OrderDetailedDish struct {
	DishID      int64  `json:"dish_id"`
	Quantity    int64  `json:"quantity"`
	Name        string `json:"name"`
	Price       int32  `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Status      string `json:"status"`
}

// OrderSetDetailed matches OrderSetDetailed from proto
type OrderSetDetailed struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Dishes      []OrderDetailedDish `json:"dishes"`
	UserID      int32             `json:"user_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	IsFavourite bool              `json:"is_favourite"`
	LikeBy      []int64           `json:"like_by"`
	IsPublic    bool              `json:"is_public"`
	Image       string            `json:"image"`
	Price       int32             `json:"price"`
	Quantity    int64             `json:"quantity"`
}

// OrderVersionSummary matches OrderVersionSummary from proto
type OrderVersionSummary struct {
	VersionNumber     int32     `json:"version_number"`

	ModificationType  string    `json:"modification_type"`
	ModifiedAt        time.Time `json:"modified_at"`
	DishesOrdered     []OrderDetailedDish `json:"dishes_ordered"`
	SetOrdered        []OrderSetDetailed  `json:"set_ordered"`
}

// OrderItemChange matches OrderItemChange from proto
type OrderItemChange struct {
	ItemType        string `json:"item_type"`
	ItemID          int64  `json:"item_id"`
	ItemName        string `json:"item_name"`
	QuantityChanged int32  `json:"quantity_changed"`
	Price           int32  `json:"price"`
}

// OrderTotalSummary matches OrderTotalSummary from proto
type OrderTotalSummary struct {
	TotalVersions        int32            `json:"total_versions"`
	TotalDishesOrdered   int32            `json:"total_dishes_ordered"`
	TotalSetsOrdered     int32            `json:"total_sets_ordered"`
	CumulativeTotalPrice int32            `json:"cumulative_total_price"`
	MostOrderedItems     []OrderItemCount `json:"most_ordered_items"`
}

// OrderItemCount matches OrderItemCount from proto
type OrderItemCount struct {
	ItemType     string `json:"item_type"`
	ItemID       int64  `json:"item_id"`
	ItemName     string `json:"item_name"`
	TotalQuantity int32  `json:"total_quantity"`
}

// OrderDetailedResponse matches OrderDetailedResponseWithDelivery from proto
type OrderDetailedResponse struct {
	ID                   int64                 `json:"id"`
	GuestID              int64                 `json:"guest_id"`
	UserID               int64                 `json:"user_id"`
	TableNumber          int64                 `json:"table_number"`
	OrderHandlerID       int64                 `json:"order_handler_id"`
	Status               string                `json:"status"`
	TotalPrice           int32                 `json:"total_price"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	DataSet              []OrderSetDetailed    `json:"data_set"`
	DataDish             []OrderDetailedDish   `json:"data_dish"`
	IsGuest              bool                  `json:"is_guest"`
	Topping              string                `json:"topping"`
	TrackingOrder        string                `json:"tracking_order"`
	TakeAway             bool                  `json:"take_away"`
	ChiliNumber          int64                 `json:"chili_number"`
	TableToken           string                `json:"table_token"`
	OrderName            string                `json:"order_name"`
	CurrentVersion       int32                 `json:"current_version"`
	ParentOrderID        int64                 `json:"parent_order_id"`
	VersionHistory       []OrderVersionSummary `json:"version_history"`
	TotalSummary         OrderTotalSummary     `json:"total_summary"`
	DeliveryHistory      []DishDelivery        `json:"delivery_history"`
	CurrentDeliveryStatus DeliveryStatus       `json:"current_delivery_status"`
	TotalItemsDelivered   int32                `json:"total_items_delivered"`
	LastDeliveryAt        time.Time            `json:"last_delivery_at"`
}

// GetOrdersRequestType matches GetOrdersRequest from proto
type GetOrdersRequestType struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

// PaginationInfo matches PaginationInfo from proto
type PaginationInfo struct {
	CurrentPage int32 `json:"current_page"`
	TotalPages  int32 `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PageSize    int32 `json:"page_size"`
}

// OrderListResponse matches OrderListResponse from proto
type OrderListResponse struct {
	Data       []OrderType    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// PayOrdersRequestType matches PayOrdersRequest from proto
type PayOrdersRequestType struct {
	GuestID *int64 `json:"guest_id,omitempty"`
	UserID  *int64 `json:"user_id,omitempty"`
}

// OrderResponse matches OrderResponse from proto
type OrderResponse struct {
	Data OrderType `json:"data"`
}

// OrderIDParam matches OrderIdParam from proto
type OrderIDParam struct {
	ID int64 `json:"id"`
}

// OrderDetailIDParam matches OrderDetailIdParam from proto
type OrderDetailIDParam struct {
	ID int64 `json:"id"`
}

// OrderDetailedListResponse matches OrderDetailedListResponse from proto
type OrderDetailedListResponse struct {
	Data       []OrderDetailedResponse `json:"data"`
	Pagination PaginationInfo          `json:"pagination"`
}

// FetchOrdersByCriteriaRequestType matches FetchOrdersByCriteriaRequest from proto
type FetchOrdersByCriteriaRequestType struct {
	OrderIds    []int64    `json:"order_ids,omitempty"`
	OrderName   string     `json:"order_name,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Page        int32      `json:"page"`
	PageSize    int32      `json:"page_size"`
}

// CreateDishDeliveryRequestType matches CreateDishDeliveryRequest from proto
type CreateDishDeliveryRequestType struct {
	OrderID           int64            `json:"order_id"`
	OrderName         string           `json:"order_name"`
	GuestID           int64            `json:"guest_id,omitempty"`
	UserID            int64            `json:"user_id,omitempty"`
	TableNumber       int64            `json:"table_number,omitempty"`
DishID             int64       `json:"dish_id"`
	QuantityDelivered int32            `json:"quantity_delivered"`
	DeliveryStatus    string           `json:"delivery_status"`
	DeliveredAt       time.Time        `json:"delivered_at,omitempty"`
	DeliveredByUserID int64            `json:"delivered_by_user_id,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	IsGuest           bool             `json:"is_guest"`
}

// CreateDishOrderItemWithOrderID matches CreateDishOrderItemWithOrderID from proto
type CreateDishOrderItemWithOrderID struct {
	OrderID   int64  `json:"order_id"`
	DishID    int64  `json:"dish_id"`
	Quantity  int64  `json:"quantity"`
	OrderName string `json:"order_name"`
}

// CreateSetOrderItemWithOrderID matches CreateSetOrderItemWithOrderID from proto
type CreateSetOrderItemWithOrderID struct {
	OrderID   int64  `json:"order_id"`
	SetID     int64  `json:"set_id"`
	Quantity  int64  `json:"quantity"`
	OrderName string `json:"order_name"`
}

// ResponseSetOrderItemWithOrderID matches ResponseSetOrderItemWithOrderID from proto
type ResponseSetOrderItemWithOrderID struct {
	Set    OrderSet           `json:"set"`
	Dishes []OrderDetailedDish `json:"dishes"`
}


type OrderDetailedResponseWithDelivery struct {
	ID                   int64                 `json:"id"`
	GuestID              int64                 `json:"guest_id"`
	UserID               int64                 `json:"user_id"`
	TableNumber          int64                 `json:"table_number"`
	OrderHandlerID       int64                 `json:"order_handler_id"`
	Status               string                `json:"status"`

	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	// DataSet and DataDish are commented out in the proto
	// DataSet              []OrderSetDetailed    `json:"data_set"`
	// DataDish             []OrderDetailedDish   `json:"data_dish"`
	IsGuest              bool                  `json:"is_guest"`
	Topping              string                `json:"topping"`
	TrackingOrder        string                `json:"tracking_order"`
	TakeAway             bool                  `json:"take_away"`
	ChiliNumber          int64                 `json:"chili_number"`
	TableToken           string                `json:"table_token"`
	OrderName            string                `json:"order_name"`
	CurrentVersion       int32                 `json:"current_version"`
	// ParentOrderID is commented out in the proto
	// ParentOrderID        int64                 `json:"parent_order_id"`
	VersionHistory       []OrderVersionSummary `json:"version_history"`
	// TotalSummary is commented out in the proto
	// TotalSummary         OrderTotalSummary     `json:"total_summary"`
	
	// New delivery-related fields
	DeliveryHistory      []DishDelivery        `json:"delivery_history"`
	CurrentDeliveryStatus DeliveryStatus       `json:"current_delivery_status"`
	TotalItemsDelivered   int32                `json:"total_items_delivered"`
	LastDeliveryAt        time.Time            `json:"last_delivery_at"`
}

type OrderDetailedListResponseWithDelivery struct {
	Data       []OrderDetailedResponseWithDelivery `json:"data"`
	Pagination PaginationInfo          `json:"pagination"`
}