package main

import (
	"go-chat-system/internal/delivery/websocket"
	"go-chat-system/internal/infrastructure/hub"
	"go-chat-system/internal/usecase"
	"log"
	"net/http"
)

func main() {
	h := hub.NewHub()

	go h.Run()

	uc := usecase.NewChatUsecase(h)
	handler := websocket.NewHandler(h, uc)

	http.HandleFunc("/ws", handler.ServeWS)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
