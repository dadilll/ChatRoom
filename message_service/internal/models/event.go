package models

import "time"

type EventType string

const (
	EventMessageSend    EventType = "message_send"
	EventMessageReceive EventType = "message_receive"
	EventMessageStatus  EventType = "message_status"
	EventUserJoin       EventType = "user_join"
	EventUserLeave      EventType = "user_leave"
)

type WSMessage[T any] struct {
	Event EventType `json:"event"`
	Data  T         `json:"data"`
}

type UserEventData struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

type MessageEventData struct {
	ID        string        `json:"id"`
	RoomID    string        `json:"room_id"`
	Content   string        `json:"content"`
	Type      MessageType   `json:"type"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}
