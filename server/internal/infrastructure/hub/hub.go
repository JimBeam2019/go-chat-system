package hub

import "sync"

type Client struct {
	ID     string
	RoomID string
	Send   chan []byte
}

type broadcastMsg struct {
	roomID string
	data   []byte
}

type Hub struct {
	rooms map[string]map[*Client]bool
	mtx   sync.RWMutex

	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMsg
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client, 1024),
		unregister: make(chan *Client, 1024),
		broadcast:  make(chan broadcastMsg, 1024),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mtx.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mtx.Unlock()

		case client := <-h.unregister:
			h.mtx.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
				}
			}
			h.mtx.Unlock()

		case msg := <-h.broadcast:
			h.mtx.RLock()
			clients := h.rooms[msg.roomID]
			for client := range clients {
				select {
				case client.Send <- msg.data:
				default:
					// Drop slow client (critical for high throughput)
					close(client.Send)
					delete(clients, client)
				}
			}
			h.mtx.RUnlock()
		}
	}
}

func (h *Hub) Register(c *Client) {
	h.register <- c
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

func (h *Hub) Broadcast(roomID string, data []byte) {
	h.broadcast <- broadcastMsg{
		roomID: roomID,
		data:   data,
	}
}
