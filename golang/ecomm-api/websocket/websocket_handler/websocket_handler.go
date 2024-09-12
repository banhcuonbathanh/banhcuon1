package websockethandler

import (
	"log"
	"net/http"

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
				return true
			},
		},
		websocketService: websocketService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := service.NewClient(conn, h.websocketService)
	h.websocketService.RegisterClient(client)

	go client.ReadPump()
	go client.WritePump()
}
