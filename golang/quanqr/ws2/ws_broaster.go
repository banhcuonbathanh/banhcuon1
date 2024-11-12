package ws2

import (
    "encoding/json"
 

)

type Broadcaster struct {
    hub *Hub
}

func NewBroadcaster(h *Hub) *Broadcaster {
    return &Broadcaster{hub: h}
}

func (b *Broadcaster) BroadcastToRole(msg Message, role Role) {
    data, _ := json.Marshal(msg)
    b.hub.mu.Lock()
    for client := range b.hub.Clients {
        if client.Role == role {
            select {
            case client.Send <- data:
            default:
                close(client.Send)
                delete(b.hub.Clients, client)
            }
        }
    }
    b.hub.mu.Unlock()
}

func (b *Broadcaster) BroadcastToRoom(roomID string, msg Message) {
    data, _ := json.Marshal(msg)
    b.hub.mu.Lock()
    if clients, ok := b.hub.RoomMap[roomID]; ok {
        for client := range clients {
            select {
            case client.Send <- data:
            default:
                close(client.Send)
                delete(b.hub.RoomMap[roomID], client)
                delete(b.hub.Clients, client)
            }
        }
    }
    b.hub.mu.Unlock()
}
