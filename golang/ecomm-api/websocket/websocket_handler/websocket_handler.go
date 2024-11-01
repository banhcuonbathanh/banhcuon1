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
	"encoding/json"
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

// In your WebSocket handler
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Get user information from the request
    userID := r.URL.Query().Get("userId")    // or from your auth middleware
    userName := r.URL.Query().Get("userName") // optional

    if userID == "" {
        log.Printf("No userID provided for connection from %s", r.RemoteAddr)
        http.Error(w, "UserID is required", http.StatusBadRequest)
        return
    }

    conn, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection from %s: %v", r.RemoteAddr, err)
        return
    }

    log.Printf("New WebSocket connection established - UserID: %s, UserName: %s, Address: %s", 
        userID, userName, r.RemoteAddr)

    client := service.NewClient(conn, h.websocketService, userID, userName)
    h.websocketService.RegisterClient(client)

    go client.ReadPump()
    go client.WritePump()
}

// To send a message to a specific user


func (h *WebSocketHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var messageRequest struct {
        FromUser  string      `json:"fromUser"`
        ToUser    string      `json:"toUser"`
        Type      string      `json:"type"`
        Content   interface{} `json:"content"`
        TableID   string      `json:"table_id,omitempty"`
        OrderID   string      `json:"order_id,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
        http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    err := h.websocketService.SendMessageToUser(
        messageRequest.FromUser,
        messageRequest.ToUser,
        messageRequest.Type,
        messageRequest.Content,
        messageRequest.TableID,
        messageRequest.OrderID,
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "message": "Message sent successfully",
    })
}
