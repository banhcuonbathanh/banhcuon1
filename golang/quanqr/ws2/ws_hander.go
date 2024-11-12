package ws2

import "log"



type WebSocketHandler struct {
    hub *Hub
}

func NewWebSocketHandler(messageHandler MessageHandler) *WebSocketHandler {
	log.Println("golang/quanqr/ws2/ws_hander.go")
    return &WebSocketHandler{
        hub: NewHub(messageHandler),
    }
}