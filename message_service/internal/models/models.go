package models

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
	MessageTypeEvent MessageType = "event"
)

type MessageStatus string

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
)

type Message struct {
	ID        string        `json:"id" db:"id"`
	RoomID    string        `json:"room_id" db:"room_id"`
	Content   string        `json:"content" db:"content"`
	Type      MessageType   `json:"type" db:"type"`
	Status    MessageStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
}

type MessageRequest struct {
	RoomID   string      `json:"room_id" validate:"required"`
	SenderID string      `json:"sender_id" validate:"required"`
	Content  string      `json:"content" validate:"required"`
	Type     MessageType `json:"type" validate:"required,oneof=text image file event"`
}

type MessageResponse struct {
	ID        string        `json:"id"`
	RoomID    string        `json:"room_id"`
	Content   string        `json:"content"`
	Type      MessageType   `json:"type"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}
