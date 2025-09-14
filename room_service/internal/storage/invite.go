package storage

import (
	"context"
	"database/sql"
	"room_service/internal/models"
	"room_service/pkg/logger"

	"go.uber.org/zap"
)

type InviteStorage interface {
	CreateInvite(ctx context.Context, invite *models.RoomInvite) error
	GetInvitesByUser(ctx context.Context, userID string) ([]*models.RoomInvite, error)
	GetInviteByID(ctx context.Context, inviteID string) (*models.RoomInvite, error)
	UpdateInviteStatus(ctx context.Context, inviteID string, status string) error
	DeleteInvite(ctx context.Context, inviteID string) error
}

type PostgresInviteStorage struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPostgresInviteStorage(db *sql.DB, log logger.Logger) InviteStorage {
	return &PostgresInviteStorage{
		db:     db,
		logger: log,
	}
}
func (s *PostgresInviteStorage) CreateInvite(ctx context.Context, invite *models.RoomInvite) error {
	query := `
		INSERT INTO room_invites (id, room_id, invited_id, sent_by_id, status, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.ExecContext(ctx, query,
		invite.ID,
		invite.RoomID,
		invite.InvitedID,
		invite.SentByID,
		invite.Status,
		invite.SentAt,
	)
	return err
}

func (s *PostgresInviteStorage) GetInvitesByUser(ctx context.Context, userID string) ([]*models.RoomInvite, error) {
	query := `SELECT id, room_id, invited_id, sent_by_id, status, sent_at FROM room_invites WHERE invited_id = $1`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			s.logger.Error(ctx, "failed to close rows", zap.Error(cerr))
		}
	}()

	var invites []*models.RoomInvite
	for rows.Next() {
		invite := &models.RoomInvite{}
		if err := rows.Scan(
			&invite.ID,
			&invite.RoomID,
			&invite.InvitedID,
			&invite.SentByID,
			&invite.Status,
			&invite.SentAt,
		); err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invites, nil
}

func (s *PostgresInviteStorage) UpdateInviteStatus(ctx context.Context, inviteID string, status string) error {
	query := `UPDATE room_invites SET status = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, status, inviteID)
	return err
}

func (s *PostgresInviteStorage) DeleteInvite(ctx context.Context, inviteID string) error {
	query := `DELETE FROM room_invites WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, inviteID)
	return err
}

func (s *PostgresInviteStorage) GetInviteByID(ctx context.Context, inviteID string) (*models.RoomInvite, error) {
	query := `SELECT id, room_id, invited_id, sent_by_id, status, sent_at FROM room_invites WHERE id = $1`
	invite := &models.RoomInvite{}
	err := s.db.QueryRowContext(ctx, query, inviteID).Scan(
		&invite.ID,
		&invite.RoomID,
		&invite.InvitedID,
		&invite.SentByID,
		&invite.Status,
		&invite.SentAt,
	)
	if err != nil {
		return nil, err
	}
	return invite, nil
}
