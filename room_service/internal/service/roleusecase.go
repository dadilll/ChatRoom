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

type RoleService interface {
	CreateRole(ctx context.Context, request models.Role) (*models.Role, error)
	GetRole(ctx context.Context, id string) (*models.Role, error)
	GetRolesByRoom(ctx context.Context, roomID string) ([]*models.Role, error)
	UpdateRole(ctx context.Context, roleID, roomID string, req models.UpdateRoleRequest) (*models.Role, error)
	DeleteRole(ctx context.Context, id string) error
	AssignRole(ctx context.Context, roomID, userID, roleID string) error
	RemoveRole(ctx context.Context, roomID, userID string) error
	GetUserRole(ctx context.Context, roomID, userID string) (*models.Role, error)
	CheckPermission(ctx context.Context, roomID, userID string, permission int64) (bool, error)
}

type roleServiceImpl struct {
	logger      logger.Logger
	roleStorage storage.RoleStorage
	roomStorage storage.RoomStorage
}

func NewRoleService(log logger.Logger, roleStorage storage.RoleStorage, roomStorage storage.RoomStorage) RoleService {
	return &roleServiceImpl{
		logger:      log,
		roleStorage: roleStorage,
		roomStorage: roomStorage,
	}
}

func (s *roleServiceImpl) CreateRole(ctx context.Context, request models.Role) (*models.Role, error) {
	_, err := s.roomStorage.GetRoomByID(ctx, request.RoomID)
	if err != nil {
		s.logger.Error(ctx, "failed to get room", zap.Error(err))
		return nil, errors.New("room not found")
	}

	role := models.Role{
		ID:          uuid.New().String(),
		RoomID:      request.RoomID,
		Name:        request.Name,
		Color:       request.Color,
		Priority:    request.Priority,
		Permissions: request.Permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.roleStorage.CreateRole(ctx, &role)
	if err != nil {
		s.logger.Error(ctx, "failed to create role", zap.Error(err))
		return nil, err
	}

	return &role, nil
}

func (s *roleServiceImpl) GetRole(ctx context.Context, id string) (*models.Role, error) {
	role, err := s.roleStorage.GetRole(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "failed to get role", zap.Error(err))
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (s *roleServiceImpl) GetRolesByRoom(ctx context.Context, roomID string) ([]*models.Role, error) {
	roles, err := s.roleStorage.GetRoomRoles(ctx, roomID)
	if err != nil {
		s.logger.Error(ctx, "failed to get room roles", zap.Error(err))
		return nil, err
	}
	return roles, nil
}

func (s *roleServiceImpl) UpdateRole(ctx context.Context, roleID, roomID string, req models.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleStorage.GetRole(ctx, roleID)
	if err != nil {
		s.logger.Error(ctx, "failed to get role", zap.Error(err))
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	if role.RoomID != roomID {
		return nil, errors.New("role does not belong to the specified room")
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Color != nil {
		role.Color = *req.Color
	}
	if req.Priority != nil {
		role.Priority = *req.Priority
	}
	if req.Permissions != nil {
		role.Permissions = *req.Permissions
	}
	role.UpdatedAt = time.Now()

	err = s.roleStorage.UpdateRole(ctx, roleID, role)
	if err != nil {
		s.logger.Error(ctx, "failed to update role", zap.Error(err))
		return nil, err
	}

	return s.roleStorage.GetRole(ctx, roleID)
}

func (s *roleServiceImpl) DeleteRole(ctx context.Context, id string) error {
	_, err := s.roleStorage.GetRole(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "failed to get role", zap.Error(err))
		return err
	}

	err = s.roleStorage.DeleteRole(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "failed to delete role", zap.Error(err))
		return err
	}
	return nil
}

func (s *roleServiceImpl) AssignRole(ctx context.Context, roomID, userID, roleID string) error {
	role, err := s.roleStorage.GetRole(ctx, roleID)
	if err != nil {
		s.logger.Error(ctx, "failed to get role", zap.Error(err))
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}
	if role.RoomID != roomID {
		return errors.New("role does not belong to this room")
	}

	err = s.roleStorage.AssignRole(ctx, roomID, userID, roleID)
	if err != nil {
		s.logger.Error(ctx, "failed to assign role", zap.Error(err))
		return err
	}
	return nil
}

func (s *roleServiceImpl) RemoveRole(ctx context.Context, roomID, userID string) error {
	err := s.roleStorage.RemoveRole(ctx, roomID, userID)
	if err != nil {
		s.logger.Error(ctx, "failed to remove role", zap.Error(err))
		return err
	}
	return nil
}

func (s *roleServiceImpl) GetUserRole(ctx context.Context, roomID, userID string) (*models.Role, error) {
	role, err := s.roleStorage.GetUserRole(ctx, roomID, userID)
	if err != nil {
		s.logger.Error(ctx, "failed to get user role", zap.Error(err))
		return nil, err
	}
	return role, nil
}

func (s *roleServiceImpl) CheckPermission(ctx context.Context, roomID, userID string, permission int64) (bool, error) {
	room, err := s.roomStorage.GetRoomByID(ctx, roomID)
	if err != nil {
		s.logger.Error(ctx, "failed to get room", zap.Error(err))
		return false, err
	}
	if room == nil {
		return false, errors.New("room not found")
	}
	if room.OwnerID == userID {
		return true, nil
	}

	role, err := s.GetUserRole(ctx, roomID, userID)
	if err != nil {
		s.logger.Error(ctx, "failed to get user role", zap.Error(err))
		return false, err
	}
	if role == nil {
		return permission == models.PermissionSendMessage, nil
	}

	return (role.Permissions & permission) == permission, nil
}
