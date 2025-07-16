package service

import (
	"context"
	"go.uber.org/zap"
	"room_service/internal/models"
	"room_service/internal/storage"
	"room_service/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type RoomService interface {
	NewRoom(request models.Room) (*models.RoomResponse, error)
	GetRoom(id string) (*models.Room, error)
	UpdateRoom(id string, request models.UpdateRoom) error
	DeleteRoom(id string) error
}

type roomServiceImpl struct {
	logger      logger.Logger
	roomStorage storage.RoomStorage
}

func NewRoomService(log logger.Logger, storage storage.RoomStorage) RoomService {
	return &roomServiceImpl{
		logger:      log,
		roomStorage: storage,
	}
}

func (s *roomServiceImpl) NewRoom(request models.Room) (*models.RoomResponse, error) {
	room := models.Room{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Private:     request.Private,
		UserCount:   1,
		Description: request.Description,
		OwnerID:     request.OwnerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.roomStorage.CreateRoom(room)
	if err != nil {
		s.logger.Error(context.Background(), "failed to create room:", zap.Error(err))
		return nil, err
	}
	return &models.RoomResponse{RoomID: room.ID}, nil
}

func (s *roomServiceImpl) GetRoom(id string) (*models.Room, error) {
	room, err := s.roomStorage.GetRoomByID(id)
	if err != nil {
		s.logger.Error(context.Background(), "failed to create room:", zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (s *roomServiceImpl) UpdateRoom(id string, request models.UpdateRoom) error {
	return s.roomStorage.UpdateRoom(id, request)
}

func (s *roomServiceImpl) DeleteRoom(id string) error {
	return s.roomStorage.DeleteRoom(id)
}
