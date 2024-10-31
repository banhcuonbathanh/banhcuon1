// package websockethandler

// import (
// 	"log"
// 	"net/http"

// 	service "english-ai-full/ecomm-api/websocket/websocket_service"
// 	"github.com/gorilla/websocket"
// )

// type WebSocketHandler struct {
// 	upgrader        websocket.Upgrader
// 	websocketService service.WebSocketService
// }

// func NewWebSocketHandler(websocketService service.WebSocketService) *WebSocketHandler {
// 	return &WebSocketHandler{
// 		upgrader: websocket.Upgrader{
// 			ReadBufferSize:  1024,
// 			WriteBufferSize: 1024,
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true
// 			},
// 		},
// 		websocketService: websocketService,
// 	}
// }

// func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := h.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("Failed to upgrade connection: %v", err)
// 		return
// 	}

// 	client := service.NewClient(conn, h.websocketService)
// 	h.websocketService.RegisterClient(client)

// 	go client.ReadPump()
// 	go client.WritePump()
// }

package websocket_handler

import (
	"log"
	"net/http"
	"time"

	service "english-ai-full/ecomm-api/websocket/websocket_service"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
    upgrader        websocket.Upgrader
    websocketService service.WebSocketService
}

func NewWebSocketHandler(websocketService service.WebSocketService) *WebSocketHandler {
    return &WebSocketHandler{
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                // In production, you might want to be more restrictive
                return true
            },
            HandshakeTimeout: 10 * time.Second,
        },
        websocketService: websocketService,
    }
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Log connection attempt
    log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)

    conn, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection from %s: %v", r.RemoteAddr, err)
        return
    }

    log.Printf("New WebSocket connection established from %s", r.RemoteAddr)

    client := service.NewClient(conn, h.websocketService)
    h.websocketService.RegisterClient(client)

    // Start client message pumps
    go client.ReadPump()
    go client.WritePump()
}
