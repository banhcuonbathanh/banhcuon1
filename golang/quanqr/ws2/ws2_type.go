package ws2

type Role string

const (
    RoleGuest    Role = "Guest"
    RoleUser     Role = "User"
    RoleEmployee Role = "Employee"
    RoleAdmin    Role = "Admin"
    RoleKitchen  Role = "Kitchen"
)

type Message struct {
    Type    string      `json:"type"`
    Action  string      `json:"action"`
    Payload interface{} `json:"payload"`
    Role    Role        `json:"role"`
    RoomID  string      `json:"roomId,omitempty"`
}

type OrderMessage struct {
    OrderID      int64  `json:"orderId"`
    GuestID      *int64 `json:"guestId,omitempty"`
    UserID       *int64 `json:"userId,omitempty"`
    TableNumber  int64  `json:"tableNumber"`
    IsGuest      bool   `json:"isGuest"`
    TableToken   string `json:"tableToken"`
    TotalPrice   int    `json:"totalPrice"`
    BowChili     int64  `json:"bowChili"`
    BowNoChili   int64  `json:"bowNoChili"`
    TakeAway     bool   `json:"takeAway"`
    ChiliNumber  int64  `json:"chiliNumber"`
    OrderName    string `json:"orderName"`
}

type DirectMessage struct {
    FromUserID string      `json:"fromUserId"`
    ToUserID   string      `json:"toUserId"`
    Type       string      `json:"type"`
    Action     string      `json:"action"`
    Payload    interface{} `json:"payload"`
}