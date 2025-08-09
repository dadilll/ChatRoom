package models

type RoomResponse struct {
	RoomID string `json:"id"`
}

type RoleWithMembers struct {
	Role
	Members []string `json:"members"`
}
