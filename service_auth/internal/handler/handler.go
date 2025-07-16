package handler

import (
	"net/http"
	"service_auth/internal/models"
	"service_auth/internal/service"
	"service_auth/pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	UserService service.UserService
	logger      logger.Logger
}

func NewUserHandler(userService service.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		logger:      log,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var request models.UserRequestRegister

	if err := c.Bind(&request); err != nil {
		h.logger.Error(c.Request().Context(), "Неверный формат запроса:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	response, err := h.UserService.Register(request)
	if err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка при регистрации пользователя:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c echo.Context) error {
	var request models.UserRequestLogin

	if err := c.Bind(&request); err != nil {
		h.logger.Error(c.Request().Context(), "Неверный формат запроса:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	response, err := h.UserService.Login(request)
	if err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка при входе пользователя:", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error(c.Request().Context(), "Ошибка получения userID из контекста")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
	}

	profile, err := h.UserService.GetProfile(userID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка получения профиля:", zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Профиль не найден"})
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) ProfileUpdate(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error(c.Request().Context(), "Ошибка получения userID из контекста")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
	}

	var request models.UserRequestProfile
	if err := c.Bind(&request); err != nil {
		h.logger.Error(c.Request().Context(), "Неверный формат запроса:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	if err := h.UserService.ProfileUpdate(userID, request); err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка обновления профиля:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления профиля"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Профиль успешно обновлен"})
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error(c.Request().Context(), "Ошибка получения userID из контекста")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
	}

	var request models.UserRequestPassword
	if err := c.Bind(&request); err != nil {
		h.logger.Error(c.Request().Context(), "Неверный формат запроса:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	if request.OldPassword == "" || request.NewPassword == "" {
		h.logger.Error(c.Request().Context(), "Пароли не должны быть пустыми")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Оба пароля должны быть заполнены"})
	}

	if err := h.UserService.ChangePassword(userID, request.OldPassword, request.NewPassword); err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка смены пароля:", zap.Error(err))
		if err.Error() == "неверный старый пароль" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный старый пароль"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка смены пароля"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Пароль успешно изменен"})
}

func (h *UserHandler) SendCodeEmail(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error(c.Request().Context(), "Ошибка получения userID из контекста")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
	}

	if err := h.UserService.SendConfirmationCode(userID); err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка отправки кода подтверждения:", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось отправить код подтверждения"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Код подтверждения отправлен"})
}

func (h *UserHandler) ConfirmEmail(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error(c.Request().Context(), "Ошибка получения userID из контекста")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
	}

	var request models.UserRequestConfirmationCode
	if err := c.Bind(&request); err != nil {
		h.logger.Error(c.Request().Context(), "Неверный формат запроса:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	if err := h.UserService.VerifyConfirmationCode(userID, request.Code); err != nil {
		h.logger.Error(c.Request().Context(), "Ошибка подтверждения кода:", zap.Error(err))
		if err.Error() == "неверный код" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный код"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при подтверждении почты"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Почта успешно подтверждена"})
}

func (h *UserHandler) ChangeEmail(c echo.Context) error {
	var req models.ChangeEmailRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error(c.Request().Context(), "не удалось распарсить тело запроса", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат запроса")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "некорректный email")
	}

	// Предполагаем, что userID извлекается из контекста после middleware
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "не авторизован")
	}

	err := h.UserService.ChangeEmail(userID, req.NewEmail)
	if err != nil {
		h.logger.Error(c.Request().Context(), "не удалось изменить email", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка сервера")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "email успешно изменён"})
}
