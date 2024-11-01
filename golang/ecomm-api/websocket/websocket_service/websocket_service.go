package websocket_service

import (
	"encoding/json"
	"log"
	"sync"

	websocket_model "english-ai-full/ecomm-api/websocket/websocker_model"
	"english-ai-full/ecomm-api/websocket/websocket_repository"
)

type WebSocketService interface {
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *websocket_model.Message)
	Run()
}

type webSocketService struct {
	clients    map[*Client]bool
	broadcast  chan *websocket_model.Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
	repo       websocket_repository.MessageRepository
}

func NewWebSocketService(repo websocket_repository.MessageRepository) WebSocketService {
	return &webSocketService{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *websocket_model.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		repo:       repo,
	}
}

func (s *webSocketService) RegisterClient(client *Client) {
    log.Printf("golang/ecomm-api/websocket/websocket_service/websocket_service.go RegisterClient %v", client)

	s.register <- client
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
    log.Printf("golang/ecomm-api/websocket/websocket_service/websocket_service.go Run %v", s.clients)
    for {
        select {
        case client := <-s.register:
            log.Printf("New client registered: %p", client)
            s.mutex.Lock()
            s.clients[client] = true
            s.mutex.Unlock()
            
        case client := <-s.unregister:
            log.Printf("Client unregistered: %p", client)
            s.mutex.Lock()
            if _, ok := s.clients[client]; ok {
                delete(s.clients, client)
                close(client.send)
            }
            s.mutex.Unlock()
            
        case message := <-s.broadcast:
            log.Printf("Broadcasting to %d clients", len(s.clients))
            s.mutex.Lock()
            for client := range s.clients {
                select {
                case client.send <- message:
                    log.Printf("Message sent to client: %p", client)
                default:
                    log.Printf("Failed to send message to client: %p", client)
                    close(client.send)
                    delete(s.clients, client)
                }
            }
            s.mutex.Unlock()

            if err := s.repo.SaveMessage(message); err != nil {
                log.Printf("Error saving message: %v", err)
            }
        }
    }
}