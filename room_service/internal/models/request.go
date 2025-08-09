package models

type CreateInviteRequest struct {
	RoomID    string `json:"room_id" validate:"required"`
	InvitedID string `json:"invited_id" validate:"required"`
}

type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Private     *bool  `json:"private" validate:"required"`
	Category    string `json:"category" validate:"required,min=3,max=50"`
	Description string `json:"description,omitempty" validate:"max=300"`
}

type UpdateRoomRequest struct {
	Name        string `json:"name" db:"name" validate:"required,min=3,max=100"`
	Private     bool   `json:"private" db:"private"`
	Description string `json:"description" db:"description" validate:"max=300"`
	Category    string `json:"category" db:"category" validate:"required,min=3,max=50"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=50"`
	Color       *string `json:"color" validate:"omitempty,len=7"`
	Priority    *int    `json:"priority" validate:"omitempty,min=0"`
	Permissions *int64  `json:"permissions" validate:"omitempty,min=0"`
}

type AddMemberRequest struct {
	RoomID string `param:"room_id" validate:"required,uuid4"`
}

type ListMembersRequest struct {
	RoomID string `param:"room_id" validate:"required,uuid4"`
}

type RemoveMemberRequest struct {
	RoomID string `param:"room_id" validate:"required,uuid4"`
	UserID string `param:"user_id" validate:"required,uuid4"`
}

type UserIDParam struct {
	UserID string `param:"user_id" validate:"required,uuid4"`
}

type InviteIDParam struct {
	InviteID string `param:"invite_id" validate:"required,uuid4"`
}
