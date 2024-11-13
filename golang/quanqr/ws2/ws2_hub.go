package ws2

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Hub struct {
    Clients        map[*Client]bool
    Broadcast      chan []byte
    Register       chan *Client
    Unregister    chan *Client
    RoomMap        map[string]map[*Client]bool
    MessageHandler MessageHandler
    mu            sync.Mutex
}

func NewHub(messageHandler MessageHandler) *Hub {
    log.Printf("Creating new Hub with message handler type: %T", messageHandler)
    return &Hub{
        Broadcast:      make(chan []byte),
        Register:       make(chan *Client),
        Unregister:    make(chan *Client),
        Clients:        make(map[*Client]bool),
        RoomMap:        make(map[string]map[*Client]bool),
        MessageHandler: messageHandler,
    }
}
func (h *Hub) registerClient(client *Client) {
	log.Println("golang/quanqr/ws2/ws2_hub.go registerClient")
    h.mu.Lock()
    defer h.mu.Unlock()
    
    h.Clients[client] = true
    if client.RoomID != "" {
        if h.RoomMap[client.RoomID] == nil {
            h.RoomMap[client.RoomID] = make(map[*Client]bool)
        }
        h.RoomMap[client.RoomID][client] = true
    }
}

func (h *Hub) unregisterClient(client *Client) {
	log.Println("golang/quanqr/ws2/ws2_hub.go unregisterClient")
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if _, ok := h.Clients[client]; ok {
        delete(h.Clients, client)
        if client.RoomID != "" && h.RoomMap[client.RoomID] != nil {
            delete(h.RoomMap[client.RoomID], client)
            if len(h.RoomMap[client.RoomID]) == 0 {
                delete(h.RoomMap, client.RoomID)
            }
        }
        close(client.Send)
    }
}

func (h *Hub) broadcastMessage(message []byte) {
	log.Println("golang/quanqr/ws2/ws2_hub.go broadcastMessage")
    h.mu.Lock()
    defer h.mu.Unlock()
    
    for client := range h.Clients {
        select {
        case client.Send <- message:
        default:
            close(client.Send)
            delete(h.Clients, client)
        }
    }
}


func (h *Hub) SendDirectMessage(fromUserID, toUserID string, msgType, action string, payload interface{}) error {
	log.Println("golang/quanqr/ws2/ws2_hub.go SendDirectMessage")
    h.mu.Lock()
    defer h.mu.Unlock()


    var targetClient *Client
    for client := range h.Clients {
        if client.Role == RoleUser && client.ID == toUserID {
            targetClient = client
            break
        }
    }

    if targetClient == nil {
        return fmt.Errorf("target user %s not found", toUserID)
    }

    directMsg := DirectMessage{
        FromUserID: fromUserID,
        ToUserID:   toUserID,
        Type:       msgType,
        Action:     action,
        Payload:    payload,
    }

    // Wrap in standard Message format
    msg := Message{
        Type:    "direct",
        Action:  action,
        Payload: directMsg,
        Role:    RoleUser,
    }

    // Marshal the message
    data, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("error marshaling message: %v", err)
    }


    select {
    case targetClient.Send <- data:
        return nil
    default:
        close(targetClient.Send)
        delete(h.Clients, targetClient)
        return fmt.Errorf("failed to send message to user %s", toUserID)
    }
}

func (h *Hub) Run() {
	log.Println("golang/quanqr/ws2/ws2_hub.go Run")
    for {
        select {
        case client := <-h.Register:
            h.registerClient(client)
        case client := <-h.Unregister:
            h.unregisterClient(client)
        case message := <-h.Broadcast:
            h.broadcastMessage(message)
        }
    }
}

// new

func (h *Hub) BroadcastToStaff(fromUserID string, msg Message) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    // Add role information to the message
    msg.Role = RoleEmployee // Default to Employee role, staff can determine their actions based on their actual role

    // Marshal the message
    data, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("error marshaling message: %v", err)
    }

    // Keep track of successful sends
    successfulSends := 0
    staffRoles := map[Role]bool{
        RoleAdmin:    true,
        RoleEmployee: true,
        RoleKitchen:  true,
    }

    // Send to all staff clients (Admin, Employee, and Kitchen)
    for client := range h.Clients {
        if staffRoles[client.Role] {
            select {
            case client.Send <- data:
                successfulSends++
                log.Printf("Message sent from user %s to %s (ID: %s)", fromUserID, client.Role, client.ID)
            default:
                close(client.Send)
                delete(h.Clients, client)
                log.Printf("Failed to send message to %s (ID: %s), client removed", client.Role, client.ID)
            }
        }
    }

    if successfulSends == 0 {
        return fmt.Errorf("no staff members available to receive the message")
    }

    return nil
}