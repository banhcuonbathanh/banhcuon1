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
            ReadBufferSize:  4096,
            WriteBufferSize: 4096,
            EnableCompression: true,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
            HandshakeTimeout: 30 * time.Second, // Increased timeout
        },
        websocketService: websocketService,
    }
}


// In your WebSocket handler
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    userName := r.Context().Value("userName").(string)
    isGuest := r.Context().Value("isGuest").(bool)

    log.Printf("WebSocket connection request from UserID: %s, UserName: %s, IsGuest: %t", userID, userName, isGuest)

    // Set response headers for WebSocket
    headers := http.Header{}

    conn, err := h.upgrader.Upgrade(w, r, headers)
    if err != nil {
        log.Printf("Failed to upgrade connection from %s: %v", r.RemoteAddr, err)
        return
    }

    // Enable compression if available
    conn.EnableWriteCompression(true)

    client := service.NewClient(conn, h.websocketService, userID, userName, isGuest)
    h.websocketService.RegisterClient(client)

    log.Printf("WebSocket client registered - UserID: %s, UserName: %s, Address: %p", userID, userName, client)

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