package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"room_service/internal/models"
	"room_service/internal/storage"
)

type InviteService interface {
	NewInvite(ctx context.Context, roomID, invitedID, sentByID string) (*models.RoomInvite, error)
	GetUserInvites(ctx context.Context, userID string) ([]*models.RoomInvite, error)
	AcceptInvite(ctx context.Context, inviteID string) error
	DeclineInvite(ctx context.Context, inviteID string) error
	DeleteInvite(ctx context.Context, inviteID string) error
}

type inviteServiceImpl struct {
	store storage.InviteStorage
}

func NewInviteService(store storage.InviteStorage) InviteService {
	return &inviteServiceImpl{store: store}
}

func (s *inviteServiceImpl) NewInvite(ctx context.Context, roomID, invitedID, sentByID string) (*models.RoomInvite, error) {
	invite := &models.RoomInvite{
		ID:        uuid.NewString(),
		RoomID:    roomID,
		InvitedID: invitedID,
		SentByID:  sentByID,
		Status:    "pending",
		SentAt:    time.Now(),
	}

	if err := s.store.CreateInvite(ctx, invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *inviteServiceImpl) GetUserInvites(ctx context.Context, userID string) ([]*models.RoomInvite, error) {
	return s.store.GetInvitesByUser(ctx, userID)
}

func (s *inviteServiceImpl) AcceptInvite(ctx context.Context, inviteID string) error {
	return s.store.UpdateInviteStatus(ctx, inviteID, "accepted")
}

func (s *inviteServiceImpl) DeclineInvite(ctx context.Context, inviteID string) error {
	return s.store.UpdateInviteStatus(ctx, inviteID, "declined")
}

func (s *inviteServiceImpl) DeleteInvite(ctx context.Context, inviteID string) error {
	return s.store.DeleteInvite(ctx, inviteID)
}
