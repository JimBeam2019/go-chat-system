package websocket

import (
	"go-chat-system/internal/infrastructure/hub"
	"go-chat-system/internal/usecase"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub     *hub.Hub
	usecase *usecase.ChatUsecase
}

func NewHandler(h *hub.Hub, u *usecase.ChatUsecase) *Handler {
	return &Handler{hub: h, usecase: u}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	userID := r.URL.Query().Get("user")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &hub.Client{
		ID:     userID,
		RoomID: roomID,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register(client)

	go h.writePump(conn, client)
	go h.readPump(conn, client)
}

func (h *Handler) readPump(conn *websocket.Conn, client *hub.Client) {
	defer func() {
		h.hub.Unregister(client)
		conn.Close()
	}()

	for {
		var msg usecase.Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		h.usecase.HandleMessage(msg)
	}
}

func (h *Handler) writePump(conn *websocket.Conn, client *hub.Client) {
	defer conn.Close()

	for msg := range client.Send {
		if err := conn.WriteMessage(
			websocket.TextMessage, msg,
		); err != nil {
			break
		}
	}
}
