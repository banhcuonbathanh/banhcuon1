package ws2

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
    Hub      *Hub
    Conn     *websocket.Conn
    Send     chan []byte
    Role     Role
    ID       string
    RoomID   string
    UserData interface{}
}

func (c *Client) ReadPump() {

	log.Printf("golang/quanqr/ws2/ws2_client.go ReadPump 1" )
    defer func() {
        c.Hub.Unregister <- c
        c.Conn.Close()
    }()

    c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
	log.Printf("golang/quanqr/ws2/ws2_client.go ReadPump 2" )
    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        // log.Printf("ReadPump received message from client %s: %s", c.ID, string(message))
        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Printf("error unmarshaling message: %v", err)
            continue
        }

        c.Hub.MessageHandler.Handle(c, msg)
    }
}

func (c *Client) WritePump() {

	log.Printf("golang/quanqr/ws2/ws2_client.go WritePump 1" )
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            log.Printf("WritePump sending message to client %s: %s", c.ID, string(message))
            w, err := c.Conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
