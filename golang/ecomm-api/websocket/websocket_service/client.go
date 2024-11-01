package websocket_service

import (
	"log"
	"time"

	websocket_model "english-ai-full/ecomm-api/websocket/websocker_model"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)


type Client struct {
    conn     *websocket.Conn
    send     chan *websocket_model.Message
    service  WebSocketService
    userID   string        // Add user identification
    userName string        // Optional: Add user name for logging
}


func NewClient(conn *websocket.Conn, service WebSocketService, userID string, userName string) *Client {
    return &Client{
        conn:     conn,
        send:     make(chan *websocket_model.Message, 256),
        service:  service,
        userID:   userID,
        userName: userName,
    }
}



func (c *Client) ReadPump() {
	log.Printf("golang/ecomm-api/websocket/websocket_service/client.go ReadPump")
	
	defer func() {
		c.service.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var message websocket_model.Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		c.service.BroadcastMessage(&message)
	}
}



func (c *Client) WritePump() {
	log.Printf("golang/ecomm-api/websocket/websocket_service/client.go WritePump")
	
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            // Add timestamp if not present
            if message.Timestamp.IsZero() {
                message.Timestamp = time.Now()
            }

            err := c.conn.WriteJSON(message)
            if err != nil {
                log.Printf("Error writing message to client: %v", err)
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

