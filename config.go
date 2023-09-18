package ewa

import (
	"github.com/ewa-go/ewa/security"
	"github.com/ewa-go/ewa/session"
	"path/filepath"
)

type Config struct {
	Port           int
	Addr           string
	Secure         *Secure
	Authorization  security.Authorization
	Session        *session.Config
	Permission     *Permission
	Static         *Static
	NotFoundPage   string
	Views          *Views
	ContextHandler ContextHandler
	ErrorHandler   ErrorHandler
}

type Views struct {
	Root   string
	Layout string
	Engine string
}

type Static struct {
	Prefix string
	Root   string
}

type Secure struct {
	Path       string
	Key        string
	Cert       string
	ClientCert string
}

// Permission структура описывает разрешения на запрос
type Permission struct {
	AllRoutes            bool
	Handler              PermissionHandler
	NotPermissionHandler ErrorHandler
}

type Handler func(c *Context) error
type ContextHandler func(handler Handler) interface{}
type PermissionHandler func(c *Context, identity *security.Identity, method, path string) bool
type ErrorHandler func(c *Context, statusCode int, err interface{}) error

// Get Вернуть сертификат и ключ с путями
func (s *Secure) Get() (cert, key, clientCert string) {
	key = filepath.Join(s.Path, s.Key)
	cert = filepath.Join(s.Path, s.Cert)
	clientCert = filepath.Join(s.Path, s.ClientCert)
	return
}

// GetMutual Вернуть сертификат, ключ и файл с паролем с путями
/*func (s *Secure) GetMutual() (cert, key, clientCert string) {
	key = filepath.Join(s.Path, s.Key)
	cert = filepath.Join(s.Path, s.Cert)
	clientCert = filepath.Join(s.Path, s.ClientCert)
	return
}*/

type Timeout struct {
	Read  int
	Write int
	Idle  int
}

func NewTimeout(read, write, idle int) Timeout {
	return Timeout{
		Read:  read,
		Write: write,
		Idle:  idle,
	}
}

func (t Timeout) Get() (read int, write int, idle int) {
	return t.Read, t.Write, t.Idle
}

type BufferSize struct {
	Read  int
	Write int
}

func (b BufferSize) Get() (read int, write int) {
	return b.Read, b.Write
}

func NewBufferSize(read, write int) BufferSize {
	return BufferSize{
		Read:  read,
		Write: write,
	}
}
