package ewa

import (
	"time"

	"github.com/ewa-go/ewa/v1/consts"
	"github.com/google/uuid"
)

// Session структура, которая описывает сессию
type Session struct {
	RedirectPath         string
	RedirectStatus       int
	AllRoutes            bool
	Expires              time.Duration
	SessionHandler       SessionHandler
	DeleteSessionHandler DeleteSessionHandler
	GenSessionIdHandler  GenSessionIdHandler
	KeyName              string
	m                    map[string]string
}

type SessionHandler func(c *Context, value string) (user string, err error)
type GenSessionIdHandler func() string

// DeleteSessionHandler Возвращаемый флаг означает переход на страницу входа.
// True - переходить сразу на страницу входа. False - не переходить
type DeleteSessionHandler func(value string) bool

func (s *Session) Default() {

	// Хэш с данными
	s.m = make(map[string]string)
	// Имя ключа сессии
	if s.KeyName == "" {
		s.KeyName = "session_id"
	}
	// Путь для перехода на страницу авторизации
	if s.RedirectPath == "" {
		s.RedirectPath = "/login"
	}
	// Время просрочки сессии
	if s.Expires == 0 {
		s.Expires = 24 * time.Hour
	}
	// Статус при переходе, по умолчанию 302 - Found
	if s.RedirectStatus == 0 {
		s.RedirectStatus = consts.StatusFound
	}
	// Обработчик генерации SessionId
	if s.GenSessionIdHandler == nil {
		s.GenSessionIdHandler = func() string {
			return uuid.New().String()
		}
	}
}

// Добавить запись в хэш сессий
func (s *Session) add(id string, user string) {
	s.m[id] = user
}

// Вернуть сессию из хэша
func (s *Session) get(id string) string {
	if value, ok := s.m[id]; ok {
		return value
	}
	return ""
}

// Удалить из хэша сессию
func (s *Session) delete(id string) {
	delete(s.m, id)
}

// Check Проверяем куки и извлекаем по ключу id по которому в бд/файле/памяти находим запись
func (s *Session) Check(c *Context, value string) (*Identity, error) {

	user, err := s.SessionHandler(c, value)
	if err != nil {
		return nil, err
	}
	identity := &Identity{
		Username: user,
		AuthName: "Session",
	}

	return identity, nil

}
