package websocket_service

import (
	"log"
	"time"

	model "english-ai-full/ecomm-api/websocket/websocker_model"

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
    send     chan *model.Message
    service  WebSocketService
    userID   int64
    userName string
    isGuest  bool
    closed   bool
}
func NewClient(conn *websocket.Conn, service WebSocketService, userID int64, userName string, isGuest bool) *Client {
    return &Client{
        conn:     conn,
        send:     make(chan *model.Message, 256),
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
        }
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        var message model.Message
        err := c.conn.ReadJSON(&message)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("ReadPump error: %v", err)
            }
            break
        }

        message.Sender = c.userID
        message.FromUser = c.userID
        message.Timestamp = time.Now()

        if message.ToUser != 0 { // Check for non-zero recipient
            if c.isGuest {
                c.service.SendMessageToGuest(message.FromUser, message.ToUser, message.Type, message.Content, message.TableID, message.OrderID)
            } else {
                c.service.SendMessageToUser(message.FromUser, message.ToUser, message.Type, message.Content, message.TableID, message.OrderID)
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
        }
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            err := c.conn.WriteJSON(message)
            if err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) getClientType() string {
    if c.isGuest {
        return "guest"
    }
    return "user"
}