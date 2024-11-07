package websocket_service

import (
	"log"
	"time"

	websocket_model "english-ai-full/ecomm-api/websocket/websocker_model"

	"github.com/gorilla/websocket"
)

const (
    writeWait      = 30 * time.Second  // Increased from 10s
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 65536
)


type Client struct {
    conn     *websocket.Conn
    send     chan *websocket_model.Message
    service  WebSocketService
    userID   string
    userName string
    isGuest  bool
    closed   bool
}
func NewClient(conn *websocket.Conn, service WebSocketService, userID string, userName string, isGuest bool) *Client {
    return &Client{
        conn:     conn,
        send:     make(chan *websocket_model.Message, 256),
        service:  service,
        userID:   userID,
        userName: userName,
        isGuest:  isGuest,
        closed:   false,
    }
}


func (c *Client) ReadPump() {
    defer func() {
        if !c.closed {
            c.closed = true
            c.service.UnregisterClient(c)
            c.conn.Close()
            log.Printf("ReadPump closed for user %s", c.userID)
        }
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    
    // Setup ping handler
    c.conn.SetPongHandler(func(string) error {
        log.Printf("Received pong from user %s", c.userID)
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        var message websocket_model.Message
        err := c.conn.ReadJSON(&message)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("ReadPump error for user %s: %v", c.userID, err)
            }
            break
        }

        log.Printf("Received message from user %s: %v", c.userID, message.Type)

        message.Sender = c.userID
        message.FromUser = c.userID
        message.Timestamp = time.Now()

        if message.ToUser != "" {
            if err := c.service.SendMessageToUser(
                message.FromUser,
                message.ToUser,
                message.Type,
                message.Content,
                message.TableID,
                message.OrderID,
            ); err != nil {
                log.Printf("Error sending direct message from user %s: %v", c.userID, err)
            }
        } else {
            c.service.BroadcastMessage(&message)
        }
    }
}



func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        if !c.closed {
            c.closed = true
            c.conn.Close()
            log.Printf("WritePump closed for user %s", c.userID)
        }
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                log.Printf("Send channel closed for user %s", c.userID)
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.conn.WriteJSON(message); err != nil {
                log.Printf("Error writing message to user %s: %v", c.userID, err)
                return
            }
            log.Printf("Successfully sent message to user %s", c.userID)

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Printf("Error sending ping to user %s: %v", c.userID, err)
                return
            }
            log.Printf("Sent ping to user %s", c.userID)
        }
    }
}