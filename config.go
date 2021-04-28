package egowebapi

import "path/filepath"

type Config struct {
	Port      int
	Views     *Views
	Secure    *Secure
	Timeout   Timeout
	Static    string
	BasicAuth *BasicAuth
}

type Secure struct {
	Path string
	Key  string
	Cert string
}

type Views struct {
	Root   string
	Engine string
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
