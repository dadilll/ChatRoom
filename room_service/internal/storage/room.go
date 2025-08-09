package storage

import (
	"context"
	"database/sql"
	"room_service/internal/models"
)

type RoomStorage interface {
	CreateRoom(ctx context.Context, room models.Room) error
	GetRoomByID(ctx context.Context, id string) (*models.Room, error)
	UpdateRoom(ctx context.Context, id string, updated models.UpdateRoomRequest) error
	DeleteRoom(ctx context.Context, id string) error
	IncrementUserCount(ctx context.Context, roomID string) error
	DecrementUserCount(ctx context.Context, roomID string) error
}

type roomStorage struct {
	db *sql.DB
}

func NewRoomStorage(db *sql.DB) RoomStorage {
	return &roomStorage{db: db}
}

func (r *roomStorage) CreateRoom(ctx context.Context, room models.Room) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO rooms (id, name, private, category, user_count, description, owner_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, room.ID, room.Name, room.Private, room.Category, room.UserCount, room.Description, room.OwnerID, room.CreatedAt, room.UpdatedAt)
	return err
}

func (r *roomStorage) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, private, category, user_count, description, owner_id, created_at, updated_at FROM rooms WHERE id=$1`, id)

	var room models.Room
	err := row.Scan(&room.ID, &room.Name, &room.Private, &room.Category, &room.UserCount, &room.Description, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomStorage) UpdateRoom(ctx context.Context, id string, updated models.UpdateRoomRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE rooms SET name=$1, private=$2, description=$3, category=$4, updated_at=NOW()
		WHERE id=$5
	`, updated.Name, updated.Private, updated.Description, updated.Category, id)
	return err
}

func (r *roomStorage) DeleteRoom(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id=$1`, id)
	return err
}

func (r *roomStorage) IncrementUserCount(ctx context.Context, roomID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE rooms SET user_count = user_count + 1 WHERE id = $1`, roomID)
	return err
}

func (r *roomStorage) DecrementUserCount(ctx context.Context, roomID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE rooms SET user_count = GREATEST(user_count - 1, 0) WHERE id = $1`, roomID)
	return err
}
