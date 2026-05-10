package usecase

import "encoding/json"

type Broadcaster interface {
	Broadcast(roomID string, data []byte)
}

type ChatUsecase struct {
	hub Broadcaster
}

func NewChatUsecase(h Broadcaster) *ChatUsecase {
	return &ChatUsecase{hub: h}
}

type Message struct {
	RoomID  string `json:"room_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

func (u *ChatUsecase) HandleMessage(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	u.hub.Broadcast(msg.RoomID, data)
	return nil
}
