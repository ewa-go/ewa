package session

import (
	"github.com/ewa-go/ewa/consts"
	"github.com/ewa-go/ewa/security"
	"github.com/google/uuid"
	"time"
)

// Config структура, которая описывает сессию
type Config struct {
	RedirectPath        string
	RedirectStatus      int
	AllRoutes           bool
	Expires             time.Duration
	SessionHandler      Handler
	GenSessionIdHandler GenSessionIdHandler
	KeyName             string
	m                   map[string]string
}

type Handler func(value string) (user string, err error)
type GenSessionIdHandler func() string

func (s *Config) Default() {

	// Хэш с данными
	s.m = map[string]string{}
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
	// Обработчик сессии
	if s.SessionHandler == nil {
		s.SessionHandler = func(value string) (user string, err error) {
			return s.get(value), nil
		}
	}
	// Обработчик генерации SessionId
	if s.GenSessionIdHandler == nil {
		s.GenSessionIdHandler = func() string {
			return uuid.New().String()
		}
	}
}

// Добавить запись в хэш сессий
func (s *Config) add(id string, user string) {
	s.m[id] = user
}

// Вернуть сессию из хэша
func (s *Config) get(id string) string {
	if value, ok := s.m[id]; ok {
		return value
	}
	return ""
}

// Удалить из хэша сессию
func (s *Config) delete(id string) {
	delete(s.m, id)
}

// Check Проверяем куки и извлекаем по ключу id по которому в бд/файле/памяти находим запись
func (s *Config) Check(value string) (*security.Identity, error) {

	user, err := s.SessionHandler(value)
	if err != nil {
		return nil, err
	}
	identity := &security.Identity{
		Username: user,
		AuthName: "Session",
	}

	return identity, nil

}
