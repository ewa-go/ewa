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
	User   string
	Domain string
	// Используется для идентификации сессии в cookie
	SessionId string
}

func (i Identity) String() string {
	return fmt.Sprintf("user: %s, domain: %s", i.User, i.Domain)
}

const sessionId = "session_id"

// Проверяем куки и извлекаем по ключу id по которому в бд/файле/памяти находим запись
func (s *Session) check(c *Context) error {

	key := c.Cookies(sessionId)
	if len(key) == 0 {
		return errors.New(fmt.Sprintf("Cookies [%s] not found", sessionId))
	}

	if s.SessionHandler != nil {
		user, domain, err := s.SessionHandler(key)
		if err != nil {
			return err
		}
		c.Identity = &Identity{
			User:      user,
			Domain:    domain,
			SessionId: key,
		}
	}

	return nil

	/*return func(ctx *fiber.Ctx) (err error) {

		user := unknown
		domain := unknown

		// Если cookie не существует, то перенаправляем запрос на условно "/login"
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
			User:   user,
			Domain: domain,
		})
	}*/
}

// Формируем session_id и добавляем в куки
func (s *Session) login(c *Context) {

	key := utils.UUID()
	if s.GenSessionIdHandler != nil {
		key = s.GenSessionIdHandler()
	}
	cookie := &http.Cookie{
		Name:    sessionId,
		Value:   key,
		Expires: time.Now().Add(s.Expires),
	}
	cookie.Name = sessionId
	cookie.Value = key
	cookie.Expires = time.Now().Add(s.Expires)
	c.SetCookie(cookie)
	c.Identity = &Identity{SessionId: key}
}

// Очищаем куки, чтобы при маршрутизации сессия не была доступна
func (s *Session) logout(c *Context) {

	key := c.Cookies(sessionId)
	identity := &Identity{
		SessionId: key,
	}
	user, domain, err := s.SessionHandler(key)
	if err != nil {
		identity = nil
	} else {
		identity.User = user
		identity.Domain = domain
	}
	c.Identity = identity
	c.ClearCookie(sessionId)
}
