package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	Username     string    `json:"username" db:"username"`
	FriendsCount int       `json:"friends_count" db:"friends_count"`
	Bio          string    `json:"bio" db:"bio"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
	Description  string    `json:"description" db:"description"`
	Birthday     time.Time `json:"birthday" db:"birthday"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Message struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
