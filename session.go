package egowebapi

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"net/http"
	"time"
)

// Session структура, которая описывает сессию
type Session struct {
	// Путь для перехода на страницу авторизации
	RedirectPath string
	// Статус при переходе, по умолчанию 302 - Found
	RedirectStatus int
	// Применение сессии на все маршруты
	AllRoutes bool
	// Время просрочки сессии
	Expires time.Duration
	// Обработчик сессии
	SessionHandler SessionHandler
	// Обработчик генерации SessionId
	GenSessionIdHandler GenSessionIdHandler
	// Обработчик на случай ошибки
	ErrorHandler ErrorHandler
}

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	Username string
	AuthName string
}

func (i Identity) String() string {
	return fmt.Sprintf("user: %s, auth_name: %s", i.Username, i.AuthName)
}

const sessionId = "session_id"

// Проверяем куки и извлекаем по ключу id по которому в бд/файле/памяти находим запись
func (s *Session) check(c *Context) error {

	key := c.Cookies(sessionId)
	if len(key) == 0 {
		return errors.New(fmt.Sprintf("Cookies [%s] not found", sessionId))
	}

	c.SessionId = key

	if s.SessionHandler != nil {
		user, err := s.SessionHandler(key)
		if err != nil {
			return err
		}
		c.Identity = &Identity{
			Username: user,
			AuthName: SessionAuth,
		}
	}

	return nil

	/*return func(ctx *fiber.Ctx) (err error) {

		user := unknown
		domain := unknown

		// Если cookie не существует, то перенаправляем запрос на условно "/signIn"
		key := ctx.Cookies(sessionId)
		if len(key) == 0 {
			return ctx.Redirect(s.RedirectPath)
		}

		// Проверяем на существование Handler
		if s.Handler != nil {
			user, domain, err = s.Handler(ctx.Cookies(sessionId))
			if err != nil {
				return ctx.Redirect(s.RedirectPath)
			}
		}

		// Получаем путь, чтобы передать в WebHandler
		route := ctx.Route()
		// Проверяем на существование PermissionHandler
		if isPermission && permission != nil && route != nil {
			if !permission.Handler(key, route.Path) {
				if s.ErrorHandler != nil {
					return s.ErrorHandler(ctx, fiber.StatusForbidden)
				}
				return ctx.SendStatus(fiber.StatusForbidden)
			}
		}

		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &Identity{
			Username:   user,
			Domain: domain,
		})
	}*/
}

// Формируем session_id и добавляем в куки
func (s *Session) signIn(c *Context) {

	key := utils.UUID()
	if s.GenSessionIdHandler != nil {
		key = s.GenSessionIdHandler()
	}
	c.SessionId = key
	if s.RedirectStatus == 0 {
		s.RedirectStatus = StatusFound
	}
	cookie := &http.Cookie{
		Name:    sessionId,
		Value:   key,
		Expires: time.Now().Add(s.Expires),
	}
	c.SetCookie(cookie)
}

// Очищаем куки, чтобы при маршрутизации сессия не была доступна
func (s *Session) signOut(c *Context) {

	key := c.Cookies(sessionId)
	c.SessionId = key
	identity := &Identity{
		AuthName: SessionAuth,
	}
	user, err := s.SessionHandler(key)
	if err != nil {
		identity = nil
	} else {
		identity.Username = user
	}
	c.Identity = identity
	c.ClearCookie(sessionId)
}
