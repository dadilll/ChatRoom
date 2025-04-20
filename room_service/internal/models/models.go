package models

import "time"

type Room struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Private     bool      `json:"private" db:"private"`
	Usercount   int       `json:"user_count" db:"user_count"`
	Description string    `json:"description" db:"description"`
	OwnerID     string    `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RoomMember struct {
	RoomID   string    `json:"room_id" db:"room_id"`
	UserID   string    `json:"user_id" db:"user_id"`
	Role     string    `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type Role struct {
	Name      string `json:"name" db:"name"`
	CanSend   bool   `json:"can_send" db:"can_send"`
	CanDelete bool   `json:"can_delete" db:"can_delete"`
	CanInvite bool   `json:"can_invite" db:"can_invite"`
	CanKick   bool   `json:"can_kick" db:"can_kick"`
}

type RoomInvite struct {
	ID        string    `json:"id" db:"id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	InvitedID string    `json:"invited_id" db:"invited_id"`
	SentByID  string    `json:"sent_by_id" db:"sent_by_id"`
	Status    string    `json:"status" db:"status"`
	SentAt    time.Time `json:"sent_at" db:"sent_at"`
}

type WebSocketEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
