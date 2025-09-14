package storage

import (
	"context"
	"database/sql"
	"room_service/internal/models"
	"room_service/pkg/logger"

	"go.uber.org/zap"
)

type RoomMemberStorage interface {
	AddMember(ctx context.Context, member models.RoomMember) error
	RemoveMember(ctx context.Context, roomID, userID string) error
	ListMembers(ctx context.Context, roomID string) ([]models.RoomMember, error)
	AssignRole(ctx context.Context, roomID, userID string, roleID *string) error
	IsMember(ctx context.Context, roomID, userID string) (bool, error)
}

type RoomMemberRepo struct {
	db     *sql.DB
	logger logger.Logger
}

func NewRoomMemberStorage(db *sql.DB, log logger.Logger) *RoomMemberRepo {
	return &RoomMemberRepo{
		db:     db,
		logger: log,
	}
}

func (r *RoomMemberRepo) IsMember(ctx context.Context, roomID, userID string) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2
	)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, roomID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *RoomMemberRepo) AddMember(ctx context.Context, member models.RoomMember) error {
	query := `
		INSERT INTO room_members (room_id, user_id, role_id, joined_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (room_id, user_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		member.RoomID,
		member.UserID,
		member.RoleID,
		member.JoinedAt,
	)
	return err
}

func (r *RoomMemberRepo) RemoveMember(ctx context.Context, roomID, userID string) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, roomID, userID)
	return err
}

func (r *RoomMemberRepo) AssignRole(ctx context.Context, roomID, userID string, roleID *string) error {
	query := `
		UPDATE room_members
		SET role_id = $1
		WHERE room_id = $2 AND user_id = $3`

	_, err := r.db.ExecContext(ctx, query, roleID, roomID, userID)
	return err
}

func (r *RoomMemberRepo) ListMembers(ctx context.Context, roomID string) ([]models.RoomMember, error) {
	query := `
		SELECT room_id, user_id, role_id, joined_at
		FROM room_members
		WHERE room_id = $1
		ORDER BY joined_at ASC`

	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error(ctx, "failed to close rows", zap.Error(err))
		}
	}()

	var members []models.RoomMember
	for rows.Next() {
		var member models.RoomMember
		if err := rows.Scan(
			&member.RoomID,
			&member.UserID,
			&member.RoleID,
			&member.JoinedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}
