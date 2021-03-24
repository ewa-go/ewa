package egowebapi

import "path/filepath"

type Config struct {
	Port    int
	Root string
	Secure  *Secure
	Timeout Timeout
}

type Secure struct {
	Path string
	Key  string
	Cert string
}

func (s *Secure) Get() (key string, cert string) {
	key = filepath.Join(s.Path, s.Key)
	cert = filepath.Join(s.Path, s.Cert)
	return key, cert
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
