package websocket_service

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	websocket_model "english-ai-full/ecomm-api/websocket/websocker_model"
	"english-ai-full/ecomm-api/websocket/websocket_repository"

	"github.com/google/uuid"
)

type WebSocketService interface {
    RegisterClient(client *Client)
    UnregisterClient(client *Client)
    BroadcastMessage(message *websocket_model.Message)
    SendMessageToUser(fromUser, toUser string, messageType string, content interface{}, tableID, orderID string) error
    Run()
}

type webSocketService struct {
    clients    map[string]map[*Client]bool  // Map of userID to their clients
    broadcast  chan *websocket_model.Message
    register   chan *Client
    unregister chan *Client
    mutex      sync.Mutex
    repo       websocket_repository.MessageRepository
}
func NewWebSocketService(repo websocket_repository.MessageRepository) WebSocketService {
    return &webSocketService{
        clients:    make(map[string]map[*Client]bool),
        broadcast:  make(chan *websocket_model.Message),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        repo:       repo,
    }
}


func (s *webSocketService) RegisterClient(client *Client) {
    log.Printf("Registering client - UserID: %s, UserName: %s, Address: %p", 
        client.userID, client.userName, client)
    s.register <- client
}
func (s *webSocketService) SendToUser(userID string, message *websocket_model.Message) error {
    log.Printf("golang/ecomm-api/websocket/websocket_service/websocket_service.go SendToUser")
    s.mutex.Lock()
    defer s.mutex.Unlock()

    clients, exists := s.clients[userID]
    if !exists {
        return fmt.Errorf("no connected clients found for user %s", userID)
    }

    if len(clients) == 0 {
        return fmt.Errorf("user %s has no active connections", userID)
    }

    var lastError error
    for client := range clients {
        select {
        case client.send <- message:
            log.Printf("Message sent to client %s (%p)", client.userName, client)
        default:
            lastError = fmt.Errorf("failed to send message to client %s (%p)", 
                client.userName, client)
        
            close(client.send)
            delete(clients, client)
        }
    }

    return lastError
}

func (s *webSocketService) UnregisterClient(client *Client) {
	s.unregister <- client
}

func (s *webSocketService) BroadcastMessage(message *websocket_model.Message) {
    // Log the incoming message
    messageBytes, _ := json.Marshal(message)
    log.Printf("Broadcasting message: golang/ecomm-api/websocket/websocket_service/websocket_service.go %s", string(messageBytes))
    
    s.broadcast <- message
}
func (s *webSocketService) Run() {
    for {
        select {
        case client := <-s.register:
            s.mutex.Lock()
            // Initialize user's client map if it doesn't exist
            if _, exists := s.clients[client.userID]; !exists {
                s.clients[client.userID] = make(map[*Client]bool)
            }
            s.clients[client.userID][client] = true
            log.Printf("Client registered - UserID: %s, UserName: %s, Address: %p, Total clients for user: %d",
                client.userID, client.userName, client, len(s.clients[client.userID]))
            s.mutex.Unlock()

        case client := <-s.unregister:
            s.mutex.Lock()
            if clients, exists := s.clients[client.userID]; exists {
                if _, ok := clients[client]; ok {
                    delete(clients, client)
                    close(client.send)
                    log.Printf("Client unregistered - UserID: %s, UserName: %s, Address: %p, Remaining clients for user: %d",
                        client.userID, client.userName, client, len(clients))
                    // Clean up user entry if no more clients
                    if len(clients) == 0 {
                        delete(s.clients, client.userID)
                    }
                }
            }
            s.mutex.Unlock()

        case message := <-s.broadcast:
            // Handle broadcast messages (optional - you may want to keep or remove this)
            s.mutex.Lock()
            for _, clients := range s.clients {
                for client := range clients {
                    select {
                    case client.send <- message:
                        log.Printf("Broadcast message sent to client %s (%p)", 
                            client.userName, client)
                    default:
                        log.Printf("Failed to broadcast to client %s (%p)", 
                            client.userName, client)
                        close(client.send)
                        delete(clients, client)
                    }
                }
            }
            s.mutex.Unlock()

            if err := s.repo.SaveMessage(message); err != nil {
                log.Printf("Error saving message: %v", err)
            }
        }
    }
}

func (s *webSocketService) SendMessageToUser(fromUser, toUser string, messageType string, content interface{}, tableID, orderID string) error {
    message := &websocket_model.Message{
        ID:        uuid.New().String(),
        Type:      messageType,
        Content:   content,
        Sender:    fromUser,
        FromUser:  fromUser,
        ToUser:    toUser,
        Timestamp: time.Now(),
        TableID:   tableID,
        OrderID:   orderID,
    }
    
    log.Printf("Attempting to send message from %s to %s of type %s", fromUser, toUser, messageType)
    if err := s.SendToUser(toUser, message); err != nil {
        log.Printf("Error sending message: %v", err)
        return err
    }
    
    // Save message to repository
    if err := s.repo.SaveMessage(message); err != nil {
        log.Printf("Error saving message: %v", err)
        return err
    }
    
    return nil
}