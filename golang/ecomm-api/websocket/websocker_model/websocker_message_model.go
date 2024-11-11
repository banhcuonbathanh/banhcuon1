// websocket_model/models.go
package websocket_model

import "time"

type UserWS struct {
    ID   int64    `json:"id"`
    Name string `json:"name"`
}

type GuestInfoWS struct {
    ID   int64    `json:"id"`
    Name string `json:"name"`
}

type Message struct {
    Type      string      `json:"type"`
    Content   interface{} `json:"content"`
    Sender    int64      `json:"sender"`
    Recipient string      `json:"recipient,omitempty"` // Add this field
    Timestamp time.Time   `json:"timestamp"`
    TableID   int64      `json:"table_id,omitempty"`
    OrderID   int64      `json:"order_id,omitempty"`
    ID        int64      `json:"id,omitempty"`
    FromUser  int64      `json:"fromUser"`
    ToUser    int64      `json:"toUser"`
}

type DishOrderItem struct {
    DishID   int64 `json:"dish_id"`
    Quantity int64 `json:"quantity"`
}



type SetOrderItem struct {
   SetID   int64 `json:"set_id"`
    Quantity int64 `json:"quantity"`
}





type CreateOrderRequest struct {
    GuestID        *int64            `json:"guest_id"`
    UserID         *int64            `json:"user_id"`
    IsGuest        bool            `json:"is_guest"`
    TableNumber    int             `json:"table_number"`
    OrderHandlerID int             `json:"order_handler_id"`
    Status         string          `json:"status"`
    CreatedAt      string          `json:"created_at"`
    UpdatedAt      string          `json:"updated_at"`
    TotalPrice     float64         `json:"total_price"`
    DishItems      []DishOrderItem `json:"dish_items"`
    SetItems       []SetOrderItem  `json:"set_items"`
    BowChili       int             `json:"bow_chili"`
    BowNoChili     int             `json:"bow_no_chili"`
    TakeAway       bool            `json:"takeAway"`
    ChiliNumber    int             `json:"chiliNumber"`
    TableToken     string          `json:"table_token"`
    OrderName      string          `json:"order_name"`
}

type OrderPayload struct {
    OrderID   int               `json:"orderId"`
    OrderData CreateOrderRequest `json:"orderData"`
}

type OrderStatusUpdate struct {
    OrderID   int       `json:"orderId"`
    Status    string    `json:"status"`
    Timestamp string    `json:"timestamp"`
}


// type DishOrderItemDetail struct {
//     DishID   int64 `json:"dish_id"`
//     Quantity int64 `json:"quantity"`
// }

// type SetOrderItemDetail struct {
//     SetID   int64 `json:"set_id"`
//      Quantity int64 `json:"quantity"`
//  }