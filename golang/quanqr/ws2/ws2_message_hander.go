package ws2

import (
	"encoding/json"
	"log"
)

type MessageHandler interface {
    Handle(client *Client, message Message)
}

type DefaultMessageHandler struct {

}

func (h *DefaultMessageHandler) HandleMessage(c *Client, msg Message) {

    log.Println("golang/quanqr/ws2/ws2_message_hander.go Handle 12121212121")
    switch msg.Type {
    case "order":
        h.handleOrderMessage(c, msg)
    case "notification":
        h.handleNotificationMessage(c, msg)
    case "status_update":
        h.handleStatusUpdate(c, msg)
    }
}

func (h *DefaultMessageHandler) handleOrderMessage(c *Client, msg Message) {
    log.Println("golang/quanqr/ws2/ws2_message_hander.go handleOrderMessage")
    var order OrderMessage
    data, _ := json.Marshal(msg.Payload)
    if err := json.Unmarshal(data, &order); err != nil {
        log.Printf("error unmarshaling order: %v", err)
        return
    }

    // Broadcast to kitchen staff
    notification := Message{
        Type:    "notification",
        Action:  "new_order",
        Payload: order,
        Role:    RoleKitchen,
    }
    
    data, _ = json.Marshal(notification)
    c.Hub.Broadcast <- data
}

func (h *DefaultMessageHandler) handleNotificationMessage(c *Client, msg Message) {
    log.Println("golang/quanqr/ws2/ws2_message_hander.go handleNotificationMessage")
    switch msg.Action {
    case "order_status":
        if msg.RoomID != "" {
            data, _ := json.Marshal(msg)
            if room := c.Hub.RoomMap[msg.RoomID]; room != nil {
                for client := range room {
                    select {
                    case client.Send <- data:
                    default:
                        close(client.Send)
                        delete(room, client)
                    }
                }
            }
        }
    }
}

func (h *DefaultMessageHandler) handleStatusUpdate(c *Client, msg Message) {
    log.Println("golang/quanqr/ws2/ws2_message_hander.go handleStatusUpdate")
    data, _ := json.Marshal(msg)
    c.Hub.Broadcast <- data
}
