package storage

import (
	"context"
	"database/sql"
	"room_service/internal/models"
)

type RoomStorage interface {
	CreateRoom(room models.Room) error
	GetRoomByID(id string) (*models.Room, error)
	UpdateRoom(id string, updated models.UpdateRoom) error
	DeleteRoom(id string) error
	IncrementUserCount(ctx context.Context, roomID string) error
	DecrementUserCount(ctx context.Context, roomID string) error
}

type roomStorage struct {
	db *sql.DB
}

func NewRoomStorage(db *sql.DB) RoomStorage {
	return &roomStorage{db: db}
}

func (r *roomStorage) CreateRoom(room models.Room) error {
	_, err := r.db.Exec(`
		INSERT INTO rooms (id, name, private, user_count, description, owner_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, room.ID, room.Name, room.Private, room.UserCount, room.Description, room.OwnerID, room.CreatedAt, room.UpdatedAt)
	return err
}

func (r *roomStorage) GetRoomByID(id string) (*models.Room, error) {
	row := r.db.QueryRow(`SELECT * FROM rooms WHERE id=$1`, id)
	var room models.Room
	err := row.Scan(&room.ID, &room.Name, &room.Private, &room.UserCount, &room.Description, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomStorage) UpdateRoom(id string, updated models.UpdateRoom) error {
	_, err := r.db.Exec(`
		UPDATE rooms SET name=$1, private=$2, description=$3, updated_at=NOW() WHERE id=$4
	`, updated.Name, updated.Private, updated.Description, id)
	return err
}

func (r *roomStorage) DeleteRoom(id string) error {
	_, err := r.db.Exec(`DELETE FROM rooms WHERE id=$1`, id)
	return err
}

func (r *roomStorage) IncrementUserCount(ctx context.Context, roomID string) error {
	query := `UPDATE rooms SET user_count = user_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, roomID)
	return err
}

func (r *roomStorage) DecrementUserCount(ctx context.Context, roomID string) error {
	query := `UPDATE rooms SET user_count = GREATEST(user_count - 1, 0) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, roomID)
	return err
}
