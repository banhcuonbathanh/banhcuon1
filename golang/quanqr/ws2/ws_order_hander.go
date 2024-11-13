package ws2

import (
	"bytes"
	"context"
	"encoding/json"
	order "english-ai-full/quanqr/order"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// type OrderMessageHandler struct {
//     DefaultMessageHandler
//     orderHandler *order.OrderHandlerController
// }
// new
type OrderMessageHandler struct {
    DefaultMessageHandler
    orderHandler *order.OrderHandlerController
    broadcaster  *Broadcaster
}
func NewOrderMessageHandler(orderHandler *order.OrderHandlerController) *OrderMessageHandler {
    log.Println("golang/quanqr/ws2/ws_order_hander.go NewOrderMessageHandler")
    return &OrderMessageHandler{
        orderHandler: orderHandler,
    }
}

// Handle extends the default handler with order-specific logic
func (h *OrderMessageHandler) Handle(c *Client, msg Message) {
    log.Printf("BEGIN OrderMessageHandler.Handle - Type: %s, Action: %s", msg.Type, msg.Action)
    defer log.Printf("END OrderMessageHandler.Handle")
    
   
        switch msg.Type {
        case "direct":
            if msg.Action == "order" {
                log.Printf("Handling order message to staff")
                h.handleOrderMessageToStaff(c, msg)
            } else {
                log.Printf("Handling direct message")
                h.handleDirectMessage(c, msg)
            }
    case "order":
        log.Printf("Handling order message")
        h.handleOrderMessage(c, msg)
    case "notification":
        log.Printf("Handling notification message")
        h.handleNotificationMessage(c, msg)
    case "status_update":
        log.Printf("Handling status update")
        h.handleStatusUpdate(c, msg)
    default:
        log.Printf("Using default handler")
        h.DefaultMessageHandler.HandleMessage(c, msg)
    }
}

func (h *OrderMessageHandler) handleDirectMessage(c *Client, msg Message) {
    log.Printf("golang/quanqr/ws2/ws_order_hander.go BEGIN handleDirectMessage")
    defer log.Printf("golang/quanqr/ws2/ws_order_hander.go END handleDirectMessage")
    
    var directMsg DirectMessage
    data, _ := json.Marshal(msg.Payload)
    if err := json.Unmarshal(data, &directMsg); err != nil {
        log.Printf("error unmarshaling direct message: %v", err)
        return
    }

    log.Printf("golang/quanqr/ws2/ws_order_hander.go Sending direct message from %s to %s", directMsg.FromUserID, directMsg.ToUserID)
    log.Printf("golang/quanqr/ws2/ws_order_hander.g directMsg.Payload %v", directMsg.Payload)

    // Convert payload map to CreateOrderRequestType
    payloadMap, ok := directMsg.Payload.(map[string]interface{})
    if !ok {
        log.Printf("error: payload is not a map[string]interface{}")
        return
    }

    // Helper functions for safe type conversions
    safeFloat64 := func(v interface{}) float64 {
        if v == nil {
            return 0
        }
        switch i := v.(type) {
        case float64:
            return i
        case int:
            return float64(i)
        case int64:
            return float64(i)
        default:
            return 0
        }
    }

    safeBool := func(v interface{}) bool {
        if v == nil {
            return false
        }
        b, ok := v.(bool)
        if !ok {
            return false
        }
        return b
    }

    safeString := func(v interface{}) string {
        if v == nil {
            return ""
        }
        s, ok := v.(string)
        if !ok {
            return ""
        }
        return s
    }

    safeInt64 := func(v interface{}) int64 {
        if v == nil {
            return 0
        }
        switch i := v.(type) {
        case float64:
            return int64(i)
        case int:
            return int64(i)
        case int64:
            return i
        default:
            return 0
        }
    }

    // Helper function to get value from map with multiple possible keys
    getMapValue := func(m map[string]interface{}, keys ...string) interface{} {
        for _, key := range keys {
            if val, exists := m[key]; exists && val != nil {
                return val
            }
        }
        return nil
    }

    // Parse dish items
    var dishItems []order.OrderDish
    if rawDishItems, ok := getMapValue(payloadMap, "dish_items", "dishItems").([]interface{}); ok {
        for _, item := range rawDishItems {
            if dishMap, ok := item.(map[string]interface{}); ok {
                dishItems = append(dishItems, order.OrderDish{
                    DishID:   int64(safeFloat64(getMapValue(dishMap, "dish_id", "dishId"))),
                    Quantity: int64(safeFloat64(getMapValue(dishMap, "quantity"))),
                })
            }
        }
    }

    // Parse set items
    var setItems []order.OrderSet
    if rawSetItems, ok := getMapValue(payloadMap, "set_items", "setItems").([]interface{}); ok {
        for _, item := range rawSetItems {
            if setMap, ok := item.(map[string]interface{}); ok {
                setItems = append(setItems, order.OrderSet{
                    SetID:    int64(safeFloat64(getMapValue(setMap, "set_id", "setId"))),
                    Quantity: int64(safeFloat64(getMapValue(setMap, "quantity"))),
                })
            }
        }
    }

    isGuest := safeBool(getMapValue(payloadMap, "is_guest", "isGuest"))
    
    // Create order request with safe conversions
    orderReq := order.CreateOrderRequestType{
        TableNumber:    int64(safeFloat64(getMapValue(payloadMap, "table_number", "tableNumber"))),
        TotalPrice:     int32(safeFloat64(getMapValue(payloadMap, "total_price", "totalPrice"))),
        BowChili:       int64(safeFloat64(getMapValue(payloadMap, "bow_chili", "bowChili"))),
        BowNoChili:     int64(safeFloat64(getMapValue(payloadMap, "bow_no_chili", "bowNoChili"))),
        TakeAway:       safeBool(getMapValue(payloadMap, "take_away", "takeAway")),
        ChiliNumber:    int64(safeFloat64(getMapValue(payloadMap, "chili_number", "chiliNumber"))),
        TableToken:     safeString(getMapValue(payloadMap, "table_token", "tableToken", "Table_token")), // Added Table_token as fallback
        OrderName:      safeString(getMapValue(payloadMap, "order_name", "orderName")),
        OrderHandlerID: safeInt64(getMapValue(payloadMap, "order_handler_id", "orderHandlerId")), // Added OrderHandlerID
        IsGuest:        safeBool(getMapValue(payloadMap, "is_guest", "isGuest")),
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
        Status:         "pending",
        DishItems:      dishItems,
        SetItems:       setItems,
    }

    // Modified ID assignment logic with snake_case support
    if isGuest {
        guestId := safeInt64(getMapValue(payloadMap, "guest_id", "guestId"))
        if guestId != 0 {
            orderReq.GuestID = guestId
        }
        orderReq.UserID = 0
    } else {
        userId := safeInt64(getMapValue(payloadMap, "user_id", "userId"))
        if userId != 0 {
            orderReq.UserID = userId
        }
        orderReq.GuestID = 0
    }

    // Validate required fields
    if orderReq.IsGuest && orderReq.GuestID == 0 {
        log.Printf("error: guest order requires guest_id")
        return
    }

    if !orderReq.IsGuest && orderReq.UserID == 0 {
        log.Printf("error: user order requires user_id")
        return
    }

    // Log the created order request for debugging
    log.Printf("golang/quanqr/ws2/ws_order_hander.go Created order request: %+v", orderReq)

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Serialize the order request to JSON
    orderReqJSON, err := json.Marshal(orderReq)
    if err != nil {
        log.Printf("error marshaling order request: %v", err)
        return
    }
    bodyReader := bytes.NewReader(orderReqJSON)
    // Create the request with the serialized order data
    r := &http.Request{
        Method: "POST",  // Add HTTP method
        URL: &url.URL{}, // Add empty URL
        Header: make(http.Header),
        Body: io.NopCloser(bodyReader),
        GetBody: func() (io.ReadCloser, error) { // This allows the body to be read multiple times
            r := bytes.NewReader(orderReqJSON)
            return io.NopCloser(r), nil
        },
        ContentLength: int64(len(orderReqJSON)),
    }
    r.Header.Set("Content-Type", "application/json")
    // Create the order
    w := NewResponseWriter()
    h.orderHandler.CreateOrder2(w, r.WithContext(ctx))

    if w.StatusCode != http.StatusCreated {
        log.Printf("error creating order: received status code %d", w.StatusCode)
        return
    }

    if err := c.Hub.SendDirectMessage(directMsg.FromUserID, directMsg.ToUserID, directMsg.Type, directMsg.Action, directMsg.Payload); err != nil {
        log.Printf("error sending direct message: %v", err)
    }
}
type ResponseWriter struct {
    HeaderMap    http.Header
    Body         bytes.Buffer
    StatusCode   int
}

func NewResponseWriter() *ResponseWriter {
    return &ResponseWriter{
        HeaderMap: make(http.Header),
        StatusCode: http.StatusOK,
    }
}

func (w *ResponseWriter) Header() http.Header {
    return w.HeaderMap
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
    return w.Body.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
    w.StatusCode = statusCode
}


func (h *OrderMessageHandler) handleOrderMessage(c *Client, msg Message) {
    var order OrderMessage
    data, _ := json.Marshal(msg.Payload)
    if err := json.Unmarshal(data, &order); err != nil {
        log.Printf("error unmarshaling order: %v", err)
        return
    }


    notification := Message{
        Type:    "notification",
        Action:  "new_order",
        Payload: order,
        Role:    RoleKitchen,
    }
    
    data, _ = json.Marshal(notification)
    c.Hub.Broadcast <- data
}

// new

func (h *OrderMessageHandler) SetBroadcaster(b *Broadcaster) {
    h.broadcaster = b
}

func (h *OrderMessageHandler) handleOrderMessageToStaff(c *Client, msg Message) {
    log.Printf("BEGIN handleOrderMessageToStaff")
    defer log.Printf("END handleOrderMessageToStaff")

    // Extract the direct message
    data, _ := json.Marshal(msg.Payload)
    var directMsg DirectMessage
    if err := json.Unmarshal(data, &directMsg); err != nil {
        log.Printf("error unmarshaling direct message: %v", err)
        return
    }

    // First, create the order in the database
    if err := h.createOrder(directMsg.Payload); err != nil {
        log.Printf("Error creating order: %v", err)
        return
    }

    // After successful order creation, broadcast to staff
    staffMsg := Message{
        Type:    "order",
        Action:  "new_order",
        Payload: directMsg.Payload,
    }

    // Broadcast to all staff members
    if err := c.Hub.BroadcastToStaff(directMsg.FromUserID, staffMsg); err != nil {
        log.Printf("Error broadcasting order to staff: %v", err)
        return
    }

    log.Printf("Successfully created order and broadcasted to staff")
}

// Optional: Add a helper method to Hub to make it more semantic
func (h *Hub) BroadcastOrderToStaff(msg Message) error {
    // This is just a wrapper around BroadcastToStaff for better semantics
    return h.BroadcastToStaff("system", msg)
}

// 

func (h *OrderMessageHandler) createOrder(payload interface{}) error {
    log.Printf("BEGIN createOrder")
    defer log.Printf("END createOrder")

    // Convert payload to map
    payloadMap, ok := payload.(map[string]interface{})
    if !ok {
        return fmt.Errorf("error: payload is not a map[string]interface{}")
    }

    // Helper functions for safe type conversions
    safeFloat64 := func(v interface{}) float64 {
        if v == nil {
            return 0
        }
        switch i := v.(type) {
        case float64:
            return i
        case int:
            return float64(i)
        case int64:
            return float64(i)
        default:
            return 0
        }
    }

    safeBool := func(v interface{}) bool {
        if v == nil {
            return false
        }
        b, ok := v.(bool)
        if !ok {
            return false
        }
        return b
    }

    safeString := func(v interface{}) string {
        if v == nil {
            return ""
        }
        s, ok := v.(string)
        if !ok {
            return ""
        }
        return s
    }

    safeInt64 := func(v interface{}) int64 {
        if v == nil {
            return 0
        }
        switch i := v.(type) {
        case float64:
            return int64(i)
        case int:
            return int64(i)
        case int64:
            return i
        default:
            return 0
        }
    }

    // Helper function to get value from map with multiple possible keys
    getMapValue := func(m map[string]interface{}, keys ...string) interface{} {
        for _, key := range keys {
            if val, exists := m[key]; exists && val != nil {
                return val
            }
        }
        return nil
    }

    // Parse dish items
    var dishItems []order.OrderDish
    if rawDishItems, ok := getMapValue(payloadMap, "dish_items", "dishItems").([]interface{}); ok {
        for _, item := range rawDishItems {
            if dishMap, ok := item.(map[string]interface{}); ok {
                dishItems = append(dishItems, order.OrderDish{
                    DishID:   int64(safeFloat64(getMapValue(dishMap, "dish_id", "dishId"))),
                    Quantity: int64(safeFloat64(getMapValue(dishMap, "quantity"))),
                })
            }
        }
    }

    // Parse set items
    var setItems []order.OrderSet
    if rawSetItems, ok := getMapValue(payloadMap, "set_items", "setItems").([]interface{}); ok {
        for _, item := range rawSetItems {
            if setMap, ok := item.(map[string]interface{}); ok {
                setItems = append(setItems, order.OrderSet{
                    SetID:    int64(safeFloat64(getMapValue(setMap, "set_id", "setId"))),
                    Quantity: int64(safeFloat64(getMapValue(setMap, "quantity"))),
                })
            }
        }
    }

    isGuest := safeBool(getMapValue(payloadMap, "is_guest", "isGuest"))
    
    // Create order request with safe conversions
    orderReq := order.CreateOrderRequestType{
        TableNumber:    int64(safeFloat64(getMapValue(payloadMap, "table_number", "tableNumber"))),
        TotalPrice:     int32(safeFloat64(getMapValue(payloadMap, "total_price", "totalPrice"))),
        BowChili:       int64(safeFloat64(getMapValue(payloadMap, "bow_chili", "bowChili"))),
        BowNoChili:     int64(safeFloat64(getMapValue(payloadMap, "bow_no_chili", "bowNoChili"))),
        TakeAway:       safeBool(getMapValue(payloadMap, "take_away", "takeAway")),
        ChiliNumber:    int64(safeFloat64(getMapValue(payloadMap, "chili_number", "chiliNumber"))),
        TableToken:     safeString(getMapValue(payloadMap, "table_token", "tableToken", "Table_token")),
        OrderName:      safeString(getMapValue(payloadMap, "order_name", "orderName")),
        OrderHandlerID: safeInt64(getMapValue(payloadMap, "order_handler_id", "orderHandlerId")),
        IsGuest:        isGuest,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
        Status:         "pending",
        DishItems:      dishItems,
        SetItems:       setItems,
    }

    // Modified ID assignment logic with snake_case support
    if isGuest {
        guestId := safeInt64(getMapValue(payloadMap, "guest_id", "guestId"))
        if guestId != 0 {
            orderReq.GuestID = guestId
        }
        orderReq.UserID = 0
    } else {
        userId := safeInt64(getMapValue(payloadMap, "user_id", "userId"))
        if userId != 0 {
            orderReq.UserID = userId
        }
        orderReq.GuestID = 0
    }

    // Validate required fields
    if orderReq.IsGuest && orderReq.GuestID == 0 {
        return fmt.Errorf("error: guest order requires guest_id")
    }

    if !orderReq.IsGuest && orderReq.UserID == 0 {
        return fmt.Errorf("error: user order requires user_id")
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Serialize the order request to JSON
    orderReqJSON, err := json.Marshal(orderReq)
    if err != nil {
        return fmt.Errorf("error marshaling order request: %v", err)
    }

    bodyReader := bytes.NewReader(orderReqJSON)
    // Create the request with the serialized order data
    r := &http.Request{
        Method: "POST",
        URL: &url.URL{},
        Header: make(http.Header),
        Body: io.NopCloser(bodyReader),
        GetBody: func() (io.ReadCloser, error) {
            r := bytes.NewReader(orderReqJSON)
            return io.NopCloser(r), nil
        },
        ContentLength: int64(len(orderReqJSON)),
    }
    r.Header.Set("Content-Type", "application/json")

    // Create the order
    w := NewResponseWriter()
    h.orderHandler.CreateOrder2(w, r.WithContext(ctx))

    if w.StatusCode != http.StatusCreated {
        return fmt.Errorf("error creating order: received status code %d", w.StatusCode)
    }

    return nil
}
