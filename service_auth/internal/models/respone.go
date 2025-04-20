package models

type UserResponse struct {
	ID       string `json:"id"`
	JWTtoken string `json:"jwtoken"`
}
