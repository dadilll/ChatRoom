package models

import "time"

type Room struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Private     bool      `json:"private" db:"private"`
	Category    string    `json:"category" db:"category"`
	UserCount   int       `json:"user_count" db:"user_count"`
	Description string    `json:"description" db:"description"`
	OwnerID     string    `json:"-" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateRoom struct {
	Name        string `json:"name" db:"name"`
	Private     bool   `json:"private" db:"private"`
	Description string `json:"description" db:"description"`
	Category    string `json:"category" db:"category"`
}

type RoomMember struct {
	RoomID   string    `json:"room_id" db:"room_id"`
	UserID   string    `json:"user_id" db:"user_id"`
	RoleID   *string   `json:"role_id,omitempty" db:"role_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type Role struct {
	ID          string    `json:"id" db:"id"`
	RoomID      string    `json:"room_id" db:"room_id"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"`
	Priority    int       `json:"priority" db:"priority"`
	Permissions int64     `json:"permissions" db:"permissions"`
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

type RoleWithMembers struct {
	Role
	Members []string `json:"members"`
}

type UpdateRole struct {
	Name        *string `json:"name"`
	Color       *string `json:"color"`
	Priority    *int    `json:"priority"`
	Permissions *int64  `json:"permissions"`
}

// Рекомендуется для статусов объявить константы:
const (
	InviteStatusPending  = "pending"
	InviteStatusAccepted = "accepted"
	InviteStatusDeclined = "declined"
)

type RoomInvite struct {
	ID        string    `json:"id" db:"id"`                 // UUID приглашения
	RoomID    string    `json:"room_id" db:"room_id"`       // ID комнаты
	InvitedID string    `json:"invited_id" db:"invited_id"` // ID приглашённого пользователя
	SentByID  string    `json:"sent_by_id" db:"sent_by_id"` // ID пользователя, отправившего приглашение
	Status    string    `json:"status" db:"status"`         // Статус приглашения ("pending", "accepted", "declined" и т.п.)
	SentAt    time.Time `json:"sent_at" db:"sent_at"`       // Время отправки приглашения
}

type RoomResponse struct {
	RoomID string `json:"id"`
}

type SearchRoomParams struct {
	Name     string `query:"name"`
	Category string `query:"category"`
}

type CreateInviteRequest struct {
	RoomID    string `json:"room_id" validate:"required"`
	InvitedID string `json:"invited_id" validate:"required"`
	SentByID  string `json:"sent_by_id" validate:"required"` // или получай из токена
}
