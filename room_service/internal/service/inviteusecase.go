package service

import (
	"context"
	"fmt"
	"time"

	"room_service/internal/models"
	"room_service/internal/storage"

	"github.com/google/uuid"
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
		Status:    models.InviteStatusPending,
		SentAt:    time.Now(),
	}

	if err := s.store.CreateInvite(ctx, invite); err != nil {
		return nil, fmt.Errorf("create invite: %w", err)
	}
	return invite, nil
}

func (s *inviteServiceImpl) GetUserInvites(ctx context.Context, userID string) ([]*models.RoomInvite, error) {
	invites, err := s.store.GetInvitesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get invites by user: %w", err)
	}
	return invites, nil
}

func (s *inviteServiceImpl) AcceptInvite(ctx context.Context, inviteID string) error {
	return s.updateInviteStatus(ctx, inviteID, models.InviteStatusAccepted)
}

func (s *inviteServiceImpl) DeclineInvite(ctx context.Context, inviteID string) error {
	return s.updateInviteStatus(ctx, inviteID, models.InviteStatusDeclined)
}

func (s *inviteServiceImpl) DeleteInvite(ctx context.Context, inviteID string) error {
	if err := s.store.DeleteInvite(ctx, inviteID); err != nil {
		return fmt.Errorf("delete invite: %w", err)
	}
	return nil
}

func (s *inviteServiceImpl) updateInviteStatus(ctx context.Context, inviteID, newStatus string) error {
	invite, err := s.store.GetInviteByID(ctx, inviteID)
	if err != nil {
		return fmt.Errorf("get invite by id: %w", err)
	}

	if invite.Status != models.InviteStatusPending {
		return fmt.Errorf("invite already processed: current status = %s", invite.Status)
	}

	if err := s.store.UpdateInviteStatus(ctx, inviteID, newStatus); err != nil {
		return fmt.Errorf("update invite status: %w", err)
	}
	return nil
}
