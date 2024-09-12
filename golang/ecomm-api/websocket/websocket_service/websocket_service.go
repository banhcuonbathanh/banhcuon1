package websocket_service

import (
	"log"
	"sync"

	 "english-ai-full/ecomm-api/websocket/websocker_model"
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
	s.register <- client
}

func (s *webSocketService) UnregisterClient(client *Client) {
	s.unregister <- client
}

func (s *webSocketService) BroadcastMessage(message *websocket_model.Message) {
	s.broadcast <- message
}

func (s *webSocketService) Run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mutex.Unlock()
		case message := <-s.broadcast:
			s.mutex.Lock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mutex.Unlock()

			// Save message to repository
			if err := s.repo.SaveMessage(message); err != nil {
				log.Printf("Error saving message: %v", err)
			}
		}
	}
}
