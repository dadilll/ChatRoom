package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"service_auth/internal/models"
	"service_auth/internal/storage"
	"service_auth/pkg/jwt"
	"service_auth/pkg/logger"
	"service_auth/pkg/mailer"
	"strings"
	"time"
	"unicode"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"
)

type UserService interface {
	Register(request models.UserRequestRegister) (*models.UserResponse, error)
	Login(request models.UserRequestLogin) (*models.UserResponse, error)
	GetProfile(userID string) (*models.UserRequestProfile, error)
	ProfileUpdate(userID string, profile models.UserRequestProfile) error
	ChangePassword(userID string, oldPassword, newPassword string) error
	SendConfirmationCode(userID string) error
	VerifyConfirmationCode(userID string, inputCode string) error
	ChangeEmail(userID, newEmail string) error
}

type userService struct {
	userStorage storage.UserStorage
	logger      logger.Logger
	redis       storage.UserStorageRedis
	privateKey  *rsa.PrivateKey
	mailer      mailer.Mailer
}

func NewUserService(userStorage storage.UserStorage, redis storage.UserStorageRedis, log logger.Logger, privateKey *rsa.PrivateKey, mailer mailer.Mailer) UserService {
	return &userService{
		userStorage: userStorage,
		redis:       redis,
		logger:      log,
		privateKey:  privateKey,
		mailer:      mailer,
	}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasDigit := false
	hasSpecial := false
	specialChars := "@$!%*?&"

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case strings.ContainsRune(specialChars, c):
			hasSpecial = true
		}
	}

	return hasUpper && hasDigit && hasSpecial
}
func generateRandomCode() (string, error) {
	max := big.NewInt(1000000) // 6-значный код: от 000000 до 999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
func (s *userService) Register(request models.UserRequestRegister) (*models.UserResponse, error) {
	if !isValidEmail(request.Email) {
		return nil, errors.New("некорректный email")
	}

	if !isValidPassword(request.Password) {
		return nil, errors.New("пароль слишком слабый: должен содержать минимум 8 символов, цифру и спецсимвол")
	}

	existingUser, _ := s.userStorage.GetUserByEmail(request.Email)
	if existingUser != nil {
		s.logger.Error(context.Background(), "Пользователь с таким email уже существует:", zap.String("email", request.Email))
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка при хешировании пароля:", zap.Error(err))
		return nil, errors.New("ошибка при хешировании пароля")
	}

	// Создание пользователя
	user := &models.User{
		Email:       request.Email,
		Password:    string(hashedPassword),
		Username:    request.Username,
		Bio:         request.Bio,
		AvatarURL:   request.AvatarURL,
		Description: request.Description,
		Birthday:    request.Birthday,
		IsVerified:  false,
		IsActive:    true, //заглушка
		CreatedAt:   time.Now(),
	}

	err = s.userStorage.CreateUser(user)
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка при создании пользователя:", zap.Error(err))
		return nil, errors.New("ошибка при создании пользователя")
	}

	token, err := jwt.GenerateJWT(user.ID, s.privateKey)
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка при создании JWT-токена:", zap.Error(err))
		return nil, errors.New("ошибка при создании JWT-токена")
	}

	s.logger.Info(context.Background(), "Попытка регистрации пользователя с email:", zap.String("email", request.Email))

	return &models.UserResponse{
		ID:       string(user.ID),
		JWTtoken: token,
	}, nil
}
func (s *userService) Login(request models.UserRequestLogin) (*models.UserResponse, error) {
	user, err := s.userStorage.GetUserByEmail(request.Email)
	if err != nil {
		s.logger.Error(context.Background(), "Неудачная попытка входа: неверный email", zap.String("email", request.Email))
		return nil, errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		s.logger.Error(context.Background(), "Неудачная попытка входа: неверный пароль для email", zap.String("email", request.Email))
		return nil, errors.New("неверный email или пароль")
	}

	token, err := jwt.GenerateJWT(user.ID, s.privateKey)
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка при создании JWT-токена:", zap.Error(err))
		return nil, errors.New("ошибка при создании JWT-токена")
	}

	s.logger.Info(context.Background(), "Пользователь успешно вошел в систему:", zap.String("email", user.Email))

	return &models.UserResponse{
		ID:       string(user.ID),
		JWTtoken: token,
	}, nil
}
func (s *userService) GetProfile(userID string) (*models.UserRequestProfile, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &models.UserRequestProfile{
		Username:    user.Username,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		Description: user.Description,
		Birthday:    user.Birthday,
	}, nil
}
func (s *userService) ProfileUpdate(userID string, profile models.UserRequestProfile) error {
	return s.userStorage.UpdateUserProfile(userID, &profile)
}
func (s *userService) ChangePassword(userID string, oldPassword, newPassword string) error {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		s.logger.Error(context.Background(), "пользователь не найден:", zap.Error(err))
		return errors.New("пользователь не найден")
	}

	if !user.IsVerified {
		s.logger.Error(context.Background(), "почта пользователя не подтверждена", zap.String("userID", userID))
		return errors.New("почта не подтверждена, смена пароля невозможна")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		s.logger.Error(context.Background(), "неверный старый пароль:", zap.Error(err))
		return errors.New("неверный старый пароль")
	}

	if !isValidPassword(newPassword) {
		s.logger.Error(context.Background(), "некорректный новый пароль", zap.String("userID", userID))
		return errors.New("некорректный новый пароль: минимум 8 символов, одна заглавная буква, цифра и спецсимвол")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(context.Background(), "ошибка при хешировании пароля:", zap.Error(err))
		return errors.New("ошибка при хешировании пароля")
	}

	if err := s.userStorage.UpdatePassword(userID, string(newHashedPassword)); err != nil {
		s.logger.Error(context.Background(), "ошибка при обновлении пароля:", zap.Error(err))
		return errors.New("ошибка при обновлении пароля")
	}

	return nil
}
func (s *userService) SendConfirmationCode(userID string) error {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	code, err := generateRandomCode()
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка генерации кода", zap.Error(err))
		return fmt.Errorf("временная ошибка")
	}

	redisKey := fmt.Sprintf("email_confirm:%s", userID)
	err = s.redis.Set(context.Background(), redisKey, code, 10*time.Minute)
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка записи в Redis", zap.Error(err))
		return fmt.Errorf("временная ошибка")
	}

	subject := "Подтверждение почты"
	body := fmt.Sprintf("Ваш код подтверждения: %s", code)
	err = s.mailer.Send(user.Email, subject, body)
	if err != nil {
		s.logger.Error(context.Background(), "Ошибка отправки письма", zap.Error(err))
		return fmt.Errorf("не удалось отправить письмо")
	}

	return nil
}

func (s *userService) VerifyConfirmationCode(userID string, inputCode string) error {
	redisKey := fmt.Sprintf("email_confirm:%s", userID)
	storedCode, err := s.redis.Get(context.Background(), redisKey)
	if errors.Is(err, redis.Nil) {
		return errors.New("код не найден или срок действия истек")
	} else if err != nil {
		s.logger.Error(context.Background(), "Ошибка получения кода из Redis", zap.Error(err))
		return fmt.Errorf("временная ошибка")
	}

	if storedCode != inputCode {
		return errors.New("неверный код")
	}

	if err := s.userStorage.MarkEmailVerified(userID); err != nil {
		s.logger.Error(context.Background(), "Ошибка обновления подтверждения", zap.Error(err))
		return fmt.Errorf("ошибка при подтверждении")
	}

	if err := s.redis.Del(context.Background(), redisKey); err != nil {
		s.logger.Error(context.Background(), "Ошибка удаления кода из Redis", zap.Error(err))
		return fmt.Errorf("ошибка удаления кода из Redis")
	}

	return nil
}

func (s *userService) ChangeEmail(userID, newEmail string) error {
	return s.userStorage.UpdateEmail(userID, newEmail, false)
}
