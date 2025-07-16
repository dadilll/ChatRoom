package storage

import (
	"database/sql"
	"service_auth/internal/models"
)

type UserStorage interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	UpdateUserProfile(userID string, profile *models.UserRequestProfile) error
	UpdatePassword(userID string, newHashedPassword string) error
	MarkEmailVerified(userID string) error
	UpdateEmail(userID, newEmail string, isVerified bool) error
}

type userStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) UserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users (email, password, username, bio, avatar_url, description, birthday, is_verified, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, false, NOW()) RETURNING id`

	err := s.db.QueryRow(query,
		user.Email,
		user.Password,
		user.Username,
		user.Bio,
		user.AvatarURL,
		user.Description,
		user.Birthday,
	).Scan(&user.ID)
	return err
}

func (s *userStorage) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, username, is_verified FROM users WHERE email = $1`
	user := &models.User{}
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.IsVerified,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userStorage) GetUserByID(userID string) (*models.User, error) {
	query := `SELECT email, username, bio, avatar_url, description, birthday FROM users WHERE id = $1`
	user := &models.User{}
	err := s.db.QueryRow(query, userID).Scan(
		&user.Email,
		&user.Username,
		&user.Bio,
		&user.AvatarURL,
		&user.Description,
		&user.Birthday,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userStorage) UpdateUserProfile(userID string, profile *models.UserRequestProfile) error {
	query := `
		UPDATE users
		SET username = $1,
			bio = $2,
			avatar_url = $3,
			description = $4,
			birthday = $5
		WHERE id = $6`
	_, err := s.db.Exec(query,
		profile.Username,
		profile.Bio,
		profile.AvatarURL,
		profile.Description,
		profile.Birthday,
		userID,
	)
	return err
}

func (s *userStorage) UpdatePassword(userID string, newHashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := s.db.Exec(query, newHashedPassword, userID)
	return err
}

func (s *userStorage) MarkEmailVerified(userID string) error {
	query := `UPDATE users SET is_verified = true WHERE id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}

func (s *userStorage) UpdateEmail(userID, newEmail string, isVerified bool) error {
	query := `UPDATE users SET email = $1, is_verified = $2 WHERE id = $3`
	_, err := s.db.Exec(query, newEmail, isVerified, userID)
	return err
}
