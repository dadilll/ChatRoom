// internal/storage/message_storage.go
package storage

import (
	models "message_service/internal/DTO"

	"github.com/jmoiron/sqlx"
)

type MessageStorage interface {
	SaveMessage(msg models.Message) (models.Message, error)
	GetMessages(roomID string, limit, offset int) ([]models.Message, error)
	UpdateStatus(messageID string, status models.MessageStatus) error
}

type PgMessageStorage struct {
	db *sqlx.DB
}

func NewPgMessageStorage(db *sqlx.DB) *PgMessageStorage {
	return &PgMessageStorage{db: db}
}

func (s *PgMessageStorage) SaveMessage(msg models.Message) (models.Message, error) {
	query := `INSERT INTO messages (room_id, content, type, status, created_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(query,
		msg.RoomID, msg.Content, msg.Type, msg.Status, msg.CreatedAt,
	).Scan(&msg.ID)
	return msg, err
}

func (s *PgMessageStorage) GetMessages(roomID string, limit, offset int) ([]models.Message, error) {
	query := `SELECT id, room_id, content, type, status, created_at
	          FROM messages WHERE room_id = $1 ORDER BY created_at DESC
	          LIMIT $2 OFFSET $3`
	rows, err := s.db.Query(query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.RoomID, &m.Content, &m.Type, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *PgMessageStorage) UpdateStatus(messageID string, status models.MessageStatus) error {
	_, err := s.db.Exec(`UPDATE messages SET status=$1 WHERE id=$2`, status, messageID)
	return err
}
