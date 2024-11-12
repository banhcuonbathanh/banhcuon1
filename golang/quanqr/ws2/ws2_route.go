package ws2

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type WebSocketRouter struct {
    hub           *Hub
    deliveryQueue chan DeliveryUpdate
}

type DeliveryUpdate struct {
    Action     string      `json:"action"`
    DeliveryID string      `json:"deliveryId"`
    Payload    interface{} `json:"payload"`
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // In production, implement proper origin checking
    },
}

func NewWebSocketRouter(h *Hub) *WebSocketRouter {
    router := &WebSocketRouter{
        hub:           h,
        deliveryQueue: make(chan DeliveryUpdate, 100),
    }
    
    // Start delivery update processor
    go router.processDeliveryUpdates()
    
    return router
}

// RegisterRoutes registers all WebSocket routes
func (wr *WebSocketRouter) RegisterRoutes(r chi.Router) {
    r.Route("/ws", func(r chi.Router) {
		log.Println("golang/quanqr/ws2/ws2_route.go user RegisterRoutes")
        r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
            wr.handleWebSocket(w, r, RoleUser)
        })

        // Guest connections
        r.Get("/guest/{id}", func(w http.ResponseWriter, r *http.Request) {
			log.Println("golang/quanqr/ws2/ws2_route.go guest RegisterRoutes")
            wr.handleWebSocket(w, r, RoleGuest)
        })

        // Kitchen staff connections
        r.Get("/kitchen/{id}", func(w http.ResponseWriter, r *http.Request) {
			log.Println("golang/quanqr/ws2/ws2_route.go kitchen RegisterRoutes")
            wr.handleWebSocket(w, r, RoleKitchen)
        })

        // Employee connections
        r.Get("/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
			log.Println("golang/quanqr/ws2/ws2_route.go employee RegisterRoutes")
            wr.handleWebSocket(w, r, RoleEmployee)
        })

        // Admin connections
        r.Get("/admin/{id}", func(w http.ResponseWriter, r *http.Request) {
			log.Println("golang/quanqr/ws2/ws2_route.go admin RegisterRoutes")
            wr.handleWebSocket(w, r, RoleAdmin)
        })
    })
}

func (wr *WebSocketRouter) handleWebSocket(w http.ResponseWriter, r *http.Request, role Role) {
	log.Println("golang/quanqr/ws2/ws2_route.go handleWebSocket")
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }

    // Extract additional parameters
    userToken := r.URL.Query().Get("token")
    tableToken := r.URL.Query().Get("tableToken")

    client := &Client{
        Hub:    wr.hub,
        Conn:   conn,
        Send:   make(chan []byte, 256),
        Role:   role,
        ID:     chi.URLParam(r, "id"),
        RoomID: r.URL.Query().Get("roomId"),
        UserData: map[string]interface{}{
            "token":      userToken,
            "tableToken": tableToken,
        },
    }
	// log.Printf("golang/quanqr/ws2/ws2_route.go handleWebSocket client %+v", client)

    // Register client
    client.Hub.Register <- client

    // Start read/write pumps
    go client.ReadPump()
    go client.WritePump()
}

// BroadcastDeliveryUpdate sends delivery updates to connected clients
func (wr *WebSocketRouter) BroadcastDeliveryUpdate(action string, deliveryID string, payload interface{}) {
    wr.deliveryQueue <- DeliveryUpdate{
        Action:     action,
        DeliveryID: deliveryID,
        Payload:    payload,
    }
}

// processDeliveryUpdates handles the delivery update queue
func (wr *WebSocketRouter) processDeliveryUpdates() {
    for update := range wr.deliveryQueue {
        message := Message{
            Type:    "delivery",
            Action:  update.Action,
            Payload: update.Payload,
            Role:    RoleEmployee, // Default to employee, can be modified based on needs
        }

        // Marshal the message
        data, err := json.Marshal(message)
        if err != nil {
            log.Printf("Error marshaling delivery update: %v", err)
            continue
        }

        // Broadcast to all relevant clients
        wr.hub.Broadcast <- data
    }
}

// Helper method to get connected clients count
func (wr *WebSocketRouter) GetConnectedClientsCount() map[Role]int {
    counts := make(map[Role]int)
    wr.hub.mu.Lock()
    defer wr.hub.mu.Unlock()

    for client := range wr.hub.Clients {
        counts[client.Role]++
    }
    return counts
}
