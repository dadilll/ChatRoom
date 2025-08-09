package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          string    `json:"id" db:"id" validate:"required,uuid4"`
	Name        string    `json:"name" db:"name" validate:"required,min=3,max=100"`
	Private     bool      `json:"private" db:"private"`
	Category    string    `json:"category" db:"category" validate:"required,min=3,max=50"`
	UserCount   int       `json:"user_count" db:"user_count" validate:"min=0"`
	Description string    `json:"description" db:"description" validate:"max=300"`
	OwnerID     string    `json:"-" db:"owner_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RoomMember struct {
	RoomID   string    `json:"room_id" db:"room_id"`
	UserID   string    `json:"user_id" db:"user_id"`
	RoleID   *string   `json:"role_id,omitempty" db:"role_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type Role struct {
	ID          string    `json:"id" db:"id" validate:"required,uuid4"`
	RoomID      string    `json:"room_id" db:"room_id" validate:"required,uuid4"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=50"`
	Color       string    `json:"color" db:"color" validate:"required,len=7"`
	Priority    int       `json:"priority" db:"priority" validate:"min=0"`
	Permissions int64     `json:"permissions" db:"permissions" validate:"min=0"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

const (
	PermissionSendMessage = 1 << iota // 1 (0001) - Отправка сообщений
	PermissionInviteUsers             // 4 (0100) - Приглашение пользователей
	PermissionKickUsers               // 8 (1000) - Исключение пользователей
	PermissionManageRoles             // 16 (10000) - Управление ролями
	PermissionManageRoom              // 32 (100000) - Управление комнатой
)

const (
	InviteStatusPending  = "pending"
	InviteStatusAccepted = "accepted"
	InviteStatusDeclined = "declined"
)

type RoomInvite struct {
	ID        string    `json:"id" db:"id" validate:"required,uuid4"`
	RoomID    string    `json:"room_id" db:"room_id" validate:"required,uuid4"`
	InvitedID string    `json:"invited_id" db:"invited_id" validate:"required,uuid4"`
	SentByID  string    `json:"sent_by_id" db:"sent_by_id" validate:"required,uuid4"`
	Status    string    `json:"status" db:"status" validate:"required,oneof=pending accepted declined"`
	SentAt    time.Time `json:"sent_at" db:"sent_at"`
}

func RoomFromCreateRequest(req CreateRoomRequest, ownerID string) Room {
	now := time.Now()

	return Room{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Private:     *req.Private,
		Category:    req.Category,
		UserCount:   1,
		Description: req.Description,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
