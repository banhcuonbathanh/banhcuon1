package websocket_handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	service "english-ai-full/ecomm-api/websocket/websocket_service"

	"github.com/gorilla/websocket"

    model "english-ai-full/ecomm-api/websocket/websocker_model"
    order "english-ai-full/quanqr/order"
    
)

type WebSocketHandler struct {
    upgrader        websocket.Upgrader
    websocketService service.WebSocketService
}

func NewWebSocketHandler(websocketService service.WebSocketService) *WebSocketHandler {
    return &WebSocketHandler{
        upgrader: websocket.Upgrader{
            ReadBufferSize:     4096,
            WriteBufferSize:    4096,
            EnableCompression:  true,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
            HandshakeTimeout: 30 * time.Second,
        },
        websocketService: websocketService,
    }
}


// In your WebSocket handler
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    userName := r.Context().Value("userName").(string)
    isGuest := r.Context().Value("isGuest").(bool)
    
    conn, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    
    conn.EnableWriteCompression(true)
    
    client := service.NewClient(conn, h.websocketService, userID, userName, isGuest)
    h.websocketService.RegisterClient(client)
    
    go client.ReadPump()
    go client.WritePump()
}
// To send a message to a specific user

func (h *WebSocketHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
    var messageRequest struct {
        FromUserID string      `json:"fromUserId"`
        ToUserID   string      `json:"toUserId"`
        Type       string      `json:"type"`
        Content    interface{} `json:"content"`
        TableID    string      `json:"tableId,omitempty"`
        OrderID    string      `json:"orderId,omitempty"`
        IsGuest    bool        `json:"isGuest"` // Add this field to specify if sending to a guest
    }

    if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        log.Printf("Error decoding message request: %v", err)
        return
    }

    var err error
    if messageRequest.IsGuest {
        // Send message to guest
        err = h.websocketService.SendMessageToGuest(
            messageRequest.FromUserID,
            messageRequest.ToUserID,
            messageRequest.Type,
            messageRequest.Content,
            messageRequest.TableID,
            messageRequest.OrderID,
        )
    } else {
        // Send message to user
        err = h.websocketService.SendMessageToUser(
            messageRequest.FromUserID,
            messageRequest.ToUserID,
            messageRequest.Type,
            messageRequest.Content,
            messageRequest.TableID,
            messageRequest.OrderID,
        )
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Error sending message: %v", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "Message sent successfully",
    })
}

func (h *WebSocketHandler) HandleOrderMessage(w http.ResponseWriter, r *http.Request, userID string, isGuest bool) {
    var orderContent order.CreateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderContent); err != nil {
        http.Error(w, "Invalid order content", http.StatusBadRequest)
        return
    }

    message := &model.Message{
        Type:    "ORDER",
        Content: orderContent,
        Sender:  userID,
        ToUser:  orderContent.GuestID,
    }

    if orderContent.GuestID != nil {
        if isGuest {
            err := h.websocketService.SendMessageToGuest(
                userID,
                string(*orderContent.GuestID),
                "ORDER",
                orderContent,
                orderContent.TableToken,
                "",
            )
            if err != nil {
                http.Error(w, "Failed to send order message", http.StatusInternalServerError)
                return
            }
        }
    } else {
        h.websocketService.BroadcastMessage(message)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "Order message sent"})
}