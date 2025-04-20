package handler

import (
	"encoding/json"
	"errors"
	"notification_service/internal/models"
	"notification_service/internal/service"
)

type Handler interface {
	Handle(data []byte) error
}

type notificationHandler struct {
	mailer service.Mailer
}

func NewNotificationHandler(m service.Mailer) Handler {
	return &notificationHandler{
		mailer: m,
	}
}

func (h *notificationHandler) Handle(data []byte) error {
	var msg models.EmailMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return errors.New("не удалось распарсить сообщение")
	}
	return h.mailer.Send(msg)
}
