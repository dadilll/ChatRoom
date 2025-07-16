package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"room_service/internal/models"
	"room_service/internal/storage"
	"room_service/pkg/logger"
)

type MemberService interface {
	AddMember(ctx context.Context, member models.RoomMember) error
	ListMembers(ctx context.Context, roomID string) ([]models.RoomMember, error)
	RemoveMember(ctx context.Context, roomID, targetUserID, requesterID string) error
	IsRoomPublic(roomID string) (bool, error)
	HasInvite(ctx context.Context, roomID, userID string) (bool, error)
}

type memberServicelmpl struct {
	logger        logger.Logger
	storage       storage.RoomMemberStorage
	roleSvc       RoleService
	roomStorage   storage.RoomStorage
	inviteService InviteService
}

func NewMemberService(
	s storage.RoomMemberStorage,
	r RoleService,
	l logger.Logger,
	room storage.RoomStorage,
	inviteSvc InviteService,
) MemberService {
	return &memberServicelmpl{
		storage:       s,
		roleSvc:       r,
		logger:        l,
		roomStorage:   room,
		inviteService: inviteSvc,
	}
}

func (s *memberServicelmpl) AddMember(ctx context.Context, member models.RoomMember) error {
	exists, err := s.storage.IsMember(ctx, member.RoomID, member.UserID)
	if err != nil {
		s.logger.Error(ctx, "failed to check member existence", zap.Error(err))
		return err
	}
	if exists {
		s.logger.Error(ctx, "user already in room", zap.String("room_id", member.RoomID), zap.String("user_id", member.UserID))
		return errors.New("user already in room")
	}

	isPublic, err := s.IsRoomPublic(member.RoomID)
	if err != nil {
		s.logger.Error(ctx, "failed to check if room is public", zap.Error(err))
		return err
	}

	if !isPublic {
		hasInvite, err := s.HasInvite(ctx, member.RoomID, member.UserID)
		if err != nil {
			s.logger.Error(ctx, "failed to check invite", zap.Error(err))
			return err
		}
		if !hasInvite {
			s.logger.Error(ctx, "invite required to join room", zap.String("room_id", member.RoomID), zap.String("user_id", member.UserID))
			return errors.New("invite required")
		}
	}

	err = s.storage.AddMember(ctx, member)
	if err != nil {
		s.logger.Error(ctx, "failed to add member", zap.Error(err))
		return err
	}

	if err := s.roomStorage.IncrementUserCount(ctx, member.RoomID); err != nil {
		s.logger.Error(ctx, "failed to increment user count", zap.Error(err))
	}

	s.logger.Info(ctx, "member added", zap.String("room_id", member.RoomID), zap.String("user_id", member.UserID))
	return nil
}
func (s *memberServicelmpl) ListMembers(ctx context.Context, roomID string) ([]models.RoomMember, error) {
	userID, _ := ctx.Value("userID").(string)
	isMember, err := s.storage.IsMember(ctx, roomID, userID)
	if err != nil {
		s.logger.Error(ctx, "failed to check membership", zap.Error(err))
		return nil, err
	}
	if !isMember {
		s.logger.Error(ctx, "user is not a member of the room", zap.String("room_id", roomID), zap.String("user_id", userID))
		return nil, errors.New("forbidden: not a member of the room")
	}

	members, err := s.storage.ListMembers(ctx, roomID)
	if err != nil {
		s.logger.Error(ctx, "failed to list members", zap.String("room_id", roomID), zap.Error(err))
		return nil, err
	}
	return members, nil
}

func (s *memberServicelmpl) RemoveMember(ctx context.Context, roomID, targetUserID, requesterID string) error {
	if requesterID != targetUserID {
		s.logger.Error(ctx, "user tried to remove another member",
			zap.String("requester_id", requesterID),
			zap.String("target_user_id", targetUserID))
		return errors.New("cannot remove other members")
	}

	err := s.storage.RemoveMember(ctx, roomID, requesterID)
	if err != nil {
		s.logger.Error(ctx, "failed to leave room", zap.Error(err))
		return err
	}

	if err := s.roomStorage.DecrementUserCount(ctx, roomID); err != nil {
		s.logger.Error(ctx, "failed to decrement user count", zap.Error(err))
	}

	s.logger.Info(ctx, "user left room", zap.String("room_id", roomID), zap.String("user_id", requesterID))
	return nil
}

func (s *memberServicelmpl) IsRoomPublic(roomID string) (bool, error) {
	room, err := s.roomStorage.GetRoomByID(roomID)
	if err != nil {
		return false, err
	}
	return !room.Private, nil
}

func (s *memberServicelmpl) HasInvite(ctx context.Context, roomID, userID string) (bool, error) {
	invites, err := s.inviteService.GetUserInvites(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, invite := range invites {
		if invite.RoomID == roomID && invite.Status == "pending" {
			return true, nil
		}
	}
	return false, nil
}
