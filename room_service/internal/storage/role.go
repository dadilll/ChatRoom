package storage

import (
	"context"
	"database/sql"
	"errors"
	"room_service/internal/models"
	"time"
)

type RoleStorage interface {
	CreateRole(ctx context.Context, role *models.Role) error
	GetRole(ctx context.Context, roleID string) (*models.Role, error)
	GetRoomRoles(ctx context.Context, roomID string) ([]*models.Role, error)
	UpdateRole(ctx context.Context, roleID string, update *models.UpdateRole) error
	DeleteRole(ctx context.Context, roleID string) error
	AssignRole(ctx context.Context, roomID, userID, roleID string) error
	RemoveRole(ctx context.Context, roomID, userID string) error
	GetUserRole(ctx context.Context, roomID, userID string) (*models.Role, error)
}

type roleStorage struct {
	db *sql.DB
}

func NewRoleStorage(db *sql.DB) RoleStorage {
	return &roleStorage{db: db}
}

func (s *roleStorage) CreateRole(ctx context.Context, role *models.Role) error {
	query := `
        INSERT INTO roles (id, room_id, name, color, priority, permissions, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	_, err := s.db.ExecContext(ctx, query,
		role.ID,
		role.RoomID,
		role.Name,
		role.Color,
		role.Priority,
		role.Permissions,
		role.CreatedAt,
		role.UpdatedAt,
	)
	return err
}

func (s *roleStorage) GetRole(ctx context.Context, roleID string) (*models.Role, error) {
	query := `
        SELECT id, room_id, name, color, priority, permissions, created_at, updated_at
        FROM roles WHERE id = $1
    `

	var role models.Role
	err := s.db.QueryRowContext(ctx, query, roleID).Scan(
		&role.ID,
		&role.RoomID,
		&role.Name,
		&role.Color,
		&role.Priority,
		&role.Permissions,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (s *roleStorage) GetRoomRoles(ctx context.Context, roomID string) ([]*models.Role, error) {
	query := `
        SELECT id, room_id, name, color, priority, permissions, created_at, updated_at
        FROM roles WHERE room_id = $1 ORDER BY priority DESC
    `

	rows, err := s.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.RoomID,
			&role.Name,
			&role.Color,
			&role.Priority,
			&role.Permissions,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *roleStorage) UpdateRole(ctx context.Context, roleID string, update *models.UpdateRole) error {
	query := `
        UPDATE roles 
        SET name = COALESCE($1, name),
            color = COALESCE($2, color),
            priority = COALESCE($3, priority),
            permissions = COALESCE($4, permissions),
            updated_at = $5
        WHERE id = $6
    `

	_, err := s.db.ExecContext(ctx, query,
		update.Name,
		update.Color,
		update.Priority,
		update.Permissions,
		time.Now(),
		roleID,
	)
	return err
}

func (s *roleStorage) DeleteRole(ctx context.Context, roleID string) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, roleID)
	return err
}

func (s *roleStorage) AssignRole(ctx context.Context, roomID, userID, roleID string) error {
	// Сначала удаляем текущую роль пользователя
	err := s.RemoveRole(ctx, roomID, userID)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO room_members (room_id, user_id, role_id, joined_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (room_id, user_id) DO UPDATE SET role_id = $3
    `

	_, err = s.db.ExecContext(ctx, query, roomID, userID, roleID, time.Now())
	return err
}

func (s *roleStorage) RemoveRole(ctx context.Context, roomID, userID string) error {
	query := `UPDATE room_members SET role_id = NULL WHERE room_id = $1 AND user_id = $2`
	_, err := s.db.ExecContext(ctx, query, roomID, userID)
	return err
}

func (s *roleStorage) GetUserRole(ctx context.Context, roomID, userID string) (*models.Role, error) {
	query := `
        SELECT r.id, r.room_id, r.name, r.color, r.priority, r.permissions, r.created_at, r.updated_at
        FROM roles r
        JOIN room_members rm ON r.id = rm.role_id
        WHERE rm.room_id = $1 AND rm.user_id = $2
    `

	var role models.Role
	err := s.db.QueryRowContext(ctx, query, roomID, userID).Scan(
		&role.ID,
		&role.RoomID,
		&role.Name,
		&role.Color,
		&role.Priority,
		&role.Permissions,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
