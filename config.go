package egowebapi

import (
	"github.com/egovorukhin/egowebapi/swagger"
	"path/filepath"
)

type Config struct {
	Port          int
	Secure        *Secure
	Authorization Authorization
	Session       *Session
	Permission    *Permission
	Static        *Static
	NotFoundPage  string
	Views         *Views
	Swagger       *swagger.Config
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
	Path string
	Key  string
	Cert string
}

func (s *Secure) Get() (cert string, key string) {
	key = filepath.Join(s.Path, s.Key)
	cert = filepath.Join(s.Path, s.Cert)
	return cert, key
}

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
