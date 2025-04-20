package models

type EmailMessage struct {
	To      string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
