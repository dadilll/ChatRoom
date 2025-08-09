package service

import (
	"context"
	"errors"
	"room_service/internal/models"
	"room_service/internal/storage"
	"room_service/pkg/logger"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
)

type RoomService interface {
	NewRoom(ctx context.Context, request models.Room) (*models.RoomResponse, error)
	GetRoom(ctx context.Context, id string) (*models.Room, error)
	UpdateRoom(ctx context.Context, id string, request models.UpdateRoomRequest) error
	DeleteRoom(ctx context.Context, id string) error
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

func (s *roomServiceImpl) NewRoom(ctx context.Context, request models.Room) (*models.RoomResponse, error) {
	if request.Name == "" {
		return nil, errors.New("room name is required")
	}

	now := time.Now()
	room := models.Room{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Private:     request.Private,
		Category:    request.Category,
		UserCount:   1,
		Description: request.Description,
		OwnerID:     request.OwnerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.roomStorage.CreateRoom(ctx, room); err != nil {
		s.logger.Error(ctx, "failed to create room", zap.Error(err))
		return nil, err
	}

	return &models.RoomResponse{RoomID: room.ID}, nil
}

func (s *roomServiceImpl) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	room, err := s.roomStorage.GetRoomByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "failed to get room", zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (s *roomServiceImpl) UpdateRoom(ctx context.Context, id string, request models.UpdateRoomRequest) error {
	err := s.roomStorage.UpdateRoom(ctx, id, request)
	if err != nil {
		s.logger.Error(ctx, "failed to update room", zap.String("room_id", id), zap.Error(err))
	}
	return err
}

func (s *roomServiceImpl) DeleteRoom(ctx context.Context, id string) error {
	err := s.roomStorage.DeleteRoom(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "failed to delete room", zap.String("room_id", id), zap.Error(err))
	}
	return err
}
