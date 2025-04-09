package models

import "time"

type UserRequestRegister struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Username    string    `json:"username"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatarURL"`
	Description string    `json:"description"`
	Birthday    time.Time `json:"birthday"`
}

type UserRequestProfile struct {
	Username    string    `json:"username"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatarURL"`
	Description string    `json:"description"`
	Birthday    time.Time `json:"birthday"`
}

type UserRequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequestPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserRequestConfirmationCode struct {
	Code string `json:"code"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}
