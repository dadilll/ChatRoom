package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"room_service/internal/models"
	"room_service/internal/storage"
	"room_service/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type RoleService interface {
	CreateRole(request models.Role) (*models.Role, error)
	GetRole(id string) (*models.Role, error)
	GetRolesByRoom(roomID string) ([]*models.Role, error)
	UpdateRole(id string, request models.UpdateRole) (*models.Role, error)
	DeleteRole(id string) error
	AssignRole(roomID, userID, roleID string) error
	RemoveRole(roomID, userID string) error
	GetUserRole(roomID, userID string) (*models.Role, error)
	CheckPermission(roomID, userID string, permission int64) (bool, error)
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

func (s *roleServiceImpl) CreateRole(request models.Role) (*models.Role, error) {
	_, err := s.roomStorage.GetRoomByID(request.RoomID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get room:", zap.Error(err))
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

	err = s.roleStorage.CreateRole(context.Background(), &role)
	if err != nil {
		s.logger.Error(context.Background(), "failed to create role:", zap.Error(err))
		return nil, err
	}

	return &role, nil
}

func (s *roleServiceImpl) GetRole(id string) (*models.Role, error) {
	role, err := s.roleStorage.GetRole(context.Background(), id)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get role:", zap.Error(err))
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	return role, nil
}

func (s *roleServiceImpl) GetRolesByRoom(roomID string) ([]*models.Role, error) {
	roles, err := s.roleStorage.GetRoomRoles(context.Background(), roomID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get room roles:", zap.Error(err))
		return nil, err
	}

	return roles, nil
}

func (s *roleServiceImpl) UpdateRole(id string, request models.UpdateRole) (*models.Role, error) {
	role, err := s.roleStorage.GetRole(context.Background(), id)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get role:", zap.Error(err))
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	if request.Name != nil {
		role.Name = *request.Name
	}
	if request.Color != nil {
		role.Color = *request.Color
	}
	if request.Priority != nil {
		role.Priority = *request.Priority
	}
	if request.Permissions != nil {
		role.Permissions = *request.Permissions
	}
	role.UpdatedAt = time.Now()

	err = s.roleStorage.UpdateRole(context.Background(), id, &request)
	if err != nil {
		s.logger.Error(context.Background(), "failed to update role:", zap.Error(err))
		return nil, err
	}

	return s.roleStorage.GetRole(context.Background(), id)
}

func (s *roleServiceImpl) DeleteRole(id string) error {
	_, err := s.roleStorage.GetRole(context.Background(), id)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get role:", zap.Error(err))
		return err
	}

	err = s.roleStorage.DeleteRole(context.Background(), id)
	if err != nil {
		s.logger.Error(context.Background(), "failed to delete role:", zap.Error(err))
		return err
	}

	return nil
}

func (s *roleServiceImpl) AssignRole(roomID, userID, roleID string) error {
	role, err := s.roleStorage.GetRole(context.Background(), roleID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get role:", zap.Error(err))
		return err
	}

	if role == nil {
		return errors.New("role not found")
	}

	if role.RoomID != roomID {
		return errors.New("role does not belong to this room")
	}

	err = s.roleStorage.AssignRole(context.Background(), roomID, userID, roleID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to assign role:", zap.Error(err))
		return err
	}

	return nil
}

func (s *roleServiceImpl) RemoveRole(roomID, userID string) error {
	err := s.roleStorage.RemoveRole(context.Background(), roomID, userID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to remove role:", zap.Error(err))
		return err
	}

	return nil
}

func (s *roleServiceImpl) GetUserRole(roomID, userID string) (*models.Role, error) {
	role, err := s.roleStorage.GetUserRole(context.Background(), roomID, userID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get user role:", zap.Error(err))
		return nil, err
	}

	return role, nil
}

func (s *roleServiceImpl) CheckPermission(roomID, userID string, permission int64) (bool, error) {
	room, err := s.roomStorage.GetRoomByID(roomID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get room:", zap.Error(err))
		return false, err
	}

	if room == nil {
		return false, errors.New("room not found")
	}

	if room.OwnerID == userID {
		return true, nil
	}

	role, err := s.GetUserRole(roomID, userID)
	if err != nil {
		s.logger.Error(context.Background(), "failed to get user role:", zap.Error(err))
		return false, err
	}

	if role == nil {
		// По умолчанию разрешено только отправлять сообщения
		return permission == models.PermissionSendMessage, nil
	}

	return (role.Permissions & permission) == permission, nil
}
