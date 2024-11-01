package websocket_model

import "time"

type Message struct {
    Type      string      `json:"type"`
    Content   interface{} `json:"content"`    // Changed to interface{} to handle various content types
    Sender    string      `json:"sender"`
    Timestamp time.Time   `json:"timestamp"`
    TableID   string      `json:"table_id,omitempty"`
    OrderID   string      `json:"order_id,omitempty"`

    ID        string    `json:"id,omitempty"`
    FromUser  string    `json:"fromUser"`     // Sender's userID
    ToUser    string    `json:"toUser"`       // Recipient's userID
}


type OrderMessage struct {
    OrderID      string    `json:"order_id"`
    TableNumber  string    `json:"table_number"`
    Status       string    `json:"status"`
    Timestamp    time.Time `json:"timestamp"`
    OrderData    interface{} `json:"order_data"`
}
