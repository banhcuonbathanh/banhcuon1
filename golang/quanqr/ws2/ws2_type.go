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



type DirectMessage struct {
    FromUserID string      `json:"fromUserId"`
    ToUserID   string      `json:"toUserId"`
    Type       string      `json:"type"`
    Action     string      `json:"action"`
    Payload    interface{} `json:"payload"`
}

