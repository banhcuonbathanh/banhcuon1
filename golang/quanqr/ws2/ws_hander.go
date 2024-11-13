package ws2

import "log"



type WebSocketHandler struct {
    hub *Hub
}

func NewWebSocketHandler(   CombinedMessageHandler *CombinedMessageHandler ) *WebSocketHandler {
	log.Println("golang/quanqr/ws2/ws_hander.go")
    return &WebSocketHandler{
        hub: NewHub(CombinedMessageHandler),
    }
}