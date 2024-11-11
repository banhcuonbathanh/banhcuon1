package websocket_route

import (
	"encoding/json"

	"english-ai-full/token"
	"net/http"

	"github.com/go-chi/chi"

	Handler "english-ai-full/ecomm-api/websocket/websocket_handler"
	order "english-ai-full/quanqr/order"

	model "english-ai-full/ecomm-api/websocket/websocker_model"
)

type WebSocketRoutes struct {
    wsHandler    *Handler.WebSocketHandler
    TokenMaker   *token.JWTMaker
    orderHandler *order.OrderHandlerController
}


func NewWebSocketRoutes(wsHandler  *Handler.WebSocketHandler, orderHandler  *order.OrderHandlerController, secretKey string) *WebSocketRoutes {
    return &WebSocketRoutes{
        wsHandler:    wsHandler,
        TokenMaker: token.NewJWTMaker(secretKey),
        orderHandler: orderHandler,
    }
}

// TokenMaker interface for authentication


func (wr *WebSocketRoutes) RegisterWebSocketRoutes(r *chi.Mux) {
    r.Route("/ws", func(r chi.Router) {
        // Authentication middleware
        // r.Use(wr.authenticateWS)
        
        // WebSocket connection endpoint
        r.Get("/", wr.wsHandler.HandleWebSocket)
        
        // Handle different message types
        r.HandleFunc("/", wr.handleWSConnection)
    })
}



func (wr *WebSocketRoutes) handleWSConnection(w http.ResponseWriter, r *http.Request) {
    // Extract user info from context
    userID := r.Context().Value("userID").(int64)
    isGuest := r.Context().Value("isGuest").(bool)

    // Handle different message types
    messageType := r.URL.Query().Get("type")
    switch messageType {
    case "ORDER":
        wr.handleOrderMessage(w, r, userID, isGuest)
    case "BROADCAST":
        wr.handleBroadcastMessage(w, r, userID)
    default:
        wr.handleDirectMessage(w, r, userID, isGuest)
    }
}

func (wr *WebSocketRoutes) handleOrderMessage(w http.ResponseWriter, r *http.Request, userID int64, isGuest bool) {
    var orderContent order.CreateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderContent); err != nil {
        http.Error(w, "Invalid order content", http.StatusBadRequest)
        return
    }

    message := &model.Message{
        Type:    "ORDER",
        Content: orderContent,
        Sender:  userID,
        ToUser:  *orderContent.GuestID, // Changed to dereference pointer
    }

    // Check if GuestID is not nil instead of comparing with empty string
    if orderContent.GuestID != nil {
        if isGuest {
            // Use the websocketService through the handler instead
            err := wr.wsHandler.websocketService.SendMessageToGuest(
                userID,
                *orderContent.GuestID, // Dereference the pointer
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
        // Use the websocketService through the handler
        wr.wsHandler.websocketService.BroadcastMessage(message)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "Order message sent"})
}

func (wr *WebSocketRoutes) handleDirectMessage(w http.ResponseWriter, r *http.Request, userID int64, isGuest bool) {
    var msgContent struct {
        ToUser  int64       `json:"to_user"`    // Changed to int64 to match Message type
        Content interface{} `json:"content"`
        TableID string      `json:"table_id"`
        OrderID string      `json:"order_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&msgContent); err != nil {
        http.Error(w, "Invalid message content", http.StatusBadRequest)
        return
    }

    if isGuest {
        err := wr.wsHandler.websocketService.SendMessageToGuest(
            userID,
            msgContent.ToUser,
            "DIRECT",
            msgContent.Content,
            msgContent.TableID,
            msgContent.OrderID,
        )
        if err != nil {
            http.Error(w, "Failed to send message", http.StatusInternalServerError)
            return
        }
    } else {
        err := wr.wsHandler.websocketService.SendMessageToUser(
            userID,
            msgContent.ToUser,
            "DIRECT",
            msgContent.Content,
            msgContent.TableID,
            msgContent.OrderID,
        )
        if err != nil {
            http.Error(w, "Failed to send message", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "Message sent"})
}
func (wr *WebSocketRoutes) handleBroadcastMessage(w http.ResponseWriter, r *http.Request, userID int64) {
    var msgContent struct {
        Content interface{} `json:"content"`
    }

    if err := json.NewDecoder(r.Body).Decode(&msgContent); err != nil {
        http.Error(w, "Invalid message content", http.StatusBadRequest)
        return
    }

    message := &model.Message{
        Type:    "BROADCAST",
        Content: msgContent.Content,
        Sender:  userID,
    }

    wr.wsHandler.BroadcastMessage(message)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "Broadcast sent"})
}