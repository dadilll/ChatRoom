package service

import (
	models "message_service/internal/DTO"
	"message_service/internal/storage"
	"time"
)

type MessageService interface {
	SendMessage(req models.MessageRequest) (models.MessageResponse, error)
	GetRoomMessages(roomID string, limit, offset int) ([]models.MessageResponse, error)
	UpdateMessageStatus(messageID string, status models.MessageStatus) error
}

type messageService struct {
	storage storage.MessageStorage
}

func NewMessageService(s storage.MessageStorage) MessageService {
	return &messageService{storage: s}
}

func (s *messageService) SendMessage(req models.MessageRequest) (models.MessageResponse, error) {
	msg := models.Message{
		RoomID:    req.RoomID,
		Content:   req.Content,
		Type:      req.Type,
		Status:    models.StatusSent,
		CreatedAt: time.Now(),
	}
	saved, err := s.storage.SaveMessage(msg)
	if err != nil {
		return models.MessageResponse{}, err
	}

	return models.MessageResponse{
		ID:        saved.ID,
		RoomID:    saved.RoomID,
		Content:   saved.Content,
		Type:      saved.Type,
		Status:    saved.Status,
		CreatedAt: saved.CreatedAt,
	}, nil
}

func (s *messageService) GetRoomMessages(roomID string, limit, offset int) ([]models.MessageResponse, error) {
	msgs, err := s.storage.GetMessages(roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	res := make([]models.MessageResponse, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, models.MessageResponse{
			ID:        m.ID,
			RoomID:    m.RoomID,
			Content:   m.Content,
			Type:      m.Type,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
		})
	}
	return res, nil
}

func (s *messageService) UpdateMessageStatus(messageID string, status models.MessageStatus) error {
	return s.storage.UpdateStatus(messageID, status)
}
