package domain

type Client struct {
	ID     string
	RoomID string
	Send   chan []byte
}
